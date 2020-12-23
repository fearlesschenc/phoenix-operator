package workspaceclaim

import (
	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"
)

func (r *reconciliation) EnsureFinalizerAppended() (reconcile.Result, error) {
	if controllerutil.ContainsFinalizer(r.claim, tenantv1alpha1.WorkspaceClaimFinalizer) {
		return reconcile.Continue()
	}

	controllerutil.AddFinalizer(r.claim, tenantv1alpha1.WorkspaceClaimFinalizer)
	if err := r.Update(r.ctx, r.claim); err != nil {
		return reconcile.RequeueWithError(err)
	}

	return reconcile.Stop()
}

func (r *reconciliation) isWorkspaceClaimBeingDeleted() bool {
	return !r.claim.ObjectMeta.DeletionTimestamp.IsZero()
}

func (r *reconciliation) finalizeWorkspaceClaim() (bool, error) {
	changed := false

	for _, name := range r.claim.Status.Node {
		node := &corev1.Node{}
		if err := r.Get(r.ctx, types.NamespacedName{Name: name}, node); err != nil {
			return false, err
		}

		r.removeWorkspacePossessionOfNode(node)
		if err := r.Update(r.ctx, node); err != nil {
			return false, err
		}
		changed = true
	}

	return changed, nil
}

func (r *reconciliation) EnsureWorkspaceClaimDeletionProcessed() (reconcile.Result, error) {
	if !r.isWorkspaceClaimBeingDeleted() {
		return reconcile.Continue()
	}

	if len(r.claim.Status.Node) != 0 {
		changed, err := r.finalizeWorkspaceClaim()
		if err != nil {
			return reconcile.RequeueWithError(err)
		} else if changed {
			return reconcile.RequeueAfter(2*time.Second, nil)
		}
	}

	controllerutil.RemoveFinalizer(r.claim, tenantv1alpha1.WorkspaceClaimFinalizer)
	if err := r.Update(r.ctx, r.claim); err != nil {
		return reconcile.RequeueWithError(err)
	}

	return reconcile.Stop()
}
