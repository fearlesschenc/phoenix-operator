package networkpolicy

import (
	"context"
	networkingv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/networking/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/controllers/networking/networkpolicy/finalize"
	"github.com/fearlesschenc/phoenix-operator/controllers/networking/networkpolicy/initialize"
	"github.com/fearlesschenc/phoenix-operator/controllers/networking/networkpolicy/status"
	"github.com/fearlesschenc/phoenix-operator/controllers/networking/networkpolicy/validate"
	"github.com/fearlesschenc/phoenix-operator/pkg/constants"
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

type Reconciliation struct {
	initialize.Initializer
	finalize.Finalizer
	validate.Validator
	status.Updater

	// Reconciler variable
	client.Client
	scheme *runtime.Scheme

	// reconciliation specific variable
	ctx context.Context
	log logr.Logger
	obj *networkingv1alpha1.NetworkPolicy
}

func newReconciliation(r *Reconciler, log logr.Logger, obj *networkingv1alpha1.NetworkPolicy) *Reconciliation {
	reconciliation := &Reconciliation{
		Client: r.Client,
		scheme: r.Scheme,
		log:    log,
		obj:    obj,
	}

	reconciliation.Initializer = initialize.NewInitializer(reconciliation.Client, reconciliation.obj)
	reconciliation.Finalizer = finalize.NewFinalizer(reconciliation.Client, reconciliation.obj)
	reconciliation.Validator = validate.NewValidator(reconciliation.obj)
	reconciliation.Updater = status.NewUpdater(reconciliation.Client, reconciliation.obj)

	return reconciliation
}

func (r *Reconciliation) getWorkspaceNamespaces(workspace string, namespaceSelector *metav1.LabelSelector) ([]string, error) {
	namespaceSelector.MatchLabels[constants.WorkspaceLabelKey] = workspace
	selector, err := metav1.LabelSelectorAsSelector(namespaceSelector)
	if err != nil {
		return nil, errors.NewBadRequest(err.Error())
	}

	namespaceList := &corev1.NamespaceList{}
	if err := r.List(r.ctx, namespaceList, &client.ListOptions{LabelSelector: selector}); err != nil {
		return nil, err
	}

	namespaces := []string{}
	for _, item := range namespaceList.Items {
		namespaces = append(namespaces, item.Name)
	}

	return namespaces, nil
}

func (r *Reconciliation) getNetworkPolicySpecifiedNamespaces() ([]corev1.Namespace, error) {
	workspaceName := r.obj.GetLabels()[constants.WorkspaceLabelKey]

	namespaceLabelSelector := r.obj.Spec.NamespaceSelector.DeepCopy()
	namespaceLabelSelector.MatchLabels[constants.WorkspaceLabelKey] = workspaceName
	namespaceSelector, err := metav1.LabelSelectorAsSelector(namespaceLabelSelector)
	if err != nil {
		return nil, errors.NewBadRequest(err.Error())
	}

	namespaces := &corev1.NamespaceList{}
	if err := r.List(r.ctx, namespaces, &client.ListOptions{LabelSelector: namespaceSelector}); err != nil {
		return nil, err
	}

	return namespaces.Items, nil
}

func ensureNetworkPolicyAllowFrom(np *networkingv1.NetworkPolicy, namespaces []string) bool {
	updated := false

	for _, namespace := range namespaces {
		found := false
		for _, peer := range np.Spec.Ingress[0].From {
			// ignore error
			selector, _ := metav1.LabelSelectorAsSelector(peer.NamespaceSelector)
			if selector.Matches(labels.Set{constants.NamespaceLabelKey: namespace}) {
				found = true
			}
		}

		if !found {
			np.Spec.Ingress[0].From = append(np.Spec.Ingress[0].From, networkingv1.NetworkPolicyPeer{
				NamespaceSelector: metav1.SetAsLabelSelector(labels.Set{constants.NamespaceLabelKey: namespace}),
			})

			updated = true
		}
	}

	return updated
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
			if err := r.Get(r.ctx, types.NamespacedName{Namespace: ref.Namespace, Name: r.obj.Name}, policy); err != nil {
				return false, err
			}

			np = policy
		}
	}

	// exist
	if np != nil {
		changed := ensureNetworkPolicyAllowFrom(np, from)
		if changed {
			if err := r.Update(r.ctx, np); err != nil {
				return false, err
			}

			return true, nil
		}

		return false, nil
	}

	np = r.newNamespaceNetworkPolicy(namespace, from)
	if err := r.Create(r.ctx, np); err != nil {
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
			if err := r.Delete(r.ctx, &networkingv1.NetworkPolicy{
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
	if err := r.List(r.ctx, list, &client.ListOptions{LabelSelector: selector}); err != nil {
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
				Values:   []string{constants.SystemWorkspaceName},
			},
		},
	})

	list := &corev1.NamespaceList{}
	if err := r.List(r.ctx, list, &client.ListOptions{LabelSelector: selector}); err != nil {
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
	workspaceName := r.obj.GetLabels()[constants.WorkspaceLabelKey]

	// get specified targetNamespaces
	targetNamespaces, err := r.getWorkspaceNamespaces(workspaceName, r.obj.Spec.NamespaceSelector.DeepCopy())
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
