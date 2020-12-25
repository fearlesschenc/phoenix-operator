package initialize

import (
	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type initializer struct {
	obj    *tenantv1alpha1.WorkspaceClaim
	client client.Client
}

func (init *initializer) EnsureInitialized() (reconcile.Result, error) {
	return reconcile.RunSubRoutine([]reconcile.SubroutineFunc{
		init.ensureFinalizerAppended,
	})
}
