package schedule

import (
	"github.com/fearlesschenc/kubesphere/pkg/constants"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/kubernetes/pkg/util/labels"
	"k8s.io/kubernetes/pkg/util/taints"
)

func workspaceTaints(workspace string) []corev1.Taint {
	return []corev1.Taint{
		{
			Key:    constants.WorkspaceLabelKey,
			Value:  workspace,
			Effect: corev1.TaintEffectNoSchedule,
		},
		{
			Key:    constants.WorkspaceLabelKey,
			Value:  workspace,
			Effect: corev1.TaintEffectNoExecute,
		},
	}
}

func taintMatch(taint corev1.Taint, taintToMatch corev1.Taint) bool {
	return taint.Key == taintToMatch.Key &&
		taint.Effect == taintToMatch.Effect &&
		taint.Value == taintToMatch.Value
}

func taintExists(node corev1.Node, taint corev1.Taint) bool {
	for _, t := range node.Spec.Taints {
		if taintMatch(t, taint) {
			return true
		}
	}

	return false
}

func GetNodeWorkspace(node *corev1.Node) string {
	workspace := node.ObjectMeta.Labels[constants.WorkspaceLabelKey]
	if workspace == "" {
		return ""
	}

	for _, taint := range workspaceTaints(workspace) {
		if !taintExists(*node, taint) {
			return ""
		}
	}

	return workspace
}

//func IsNodePossessedByWorkspace(node *corev1.Node, workspace string) bool {
//	return taintExists(*node, workspaceTaints(workspace)...) &&
//		node.ObjectMeta.Labels[constants.WorkspaceLabelKey] == workspace
//}

func AddWorkspacePossessionOfNode(node *corev1.Node, workspace string) {
	// TODO: add taints
	for _, taint := range workspaceTaints(workspace) {
		if taintExists(*node, taint) {
			continue
		}

		node.Spec.Taints = append(node.Spec.Taints, taint)
	}

	// label
	node.Labels = labels.AddLabel(node.Labels, constants.WorkspaceLabelKey, workspace)
}

func RemoveWorkspacePossessionOfNode(node *corev1.Node, workspace string) {
	newTaints := node.Spec.Taints
	for _, taint := range workspaceTaints(workspace) {
		newTaints, _ = taints.DeleteTaint(newTaints, &taint)
	}
	node.Spec.Taints = newTaints

	// label
	if _, ok := node.ObjectMeta.Labels[constants.WorkspaceLabelKey]; ok {
		node.ObjectMeta.Labels = labels.CloneAndRemoveLabel(node.ObjectMeta.Labels, constants.WorkspaceLabelKey)
	}
}
