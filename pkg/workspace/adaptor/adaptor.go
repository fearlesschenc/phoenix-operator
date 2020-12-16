package adaptor

import (
	"context"
	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/util"
	"github.com/fearlesschenc/phoenix-operator/pkg/workspace"
	"github.com/go-logr/logr"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"
)

type Adaptor struct {
	client.Client

	ctx             context.Context
	log             logr.Logger
	scheme          *runtime.Scheme
	claim           *tenantv1alpha1.WorkspaceClaim
	workspaceTaints []v1.Taint
}

func NewAdaptor(ctx context.Context, client client.Client, logger logr.Logger, scheme *runtime.Scheme, claim *tenantv1alpha1.WorkspaceClaim) *Adaptor {
	taints := NewWorkspaceTaints(claim.Spec.WorkspaceRef.Name)

	return &Adaptor{client, ctx, logger, scheme, claim, taints}
}

func (adaptor *Adaptor) EnsureFinalizerAppended() (util.OperationResult, error) {
	if controllerutil.ContainsFinalizer(adaptor.claim, workspace.ClaimFinalizer) {
		return util.ContinueProcessing()
	}

	controllerutil.AddFinalizer(adaptor.claim, workspace.ClaimFinalizer)
	if err := adaptor.Update(adaptor.ctx, adaptor.claim); err != nil {
		return util.RequeueWithError(err)
	}

	return util.StopProcessing()
}

func (adaptor *Adaptor) isWorkspaceClaimBeingDeleted() bool {
	return !adaptor.claim.ObjectMeta.DeletionTimestamp.IsZero()
}

func (adaptor *Adaptor) FinalizeWorkspaceClaim() (bool, error) {
	changed := false

	for _, name := range adaptor.claim.Status.Node {
		node := &v1.Node{}
		if err := adaptor.Get(adaptor.ctx, types.NamespacedName{Name: name}, node); err != nil {
			return false, err
		}

		adaptor.removeWorkspacePossessionOfNode(node)
		if err := adaptor.Update(adaptor.ctx, node); err != nil {
			return false, err
		}
		changed = true
	}

	return changed, nil
}

func (adaptor *Adaptor) EnsureWorkspaceClaimDeletionProcessed() (util.OperationResult, error) {
	if !adaptor.isWorkspaceClaimBeingDeleted() {
		return util.ContinueProcessing()
	}

	if len(adaptor.claim.Status.Node) != 0 {
		changed, err := adaptor.FinalizeWorkspaceClaim()
		if err != nil {
			return util.RequeueWithError(err)
		} else if changed {
			return util.RequeueAfter(2*time.Second, nil)
		}
	}

	controllerutil.RemoveFinalizer(adaptor.claim, workspace.ClaimFinalizer)
	if err := adaptor.Update(adaptor.ctx, adaptor.claim); err != nil {
		return util.RequeueWithError(err)
	}

	return util.StopProcessing()
}

func (adaptor *Adaptor) EnsureWorkspaceClaimPossessionProcessed() (util.OperationResult, error) {
	nodeList := &v1.NodeList{}
	if err := adaptor.List(adaptor.ctx, nodeList); err != nil {
		return util.RequeueWithError(err)
	}

	possessionStatus := make(map[string]*NodePossessionStatus)
	for _, node := range nodeList.Items {
		possessionStatus[node.Name] = &NodePossessionStatus{claimed: false, possessed: false}
	}
	for _, node := range adaptor.claim.Spec.Node {
		possessionStatus[node].claimed = true
	}
	for _, node := range adaptor.claim.Status.Node {
		possessionStatus[node].possessed = true
	}

	changed := false
	for nodeName, status := range possessionStatus {
		if status.claimed == status.possessed {
			continue
		}

		updated := false
		node := &v1.Node{}
		if err := adaptor.Get(adaptor.ctx, types.NamespacedName{Name: nodeName}, node); err != nil {
			return util.RequeueWithError(err)
		}

		if !status.claimed {
			updated = adaptor.removeWorkspacePossessionOfNode(node)
		} else {
			updated = adaptor.addWorkspacePossessionOfNode(node)
		}

		if updated {
			if err := adaptor.Update(adaptor.ctx, node); err != nil {
				return util.RequeueWithError(err)
			}
			changed = true
		}
	}

	if changed {
		return util.RequeueAfter(2*time.Second, nil)
	}

	return util.ContinueProcessing()
}
