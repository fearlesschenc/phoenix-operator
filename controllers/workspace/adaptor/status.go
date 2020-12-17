package adaptor

import (
	"github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/util"
	v1 "k8s.io/api/core/v1"
	"reflect"
	"sort"
)

type NodePossessionStatus struct {
	claimed   bool
	possessed bool
}

func (adaptor *Adaptor) getLatestWorkspaceNodeStatus(status *v1alpha1.WorkspaceClaimStatus) error {
	// Update Nodes
	nodeList := &v1.NodeList{}
	if err := adaptor.List(adaptor.ctx, nodeList); err != nil {
		return err
	}

	nodes := make([]string, 0)
	for _, node := range nodeList.Items {
		if adaptor.isNodePossessedByWorkspace(node) {
			nodes = append(nodes, node.Name)
		}
	}
	sort.Strings(nodes)
	status.Node = nodes

	return nil
}

// UpdateWorkspaceClaimStatus initialize status of workspaceClaim
func (adaptor *Adaptor) UpdateWorkspaceClaimStatus() (util.OperationResult, error) {
	status := &v1alpha1.WorkspaceClaimStatus{}

	if err := adaptor.getLatestWorkspaceNodeStatus(status); err != nil {
		return util.RequeueWithError(err)
	}

	if !reflect.DeepEqual(status, &adaptor.claim.Status) {
		status.DeepCopyInto(&adaptor.claim.Status)
		if err := adaptor.Status().Update(adaptor.ctx, adaptor.claim); err != nil {
			return util.RequeueWithError(err)
		}

		return util.StopProcessing()
	}

	return util.ContinueProcessing()
}
