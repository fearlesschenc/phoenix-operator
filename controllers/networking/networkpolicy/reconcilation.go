package networkpolicy

import (
	"context"
	networkingv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/networking/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/controllers/networking/networkpolicy/finalize"
	"github.com/fearlesschenc/phoenix-operator/controllers/networking/networkpolicy/initialize"
	"github.com/fearlesschenc/phoenix-operator/controllers/networking/networkpolicy/validate"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Reconciliation struct {
	initialize.Initializer
	finalize.Finalizer
	validate.Validator

	// Reconciler variable
	client.Client
	scheme *runtime.Scheme

	// reconciliation specific variable
	ctx context.Context
	log logr.Logger
	obj *networkingv1alpha1.NetworkPolicy
}

func newReconciliation(r *Reconciler, log logr.Logger, obj *networkingv1alpha1.NetworkPolicy) *Reconciliation {
	reconciliation := &Reconciliation{
		Client: r.Client,
		scheme: r.Scheme,
		log:    log,
		obj:    obj,
	}

	reconciliation.Initializer = initialize.NewInitializer(reconciliation.Client, reconciliation.obj)
	reconciliation.Finalizer = finalize.NewFinalizer(reconciliation.Client, reconciliation.obj)
	reconciliation.Validator = validate.NewValidator(reconciliation.obj)

	return reconciliation
}
