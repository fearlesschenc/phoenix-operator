package initialize

import (
	tenantv1alpha1 "github.com/fearlesschenc/kubesphere/pkg/apis/tenant/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Initializer interface {
	EnsureInitialized() (reconcile.Result, error)
}

func NewInitializer(client client.Client, obj *tenantv1alpha1.Workspace) Initializer {
	return &initializer{
		obj:    obj,
		client: client,
	}
}
