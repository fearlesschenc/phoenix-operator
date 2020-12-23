package status

import (
	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Updater interface {
	UpdateStatus() (reconcile.Result, error)
}

func NewUpdater(client client.Client, obj *tenantv1alpha1.WorkspaceClaim) Updater {
	return &updater{
		client: client,
		obj:    obj,
	}
}
