package finalize

import (
	"context"
	tenantv1alpha1 "github.com/fearlesschenc/kubesphere/pkg/apis/tenant/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/apis/tenant"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type finalizer struct {
	obj    *tenantv1alpha1.Workspace
	client client.Client
}

func (f *finalizer) isObjBeingDeleted() bool {
	return !f.obj.ObjectMeta.DeletionTimestamp.IsZero()
}

func (f *finalizer) EnsureFinalized() (reconcile.Result, error) {
	if !f.isObjBeingDeleted() {
		return reconcile.Continue()
	}

	controllerutil.RemoveFinalizer(f.obj, tenant.Finalizer)
	if err := f.client.Update(context.TODO(), f.obj); err != nil {
		return reconcile.RequeueWithError(err)
	}

	return reconcile.Stop()
}
