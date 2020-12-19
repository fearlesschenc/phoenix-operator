package workspaceclaim

import (
	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile/task"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"
)

func (r *reconcilerWrapper) EnsureFinalizerAppended() (task.Result, error) {
	if controllerutil.ContainsFinalizer(r.claim, tenantv1alpha1.ClaimFinalizer) {
		return task.ContinueProcessing()
	}

	controllerutil.AddFinalizer(r.claim, tenantv1alpha1.ClaimFinalizer)
	if err := r.Update(r.ctx, r.claim); err != nil {
		return task.RequeueWithError(err)
	}

	return task.StopProcessing()
}

func (r *reconcilerWrapper) isWorkspaceClaimBeingDeleted() bool {
	return !r.claim.ObjectMeta.DeletionTimestamp.IsZero()
}

func (r *reconcilerWrapper) finalizeWorkspaceClaim() (bool, error) {
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

func (r *reconcilerWrapper) EnsureWorkspaceClaimDeletionProcessed() (task.Result, error) {
	if !r.isWorkspaceClaimBeingDeleted() {
		return task.ContinueProcessing()
	}

	if len(r.claim.Status.Node) != 0 {
		changed, err := r.finalizeWorkspaceClaim()
		if err != nil {
			return task.RequeueWithError(err)
		} else if changed {
			return task.RequeueAfter(2*time.Second, nil)
		}
	}

	controllerutil.RemoveFinalizer(r.claim, tenantv1alpha1.ClaimFinalizer)
	if err := r.Update(r.ctx, r.claim); err != nil {
		return task.RequeueWithError(err)
	}

	return task.StopProcessing()
}
