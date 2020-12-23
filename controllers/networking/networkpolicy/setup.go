package networkpolicy

import (
	"context"
	networkingv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/networking/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

func (r *Reconciler) handleNetworkPolicy(object handler.MapObject) []reconcile.Request {
	return []reconcile.Request{{types.NamespacedName{Name: object.Meta.GetName()}}}
}

func (r *Reconciler) filterNetworkPolicy(meta metav1.Object, object runtime.Object) bool {
	np := &networkingv1alpha1.NetworkPolicy{}
	if err := r.Client.Get(context.TODO(), types.NamespacedName{Name: meta.GetName()}, np); err != nil {
		return false
	}

	return true
}

func (r *Reconciler) handleNamespace(object handler.MapObject) []reconcile.Request {
	ret := []reconcile.Request{}
	npList := &networkingv1alpha1.NetworkPolicyList{}
	if err := r.Client.List(context.TODO(), npList); err != nil {
		return nil
	}

	for _, np := range npList.Items {
		ret = append(ret, reconcile.Request{NamespacedName: types.NamespacedName{Name: np.Name}})
	}

	return ret
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.TODO(), &networkingv1.NetworkPolicy{}, networkPolicyOwnerKey, func(object runtime.Object) []string {
		policy := object.(*networkingv1.NetworkPolicy)
		return []string{policy.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&networkingv1alpha1.NetworkPolicy{}).
		Watches(
			&source.Kind{Type: &networkingv1.NetworkPolicy{}},
			&handler.EnqueueRequestsFromMapFunc{ToRequests: handler.ToRequestsFunc(r.handleNetworkPolicy)},
			builder.WithPredicates(predicate.NewPredicateFuncs(r.filterNetworkPolicy))).
		Watches(
			&source.Kind{Type: &corev1.Namespace{}},
			&handler.EnqueueRequestsFromMapFunc{ToRequests: handler.ToRequestsFunc(r.handleNamespace)}).
		Complete(r)
}
