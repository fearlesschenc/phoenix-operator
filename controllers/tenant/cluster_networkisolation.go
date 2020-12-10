package controllers

import (
	"context"
	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	v1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func clusterNetworkIsolationEnabled(cluster *tenantv1alpha1.Cluster) bool {
	if cluster.Spec.NetworkIsolationEnabled == nil {
		return false
	}

	return *cluster.Spec.NetworkIsolationEnabled
}

func (r *ClusterReconciler) removeNetworkPolicyIfExists(ctx context.Context, cluster *tenantv1alpha1.Cluster) error {
	if !cluster.Status.NetworkIsolated {
		return nil
	}

	if err := r.Delete(ctx, &v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "",
			Name:      cluster.Name,
		},
	}); err != nil {
		return client.IgnoreNotFound(err)
	}

	cluster.Status.NetworkIsolated = false
	if err := r.Update(ctx, cluster); err != nil {
		return err
	}

	return nil
}

func (r *ClusterReconciler) createNetworkPolicyForCluster(cluster *tenantv1alpha1.Cluster) *v1.NetworkPolicy {
	return &v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Name: cluster.Name,
		},
		Spec: v1.NetworkPolicySpec{
			PodSelector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					ClusterLabel: cluster.Name,
				},
			},
			PolicyTypes: []v1.PolicyType{v1.PolicyTypeIngress},
			Ingress: []v1.NetworkPolicyIngressRule{
				{
					From: []v1.NetworkPolicyPeer{
						{
							PodSelector: &metav1.LabelSelector{
								MatchExpressions: []metav1.LabelSelectorRequirement{
									// allow same cluster pod communication
									// allow non-workspace pod, aka common metacluster components' pod.
									{Key: ClusterLabel, Operator: metav1.LabelSelectorOpIn, Values: []string{cluster.Name, ""}},
								},
							},
						},
					},
				},
			},
		},
	}
}

func (r *ClusterReconciler) reconcileNetwork(ctx context.Context, cluster *tenantv1alpha1.Cluster) error {
	if !clusterNetworkIsolationEnabled(cluster) {
		return r.removeNetworkPolicyIfExists(ctx, cluster)
	}

	if cluster.Status.NetworkIsolated {
		return nil
	}

	networkPolicy := r.createNetworkPolicyForCluster(cluster)
	if err := r.Create(ctx, networkPolicy); err != nil {
		return err
	}

	cluster.Status.NetworkIsolated = true
	if err := r.Status().Update(ctx, cluster); err != nil {
		return err
	}

	return nil
}
