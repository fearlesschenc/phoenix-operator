package validate

import (
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Validator interface {
	EnsureValidated() (reconcile.Result, error)
}

func NewValidator(client client.Client, logger logr.Logger, obj *corev1.Namespace) Validator {
	return &validator{
		client: client,
		logger: logger,
		obj:    obj,
	}
}
