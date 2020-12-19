package workspaceclaim

import (
	"context"
	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile/task"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type reconciler interface {
	// status
	UpdateWorkspaceClaimStatus() (task.Result, error)

	// deletion
	EnsureFinalizerAppended() (task.Result, error)
	EnsureWorkspaceClaimDeletionProcessed() (task.Result, error)

	// possession
	EnsureWorkspaceClaimPossessionProcessed() (task.Result, error)
}

type reconcilerWrapper struct {
	client.Client

	ctx             context.Context
	log             logr.Logger
	scheme          *runtime.Scheme
	claim           *tenantv1alpha1.WorkspaceClaim
	workspaceTaints []corev1.Taint
}

func newReconciler(ctx context.Context, client client.Client, logger logr.Logger, scheme *runtime.Scheme, claim *tenantv1alpha1.WorkspaceClaim) reconciler {
	taints := NewWorkspaceTaints(claim.Spec.WorkspaceRef.Name)

	return &reconcilerWrapper{client, ctx, logger, scheme, claim, taints}
}
