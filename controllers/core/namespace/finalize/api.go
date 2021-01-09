package finalize

import (
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Finalizer interface {
	EnsureFinalized() (reconcile.Result, error)
}

func NewFinalizer(client client.Client, obj *corev1.Namespace) Finalizer {
	return &finalizer{
		obj:    obj,
		client: client,
	}
}
