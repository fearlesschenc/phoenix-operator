package initialize

import (
	"context"
	"github.com/fearlesschenc/phoenix-operator/apis/tenant"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (init *initializer) ensureFinalizerAppended() (reconcile.Result, error) {
	if controllerutil.ContainsFinalizer(init.obj, tenant.Finalizer) {
		return reconcile.Continue()
	}

	controllerutil.AddFinalizer(init.obj, tenant.Finalizer)
	if err := init.client.Update(context.TODO(), init.obj); err != nil {
		return reconcile.RequeueWithError(err)
	}

	return reconcile.Stop()
}
