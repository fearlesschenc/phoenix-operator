package workspaceclaim

import (
	"github.com/fearlesschenc/phoenix-operator/pkg/constants"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/kubernetes/pkg/util/labels"
	"k8s.io/kubernetes/pkg/util/taints"
	"time"
)

func NewWorkspaceTaints(workspace string) []corev1.Taint {
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

func taintExists(node corev1.Node, taints ...corev1.Taint) bool {
	for _, t := range node.Spec.Taints {
		for _, taint := range taints {
			if taintMatch(t, taint) {
				return true
			}
		}
	}

	return false
}

func (r *reconciliation) isNodePossessedByWorkspace(node corev1.Node) bool {
	return taintExists(node, r.workspaceTaints...) ||
		node.ObjectMeta.Labels[constants.WorkspaceLabelKey] == r.claim.Spec.WorkspaceRef.Name
}

func (r *reconciliation) removeWorkspacePossessionOfNode(node *corev1.Node) {
	// taint
	newTaints := node.Spec.Taints
	for _, taint := range r.workspaceTaints {
		newTaints, _ = taints.DeleteTaint(newTaints, &taint)
	}
	node.Spec.Taints = newTaints

	// label
	if _, ok := node.ObjectMeta.Labels[constants.WorkspaceLabelKey]; ok {
		node.ObjectMeta.Labels = labels.CloneAndRemoveLabel(node.ObjectMeta.Labels, constants.WorkspaceLabelKey)
	}
}

func (r *reconciliation) addWorkspacePossessionOfNode(node *corev1.Node) {
	// taint
	for _, taint := range r.workspaceTaints {
		node.Spec.Taints = append(node.Spec.Taints, taint)
	}

	// label
	node.Labels = labels.AddLabel(node.Labels, constants.WorkspaceLabelKey, r.claim.Spec.WorkspaceRef.Name)
}

func (r *reconciliation) EnsureWorkspaceClaimPossessionProcessed() (reconcile.Result, error) {
	nodeList := &corev1.NodeList{}
	if err := r.List(r.ctx, nodeList); err != nil {
		return reconcile.RequeueWithError(err)
	}

	possessionStatus := make(map[string]*NodePossessionStatus)
	for _, node := range nodeList.Items {
		possessionStatus[node.Name] = &NodePossessionStatus{claimed: false, possessed: false}
	}
	for _, node := range r.claim.Spec.Node {
		possessionStatus[node].claimed = true
	}
	for _, node := range r.claim.Status.Node {
		possessionStatus[node].possessed = true
	}

	changed := false
	for nodeName, status := range possessionStatus {
		if status.claimed == status.possessed {
			continue
		}

		node := &corev1.Node{}
		if err := r.Get(r.ctx, types.NamespacedName{Name: nodeName}, node); err != nil {
			return reconcile.RequeueWithError(err)
		}

		if !status.claimed {
			r.removeWorkspacePossessionOfNode(node)
		} else {
			r.addWorkspacePossessionOfNode(node)
		}

		if err := r.Update(r.ctx, node); err != nil {
			return reconcile.RequeueWithError(err)
		}
		changed = true
	}

	if changed {
		return reconcile.RequeueAfter(2*time.Second, nil)
	}

	return reconcile.Continue()
}
