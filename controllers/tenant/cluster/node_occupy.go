package cluster

import (
	"context"
	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/node"
	corev1 "k8s.io/api/core/v1"
)

func (r *Reconciler) getNodeOccupyStatus(cluster *tenantv1alpha1.Cluster, name string) *node.OccupyStatus {
	for _, occupy := range cluster.Status.NodeOccupies {
		if occupy.NodeName == name {
			return &node.OccupyStatus{
				Occupied: true,
				Policy:   occupy.Policy,
			}
		}
	}

	return &node.OccupyStatus{Occupied: false}
}

func (r *Reconciler) reconcileNodeOccupy(ctx context.Context, cluster *tenantv1alpha1.Cluster, op *node.OccupyOperation) error {
	status := r.getNodeOccupyStatus(cluster, op.Node.Name)
	if op.ShouldSkipOn(status) {
		return nil
	}

	if err := op.OperateOn(status); err != nil {
		return err
	}

	if err := r.Update(ctx, op.Node); err != nil {
		return err
	}

	return nil
}

func (r *Reconciler) reconcileNodeOccupies(ctx context.Context, cluster *tenantv1alpha1.Cluster) error {
	nodeList := &corev1.NodeList{}
	if err := r.List(ctx, nodeList); err != nil {
		return err
	}

	var operations []*node.OccupyOperation
	for _, n := range nodeList.Items {
		operations = append(operations, &node.OccupyOperation{
			Cluster: cluster.Name,
			Node:    &n,
			ExpectStatus: &node.OccupyStatus{
				Occupied: false,
			},
		})
	}

	for _, occupy := range cluster.Spec.NodeOccupies {
		for _, op := range operations {
			if op.Node.Name == occupy.NodeName {
				op.ExpectStatus.Occupied = true
				op.ExpectStatus.Policy = occupy.Policy
			}
		}
	}

	for _, op := range operations {
		if err := r.reconcileNodeOccupy(ctx, cluster, op); err != nil {
			return err
		}
	}

	cluster.Status.NodeOccupies = cluster.Spec.NodeOccupies
	if err := r.Status().Update(ctx, cluster); err != nil {
		return err
	}

	return nil
}
