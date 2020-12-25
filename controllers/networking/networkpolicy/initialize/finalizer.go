package initialize

import (
	"context"
	networkingv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/networking/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (init *initializer) ensureFinalizerAppended() (reconcile.Result, error) {
	if controllerutil.ContainsFinalizer(init.obj, networkingv1alpha1.NetworkPolicyFinalizer) {
		return reconcile.Continue()
	}

	controllerutil.AddFinalizer(init.obj, networkingv1alpha1.NetworkPolicyFinalizer)
	if err := init.client.Update(context.TODO(), init.obj); err != nil {
		return reconcile.RequeueWithError(err)
	}

	return reconcile.Stop()
}
