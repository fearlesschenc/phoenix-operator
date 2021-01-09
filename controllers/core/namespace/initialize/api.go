package initialize

import (
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Initializer interface {
	EnsureInitialized() (reconcile.Result, error)
}

func NewInitializer(client client.Client, logger logr.Logger, scheme *runtime.Scheme, obj *corev1.Namespace) Initializer {
	return &initializer{
		client: client,
		logger: logger,
		scheme: scheme,
		obj:    obj,
	}
}
