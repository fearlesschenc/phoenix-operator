package workspace

import (
	"context"
	tenantv1alpha1 "github.com/fearlesschenc/kubesphere/pkg/apis/tenant/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Reconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *Reconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	logger := r.Log.WithValues("workspace", req.NamespacedName)

	workspace := &tenantv1alpha1.Workspace{}
	if err := r.Get(ctx, req.NamespacedName, workspace); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	reconciliation := newReconciliation(r.Client, logger, r.Scheme, workspace)
	return reconcile.RunReconcileRoutine([]reconcile.SubroutineFunc{
		reconciliation.EnsureInitialized,
		reconciliation.EnsureFinalized,
	})
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tenantv1alpha1.Workspace{}).
		Complete(r)
}
