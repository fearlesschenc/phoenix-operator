package networkpolicy

import (
	"context"
	networkingv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/networking/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	networkingv1 "k8s.io/api/networking/v1"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sort"
)

const networkPolicyOwnerKey = ".meta.controller"

func (r *Reconciliation) UpdateStatus() (reconcile.Result, error) {
	networkPolicies := &networkingv1.NetworkPolicyList{}
	if err := r.List(context.TODO(), networkPolicies, client.MatchingFields{networkPolicyOwnerKey: r.obj.Name}); err != nil {
		return reconcile.RequeueWithError(err)
	}

	refs := []networkingv1alpha1.NetworkPolicyRef{}
	for _, policy := range networkPolicies.Items {
		refs = append(refs, networkingv1alpha1.NetworkPolicyRef{Namespace: policy.Namespace})
	}

	sort.SliceStable(refs, func(i, j int) bool {
		return refs[i].Namespace < refs[j].Namespace
	})

	if !reflect.DeepEqual(refs, r.obj.Status.NetworkPolicyRefs) {
		r.obj.Status.NetworkPolicyRefs = refs
		if err := r.Status().Update(context.TODO(), r.obj); err != nil {
			return reconcile.RequeueWithError(err)
		}

		return reconcile.Stop()
	}

	return reconcile.Continue()
}
