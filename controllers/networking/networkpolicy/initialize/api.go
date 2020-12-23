package initialize

import (
	networkingv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/networking/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Initializer interface {
	EnsureInitialized() (reconcile.Result, error)
}

func NewInitializer(client client.Client, obj *networkingv1alpha1.NetworkPolicy) Initializer {
	return &initializer{
		client: client,
		obj:    obj,
	}
}
