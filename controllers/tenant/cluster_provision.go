package controllers

import (
	"context"
	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
)

func (r *ClusterReconciler) removeClusterProvision(ctx context.Context, cluster *tenantv1alpha1.Cluster) error {
	nodes, err := r.getClusterProvisionedNodes(ctx, cluster)
	if err != nil {
		return err
	}

	// Remove node label
	for _, node := range nodes {
		delete(node.Labels, ClusterLabel)

		if err := r.Update(ctx, node); err != nil {
			return err
		}
	}

	return nil
}

func (r *ClusterReconciler) getClusterProvisionedNodes(ctx context.Context, cluster *tenantv1alpha1.Cluster) ([]*corev1.Node, error) {
	var ret []*corev1.Node

	for _, ref := range cluster.Status.OccupiedNodes {
		node, err := r.getNodesByRef(ctx, ref)
		if err != nil {
			return nil, err
		}

		ret = append(ret, node)
	}

	return ret, nil
}

func (r *ClusterReconciler) addClusterProvision(ctx context.Context, cluster *tenantv1alpha1.Cluster, nodes []*corev1.Node) error {
	for _, node := range nodes {
		node.Labels[ClusterLabel] = cluster.ObjectMeta.Name

		if err := r.Update(ctx, node); err != nil {
			return err
		}
	}

	return nil
}

func (r *ClusterReconciler) getNodesByRef(ctx context.Context, reference corev1.LocalObjectReference) (*corev1.Node, error) {
	node := &corev1.Node{}

	if err := r.Get(ctx, types.NamespacedName{Name: reference.Name, Namespace: ""}, node); err != nil {
		return nil, err
	}

	return node, nil
}

func nodeHaveAddress(node *corev1.Node, address string) bool {
	for _, nodeAddress := range node.Status.Addresses {
		if nodeAddress.Address == address {
			return true
		}
	}

	return false
}

func (r *ClusterReconciler) reconcileProvision(ctx context.Context, cluster *tenantv1alpha1.Cluster) error {
	if err := r.removeClusterProvision(ctx, cluster); err != nil {
		return err
	}

	nodeList := &corev1.NodeList{}
	if err := r.List(ctx, nodeList); err != nil {
		return err
	}

	var occupiedNodes []corev1.LocalObjectReference
	for _, address := range cluster.Spec.Nodes {
		for _, item := range nodeList.Items {
			node := &item

			if nodeHaveAddress(node, address) {
				node.Labels[ClusterLabel] = cluster.ObjectMeta.Name
				if err := r.Update(ctx, node); err != nil {
					return err
				}

				occupiedNodes = append(occupiedNodes, corev1.LocalObjectReference{Name: node.Name})
				break
			}
		}
	}

	cluster.Status.OccupiedNodes = occupiedNodes
	if err := r.Status().Update(ctx, cluster); err != nil {
		return err
	}

	return nil
}
