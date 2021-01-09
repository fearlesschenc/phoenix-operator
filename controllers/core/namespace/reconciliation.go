package namespace

import (
	"github.com/fearlesschenc/phoenix-operator/controllers/core/namespace/finalize"
	"github.com/fearlesschenc/phoenix-operator/controllers/core/namespace/initialize"
	"github.com/fearlesschenc/phoenix-operator/controllers/core/namespace/validate"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Reconciliation struct {
	initialize.Initializer
	finalize.Finalizer
	validate.Validator

	client client.Client
	log    logr.Logger
	scheme *runtime.Scheme
	obj    *corev1.Namespace
}

func newReconciliation(client client.Client, logger logr.Logger, scheme *runtime.Scheme, obj *corev1.Namespace) *Reconciliation {
	r := &Reconciliation{
		client: client,
		log:    logger,
		scheme: scheme,
		obj:    obj,
	}

	r.Initializer = initialize.NewInitializer(client, logger, scheme, obj)
	r.Finalizer = finalize.NewFinalizer(client, obj)
	r.Validator = validate.NewValidator(client, logger, obj)

	return r
}
