package initialize

import (
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type initializer struct {
	client client.Client
	logger logr.Logger
	scheme *runtime.Scheme
	obj    *corev1.Namespace
}

func (init *initializer) EnsureInitialized() (reconcile.Result, error) {
	return reconcile.RunSubRoutine([]reconcile.SubroutineFunc{
		init.ensureFinalizerAppended,
		init.ensureFieldsInitialized,
		init.ensureWorkspaceBound,
		init.ensureDefaultRoleCreated,
		init.ensureCreatorAdminRole,
	})
}
