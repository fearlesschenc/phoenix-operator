package finalize

import (
	tenantv1alpha1 "github.com/fearlesschenc/kubesphere/pkg/apis/tenant/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Finalizer interface {
	EnsureFinalized() (reconcile.Result, error)
}

func NewFinalizer(client client.Client, obj *tenantv1alpha1.Workspace) Finalizer {
	return &finalizer{
		obj:    obj,
		client: client,
	}
}
