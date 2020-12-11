package node

import (
	"fmt"
	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"time"
)

type ExclusiveOccupier struct {
}

var _ Occupier = &ExclusiveOccupier{}

func getLabel(node *corev1.Node, key string) (string, bool) {
	for k, v := range node.GetLabels() {
		if k == key {
			return v, true
		}
	}

	return "", false
}

func (o *ExclusiveOccupier) isNodeOccupyAble(clusterName string, node *corev1.Node) bool {
	val, exist := getLabel(node, tenantv1alpha1.ClusterLabel)
	if !exist || (exist && val == clusterName) {
		return true
	}

	return false
}

func addLabel(node *corev1.Node, key, val string) {
	node.Labels[key] = val
}

func findTaint(taints []corev1.Taint, taint corev1.Taint) int {
	for index, t := range taints {
		if t.Key == taint.Key && t.Value == taint.Value && t.Effect == taint.Effect {
			return index
		}
	}

	return -1
}

func addTaints(node *corev1.Node, taints ...corev1.Taint) {
	for _, taint := range taints {
		taint.TimeAdded.Time = time.Now()

		if index := findTaint(node.Spec.Taints, taint); index > 0 {
			continue
		}

		node.Spec.Taints = append(node.Spec.Taints, taint)
	}
}

func (o *ExclusiveOccupier) Occupy(clusterName string, node *corev1.Node) error {
	if !o.isNodeOccupyAble(clusterName, node) {
		return fmt.Errorf("unable to occupy node")
	}

	// two step
	// The first step is label node in order to let pod able
	// to schedule on it with nodeSelector in PodSpec
	addLabel(node, tenantv1alpha1.ClusterLabel, clusterName)

	// The second step is taint node, make pod able to run on it.
	addTaints(node, corev1.Taint{
		Key:    tenantv1alpha1.ClusterLabel,
		Value:  clusterName,
		Effect: corev1.TaintEffectNoExecute,
	}, corev1.Taint{
		Key:    tenantv1alpha1.ClusterLabel,
		Value:  clusterName,
		Effect: corev1.TaintEffectNoSchedule,
	})

	return nil
}

func removeLabel(node *corev1.Node, key string) {
	delete(node.Labels, key)
}

func removeTaints(node *corev1.Node, taints ...corev1.Taint) {
	for _, taint := range taints {
		index := -1

		if index = findTaint(node.Spec.Taints, taint); index < 0 {
			continue
		}

		// found
		node.Spec.Taints = append(node.Spec.Taints[:index], node.Spec.Taints[index+1:]...)
	}
}

func (o *ExclusiveOccupier) DeOccupy(clusterName string, node *corev1.Node) error {
	removeLabel(node, tenantv1alpha1.ClusterLabel)
	removeTaints(node, corev1.Taint{
		Key:    tenantv1alpha1.ClusterLabel,
		Value:  clusterName,
		Effect: corev1.TaintEffectNoExecute,
	}, corev1.Taint{
		Key:    tenantv1alpha1.ClusterLabel,
		Value:  clusterName,
		Effect: corev1.TaintEffectNoSchedule,
	})

	return nil
}

// TODO:
// inject pod:
// * nodeSelector
// * NoSchedule toleration
// * NoExecute toleration
