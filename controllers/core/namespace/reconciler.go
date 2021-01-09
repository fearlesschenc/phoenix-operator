package namespace

import (
	"context"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type Reconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups="",resources=namespaces/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=tenant.kubesphere.io,resources=workspaces,verbs=get
// +kubebuilder:rbac:groups=iam.kubesphere.io,resources=rolebases,verbs=list
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=roles,verbs=get;create;update
// +kubebuilder:rbac:groups=iam.kubesphere.io,resources=users,verbs=get
// +kubebuilder:rbac:groups=iam.kubesphere.io,resources=users,verbs=get
// +kubebuilder:rbac:groups=rbac.authorization.k8s.io,resources=rolebindings,verbs=create

func (r *Reconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	logger := r.Log.WithValues("namespace", req.NamespacedName)

	ns := &corev1.Namespace{}
	if err := r.Get(ctx, req.NamespacedName, ns); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	reconciliation := newReconciliation(r.Client, logger, r.Scheme, ns)
	return reconcile.RunReconcileRoutine([]reconcile.SubroutineFunc{
		reconciliation.EnsureValidated,
		reconciliation.EnsureInitialized,
		reconciliation.EnsureFinalized,
	})
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&corev1.Namespace{}).
		Complete(r)
}
