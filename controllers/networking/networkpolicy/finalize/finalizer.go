package finalize

import (
	"context"
	networkingv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/networking/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"
)

type finalizer struct {
	client client.Client
	obj    *networkingv1alpha1.NetworkPolicy
}

func (f *finalizer) isNetworkPolicyBeingDeleted() bool {
	return !f.obj.ObjectMeta.DeletionTimestamp.IsZero()
}

func (f *finalizer) finalizeNetworkPolicy() (bool, error) {
	changed := false

	for _, ref := range f.obj.Status.NetworkPolicyRefs {
		if err := f.client.Delete(context.TODO(), &networkingv1.NetworkPolicy{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: ref.Namespace,
				Name:      f.obj.Name,
			},
		}); err != nil {
			return false, err
		}

		changed = true
	}

	return changed, nil
}

func (f *finalizer) EnsureNetworkPolicyFinalized() (reconcile.Result, error) {
	if !f.isNetworkPolicyBeingDeleted() {
		return reconcile.Continue()
	}

	if len(f.obj.Status.NetworkPolicyRefs) != 0 {
		changed, err := f.finalizeNetworkPolicy()
		if err != nil {
			return reconcile.RequeueWithError(err)
		} else if changed {
			return reconcile.RequeueAfter(2*time.Second, nil)
		}
	}

	controllerutil.RemoveFinalizer(f.obj, networkingv1alpha1.NetworkPolicyFinalizer)
	if err := f.client.Update(context.TODO(), f.obj); err != nil {
		return reconcile.RequeueWithError(err)
	}

	return reconcile.Stop()
}
