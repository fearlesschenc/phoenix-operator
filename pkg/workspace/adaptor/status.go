package adaptor

import (
	"github.com/fearlesschenc/phoenix-operator/pkg/util"
	v1 "k8s.io/api/core/v1"
)

type NodePossessionStatus struct {
	claimed   bool
	possessed bool
}

func (adaptor *Adaptor) UpdateWorkspaceClaimStatus() (util.OperationResult, error) {
	// Update Nodes
	nodeList := &v1.NodeList{}
	if err := adaptor.List(adaptor.ctx, nodeList); err != nil {
		return util.RequeueWithError(err)
	}

	var nodes []string
	for _, node := range nodeList.Items {
		if adaptor.isNodePossessedByWorkspace(node) {
			nodes = append(nodes, node.Name)
		}
	}

	adaptor.claim.Status.Node = nodes
	if err := adaptor.Status().Update(adaptor.ctx, adaptor.claim); err != nil {
		return util.RequeueWithError(err)
	}

	return util.ContinueProcessing()
}
