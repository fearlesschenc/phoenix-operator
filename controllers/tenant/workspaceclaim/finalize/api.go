package finalize

import (
	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Finalizer interface {
	EnsureWorkspaceClaimFinalized() (reconcile.Result, error)
}

func NewFinalizer(client client.Client, obj *tenantv1alpha1.WorkspaceClaim) Finalizer {
	return &finalizer{
		obj:    obj,
		client: client,
	}
}
