package initialize

import (
	"context"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"github.com/fearlesschenc/phoenix-operator/pkg/utils"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (init *initializer) ensureFinalizerAppended() (reconcile.Result, error) {
	if controllerutil.ContainsFinalizer(init.obj, utils.NamespaceFinalizer) {
		return reconcile.Continue()
	}

	controllerutil.AddFinalizer(init.obj, utils.NamespaceFinalizer)
	if err := init.client.Update(context.TODO(), init.obj); err != nil {
		return reconcile.RequeueWithError(err)
	}

	return reconcile.Stop()
}
