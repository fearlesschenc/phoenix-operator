package initialize

import (
	"context"
	networkingv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/networking/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type initializer struct {
	obj    *networkingv1alpha1.NetworkPolicy
	client client.Client
}

func (init *initializer) EnsureInitialized() (reconcile.Result, error) {
	if controllerutil.ContainsFinalizer(init.obj, networkingv1alpha1.NetworkPolicyFinalizer) {
		return reconcile.Continue()
	}

	controllerutil.AddFinalizer(init.obj, networkingv1alpha1.NetworkPolicyFinalizer)
	if err := init.client.Update(context.TODO(), init.obj); err != nil {
		return reconcile.RequeueWithError(err)
	}

	return reconcile.Stop()
}
