package adaptor

import (
	"github.com/fearlesschenc/phoenix-operator/pkg/constants"
	v1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/util/labels"
	"k8s.io/kubernetes/pkg/util/taints"
)

func NewWorkspaceTaints(workspace string) []v1.Taint {
	return []v1.Taint{
		{
			Key:    constants.WorkspaceLabel,
			Value:  workspace,
			Effect: v1.TaintEffectNoSchedule,
		},
		{
			Key:    constants.WorkspaceLabel,
			Value:  workspace,
			Effect: v1.TaintEffectNoExecute,
		},
	}
}

func taintMatch(taint v1.Taint, taintToMatch v1.Taint) bool {
	return taint.Key == taintToMatch.Key &&
		taint.Effect == taintToMatch.Effect &&
		taint.Value == taintToMatch.Value
}

func taintExists(node v1.Node, taints ...v1.Taint) bool {
	for _, t := range node.Spec.Taints {
		for _, taint := range taints {
			if taintMatch(t, taint) {
				return true
			}
		}
	}

	return false
}

func (adaptor *Adaptor) isNodePossessedByWorkspace(node v1.Node) bool {
	return taintExists(node, adaptor.workspaceTaints...) ||
		node.ObjectMeta.Labels[constants.WorkspaceLabel] == adaptor.claim.Spec.WorkspaceRef.Name
}

func (adaptor *Adaptor) removeWorkspacePossessionOfNode(node *v1.Node) {
	// taint
	newTaints := node.Spec.Taints
	for _, taint := range adaptor.workspaceTaints {
		newTaints, _ = taints.DeleteTaint(newTaints, &taint)
	}
	node.Spec.Taints = newTaints

	// label
	if _, ok := node.ObjectMeta.Labels[constants.WorkspaceLabel]; ok {
		node.ObjectMeta.Labels = labels.CloneAndRemoveLabel(node.ObjectMeta.Labels, constants.WorkspaceLabel)
	}
}

func (adaptor *Adaptor) addWorkspacePossessionOfNode(node *v1.Node) {
	// taint
	for _, taint := range adaptor.workspaceTaints {
		node.Spec.Taints = append(node.Spec.Taints, taint)
	}

	// label
	node.Labels = labels.AddLabel(node.Labels, constants.WorkspaceLabel, adaptor.claim.Spec.WorkspaceRef.Name)
}