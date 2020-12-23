package status

import (
	networkingv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/networking/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Updater interface {
	UpdateStatus() (reconcile.Result, error)
}

func NewUpdater(client client.Client, obj *networkingv1alpha1.NetworkPolicy) Updater {
	return &updater{
		obj:    obj,
		client: client,
	}
}
