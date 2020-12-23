package finalize

import (
	networkingv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/networking/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Finalizer interface {
	EnsureNetworkPolicyFinalized() (reconcile.Result, error)
}

func NewFinalizer(client client.Client, obj *networkingv1alpha1.NetworkPolicy) Finalizer {
	return &finalizer{
		client: client,
		obj:    obj,
	}
}
