package status

import (
	"context"
	networkingv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/networking/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	networkingv1 "k8s.io/api/networking/v1"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sort"
)

type updater struct {
	obj    *networkingv1alpha1.NetworkPolicy
	client client.Client
}

func (u *updater) UpdateStatus() (reconcile.Result, error) {
	networkPolicies := &networkingv1.NetworkPolicyList{}
	if err := u.client.List(context.TODO(), networkPolicies, client.MatchingFields{networkingv1alpha1.NetworkPolicyOwnerKey: u.obj.Name}); err != nil {
		return reconcile.RequeueWithError(err)
	}

	refs := []networkingv1alpha1.NetworkPolicyRef{}
	for _, policy := range networkPolicies.Items {
		refs = append(refs, networkingv1alpha1.NetworkPolicyRef{Namespace: policy.Namespace})
	}

	sort.SliceStable(refs, func(i, j int) bool {
		return refs[i].Namespace < refs[j].Namespace
	})

	if !reflect.DeepEqual(refs, u.obj.Status.NetworkPolicyRefs) {
		u.obj.Status.NetworkPolicyRefs = refs
		if err := u.client.Status().Update(context.TODO(), u.obj); err != nil {
			return reconcile.RequeueWithError(err)
		}

		return reconcile.Stop()
	}

	return reconcile.Continue()
}
