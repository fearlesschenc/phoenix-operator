package initialize

import (
	networkingv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/networking/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type initializer struct {
	obj    *networkingv1alpha1.NetworkPolicy
	client client.Client
}

func (init *initializer) EnsureInitialized() (reconcile.Result, error) {
	return reconcile.RunSubRoutine([]reconcile.SubroutineFunc{
		init.ensureFinalizerAppended,
		init.ensureFieldsInitialized,
	})
}
