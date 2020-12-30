package networkpolicy

import (
	"context"
	"github.com/fearlesschenc/kubesphere/pkg/constants"
	networkingv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/networking/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/controllers/networking/networkpolicy/finalize"
	"github.com/fearlesschenc/phoenix-operator/controllers/networking/networkpolicy/initialize"
	"github.com/fearlesschenc/phoenix-operator/controllers/networking/networkpolicy/status"
	"github.com/fearlesschenc/phoenix-operator/controllers/networking/networkpolicy/validate"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sort"
)

const SystemWorkspaceName = "system-workspace"

type Reconciliation struct {
	initialize.Initializer
	finalize.Finalizer
	validate.Validator
	status.Updater

	// Reconciler variable
	client client.Client
	scheme *runtime.Scheme
	log    logr.Logger
	obj    *networkingv1alpha1.NetworkPolicy
}

func newReconciliation(client client.Client, log logr.Logger, scheme *runtime.Scheme, obj *networkingv1alpha1.NetworkPolicy) *Reconciliation {
	r := &Reconciliation{
		client: client,
		scheme: scheme,
		log:    log,
		obj:    obj,
	}

	r.Initializer = initialize.NewInitializer(r.client, r.obj)
	r.Finalizer = finalize.NewFinalizer(r.client, r.obj)
	r.Validator = validate.NewValidator(r.obj)
	r.Updater = status.NewUpdater(r.client, r.obj)

	return r
}

func (r *Reconciliation) getWorkspaceNamespaces(workspace string, namespaceSelector *metav1.LabelSelector) ([]string, error) {
	if namespaceSelector.MatchLabels == nil {
		namespaceSelector.MatchLabels = map[string]string{}
	}

	namespaceSelector.MatchLabels[constants.WorkspaceLabelKey] = workspace
	selector, err := metav1.LabelSelectorAsSelector(namespaceSelector)
	if err != nil {
		return nil, errors.NewBadRequest(err.Error())
	}

	namespaceList := &corev1.NamespaceList{}
	if err := r.client.List(context.TODO(), namespaceList, &client.ListOptions{LabelSelector: selector}); err != nil {
		return nil, err
	}

	namespaces := []string{}
	for _, item := range namespaceList.Items {
		namespaces = append(namespaces, item.Name)
	}

	return namespaces, nil
}

func isPeerMatchFromNamespaces(peers []networkingv1.NetworkPolicyPeer, namespaces []string) bool {
	if len(peers) != len(namespaces) {
		return false
	}

	for i, peer := range peers {
		selector, _ := metav1.LabelSelectorAsSelector(peer.NamespaceSelector)
		if !selector.Matches(labels.Set{constants.NamespaceLabelKey: namespaces[i]}) {
			return false
		}
	}

	return true
}

func ensureNetworkPolicyAllowFrom(np *networkingv1.NetworkPolicy, namespaces []string) bool {
	if isPeerMatchFromNamespaces(np.Spec.Ingress[0].From, namespaces) {
		return false
	}

	np.Spec.Ingress[0].From = make([]networkingv1.NetworkPolicyPeer, 0)
	for _, namespace := range namespaces {
		np.Spec.Ingress[0].From = append(np.Spec.Ingress[0].From, networkingv1.NetworkPolicyPeer{
			NamespaceSelector: metav1.SetAsLabelSelector(labels.Set{constants.NamespaceLabelKey: namespace}),
		})
	}

	return true
}

func (r *Reconciliation) newNamespaceNetworkPolicy(namespace string, from []string) *networkingv1.NetworkPolicy {
	rule := &networkingv1.NetworkPolicyIngressRule{}
	rule.From = make([]networkingv1.NetworkPolicyPeer, 0)

	np := &networkingv1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name:      r.obj.Name,
			Namespace: namespace,
		},
		Spec: networkingv1.NetworkPolicySpec{
			PodSelector: *metav1.SetAsLabelSelector(labels.Set{}),
			Ingress:     []networkingv1.NetworkPolicyIngressRule{*rule},
		},
	}
	ensureNetworkPolicyAllowFrom(np, from)

	return np
}

func (r *Reconciliation) ensureTargetAllowAccessFrom(namespace string, from []string) (bool, error) {
	var np *networkingv1.NetworkPolicy

	for _, ref := range r.obj.Status.NetworkPolicyRefs {
		if ref.Namespace == namespace {
			policy := &networkingv1.NetworkPolicy{}
			if err := r.client.Get(context.TODO(), types.NamespacedName{Namespace: ref.Namespace, Name: r.obj.Name}, policy); err != nil {
				return false, err
			}

			np = policy
		}
	}

	// exist
	if np != nil {
		changed := ensureNetworkPolicyAllowFrom(np, from)
		if changed {
			if err := r.client.Update(context.TODO(), np); err != nil {
				return false, err
			}

			return true, nil
		}

		return false, nil
	}

	np = r.newNamespaceNetworkPolicy(namespace, from)
	if err := r.client.Create(context.TODO(), np); err != nil {
		return false, err
	}

	return true, nil
}

func (r *Reconciliation) ensureUnspecifiedNamespaceNetworkPolicyDeleted(namespaces []string) (bool, error) {
	changed := false

	for _, ref := range r.obj.Status.NetworkPolicyRefs {
		found := false

		for _, ns := range namespaces {
			if ref.Namespace == ns {
				found = true
				break
			}
		}

		if !found {
			if err := r.client.Delete(context.TODO(), &networkingv1.NetworkPolicy{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: ref.Namespace,
					Name:      r.obj.Name,
				},
			}); err != nil && !errors.IsNotFound(err) {
				return false, err
			}

			changed = true
		}
	}

	return changed, nil
}

func (r *Reconciliation) addNonWorkspaceNamespace(namespaces []string) ([]string, error) {
	selector, _ := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      constants.WorkspaceLabelKey,
				Operator: metav1.LabelSelectorOpDoesNotExist,
			},
		},
	})

	list := &corev1.NamespaceList{}
	if err := r.client.List(context.TODO(), list, &client.ListOptions{LabelSelector: selector}); err != nil {
		return nil, err
	}

	for _, item := range list.Items {
		namespaces = append(namespaces, item.Name)
	}

	return namespaces, nil
}

func (r *Reconciliation) addSystemWorkspaceNamespace(namespaces []string) ([]string, error) {
	selector, _ := metav1.LabelSelectorAsSelector(&metav1.LabelSelector{
		MatchExpressions: []metav1.LabelSelectorRequirement{
			{
				Key:      constants.WorkspaceLabelKey,
				Operator: metav1.LabelSelectorOpIn,
				Values:   []string{SystemWorkspaceName},
			},
		},
	})

	list := &corev1.NamespaceList{}
	if err := r.client.List(context.TODO(), list, &client.ListOptions{LabelSelector: selector}); err != nil {
		return nil, err
	}

	for _, item := range list.Items {
		namespaces = append(namespaces, item.Name)
	}

	return namespaces, nil
}

func (r *Reconciliation) getNetworkPolicySpecifiedFromNamespaces() ([]string, error) {
	fromNamespaces := make([]string, 0, len(r.obj.Spec.From)-1)
	for _, peer := range r.obj.Spec.From {
		namespaces, err := r.getWorkspaceNamespaces(peer.Workspace, peer.NamespaceSelector.DeepCopy())
		if err != nil {
			return nil, err
		}

		fromNamespaces = append(fromNamespaces, namespaces...)
	}

	var err error

	fromNamespaces, err = r.addNonWorkspaceNamespace(fromNamespaces)
	if err != nil {
		return nil, err
	}
	fromNamespaces, err = r.addSystemWorkspaceNamespace(fromNamespaces)
	if err != nil {
		return nil, err
	}

	sort.Strings(fromNamespaces)
	return fromNamespaces, nil
}

func (r *Reconciliation) EnsureNetworkPolicyProcessed() (reconcile.Result, error) {
	// get specified targetNamespaces
	targetNamespaces, err := r.getWorkspaceNamespaces(r.obj.Spec.Workspace, r.obj.Spec.NamespaceSelector.DeepCopy())
	if err != nil {
		if errors.IsBadRequest(err) {
			return reconcile.Stop()
		}

		return reconcile.RequeueWithError(err)
	}

	updated := false

	// delete unspecified namespace obj
	updated, err = r.ensureUnspecifiedNamespaceNetworkPolicyDeleted(targetNamespaces)
	if err != nil {
		return reconcile.RequeueWithError(err)
	} else if updated {
		return reconcile.Stop()
	}

	fromNamespaces, err := r.getNetworkPolicySpecifiedFromNamespaces()
	if err != nil {
		return reconcile.RequeueWithError(err)
	}

	for _, namespace := range targetNamespaces {
		changed, err := r.ensureTargetAllowAccessFrom(namespace, fromNamespaces)
		if err != nil {
			return reconcile.RequeueWithError(err)
		}

		if changed {
			updated = true
		}
	}

	if updated {
		return reconcile.Stop()
	}

	return reconcile.Continue()
}
