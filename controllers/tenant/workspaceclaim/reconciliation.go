package workspaceclaim

import (
	"context"
	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// Reconciliation contains the descriptive steps that
// should be done in one reconciliation
type Reconciliation interface {
	// status
	UpdateWorkspaceClaimStatus() (reconcile.Result, error)
	// deletion
	EnsureFinalizerAppended() (reconcile.Result, error)
	EnsureWorkspaceClaimDeletionProcessed() (reconcile.Result, error)
	// possession
	EnsureWorkspaceClaimPossessionProcessed() (reconcile.Result, error)
}

// reconciliation contains all of information that's needed to do
// one time reconciliation
type reconciliation struct {
	client.Client

	ctx             context.Context
	log             logr.Logger
	scheme          *runtime.Scheme
	claim           *tenantv1alpha1.WorkspaceClaim
	workspaceTaints []corev1.Taint
}

func newReconciliation(ctx context.Context, client client.Client, logger logr.Logger, scheme *runtime.Scheme, claim *tenantv1alpha1.WorkspaceClaim) Reconciliation {
	taints := NewWorkspaceTaints(claim.Spec.WorkspaceRef.Name)

	return &reconciliation{client, ctx, logger, scheme, claim, taints}
}
