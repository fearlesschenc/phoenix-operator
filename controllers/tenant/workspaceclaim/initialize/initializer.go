package initialize

import (
	"context"
	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type initializer struct {
	obj    *tenantv1alpha1.WorkspaceClaim
	client client.Client
}

func (init *initializer) EnsureInitialized() (reconcile.Result, error) {
	if controllerutil.ContainsFinalizer(init.obj, tenantv1alpha1.WorkspaceClaimFinalizer) {
		return reconcile.Continue()
	}

	controllerutil.AddFinalizer(init.obj, tenantv1alpha1.WorkspaceClaimFinalizer)
	if err := init.client.Update(context.TODO(), init.obj); err != nil {
		return reconcile.RequeueWithError(err)
	}

	return reconcile.Stop()
}
