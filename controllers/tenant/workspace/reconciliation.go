package workspace

import (
	tenantv1alpha1 "github.com/fearlesschenc/kubesphere/pkg/apis/tenant/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/controllers/tenant/workspace/finalize"
	"github.com/fearlesschenc/phoenix-operator/controllers/tenant/workspace/initialize"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Reconciliation struct {
	initialize.Initializer
	finalize.Finalizer

	client client.Client
	log    logr.Logger
	scheme *runtime.Scheme
	obj    *tenantv1alpha1.Workspace
}

func newReconciliation(client client.Client, logger logr.Logger, scheme *runtime.Scheme, obj *tenantv1alpha1.Workspace) *Reconciliation {
	r := &Reconciliation{
		client: client,
		log:    logger,
		scheme: scheme,
		obj:    obj,
	}

	r.Initializer = initialize.NewInitializer(client, obj)
	r.Finalizer = finalize.NewFinalizer(client, obj)

	return r
}
