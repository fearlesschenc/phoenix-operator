package workspaceclaim

import (
	"github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile/task"
	v1 "k8s.io/api/core/v1"
	"reflect"
	"sort"
)

type NodePossessionStatus struct {
	claimed   bool
	possessed bool
}

func (r *reconcilerWrapper) getLatestWorkspaceNodeStatus(status *v1alpha1.WorkspaceClaimStatus) error {
	// Update Nodes
	nodeList := &v1.NodeList{}
	if err := r.List(r.ctx, nodeList); err != nil {
		return err
	}

	nodes := make([]string, 0)
	for _, node := range nodeList.Items {
		if r.isNodePossessedByWorkspace(node) {
			nodes = append(nodes, node.Name)
		}
	}
	sort.Strings(nodes)
	status.Node = nodes

	return nil
}

// UpdateWorkspaceClaimStatus initialize status of workspaceClaim
func (r *reconcilerWrapper) UpdateWorkspaceClaimStatus() (task.Result, error) {
	status := &v1alpha1.WorkspaceClaimStatus{}

	if err := r.getLatestWorkspaceNodeStatus(status); err != nil {
		return task.RequeueWithError(err)
	}

	if !reflect.DeepEqual(status, &r.claim.Status) {
		status.DeepCopyInto(&r.claim.Status)
		if err := r.Status().Update(r.ctx, r.claim); err != nil {
			return task.RequeueWithError(err)
		}

		return task.StopProcessing()
	}

	return task.ContinueProcessing()
}
