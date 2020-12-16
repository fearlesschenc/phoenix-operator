package cluster

import (
	"context"
	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	workloadv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/workload/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/constant"
	v1 "k8s.io/api/networking/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func clusterNetworkIsolationEnabled(cluster *tenantv1alpha1.Cluster) bool {
	return cluster.Spec.NetworkIsolation != nil && *cluster.Spec.NetworkIsolation
}

func (r *Reconciler) removeNetworkPolicyIfExists(ctx context.Context, cluster *tenantv1alpha1.Cluster) error {
	if err := r.Delete(ctx, &v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: metav1.NamespaceDefault,
			Name:      cluster.Name,
		},
	}); err != nil {
		return client.IgnoreNotFound(err)
	}

	return nil
}

func (r *Reconciler) isolateApplication(ctx context.Context, application workloadv1alpha1.Application) error {
	policy := &v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: application.Name,
			Name:      constants.ClusterNetworkPolicy,
		},
		Spec: v1.NetworkPolicySpec{
			PolicyTypes: []v1.PolicyType{v1.PolicyTypeIngress},
			Ingress: []v1.NetworkPolicyIngressRule{
				{
					From: []v1.NetworkPolicyPeer{
						{
							NamespaceSelector: &metav1.LabelSelector{
								MatchLabels: map[string]string{
									tenantv1alpha1.ClusterLabel: application.Spec.Cluster,
								},
							},
						},
					},
				},
			},
		},
	}

	if err := r.Create(ctx, policy); err != nil {
		if apierrors.IsAlreadyExists(err) {
			return nil
		}

		return err
	}

	return nil
}

func (r *Reconciler) removeNamespaceIsolation(ctx context.Context, application workloadv1alpha1.Application) error {
	if err := r.Delete(ctx, &v1.NetworkPolicy{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: application.Name,
			Name:      constants.ClusterNetworkPolicy,
		},
	}); err != nil {
		return client.IgnoreNotFound(err)
	}

	return nil
}

func (r *Reconciler) reconcileNetwork(ctx context.Context, cluster *tenantv1alpha1.Cluster) error {
	if cluster.Status.NetworkIsolated == *cluster.Spec.NetworkIsolation {
		return nil
	}

	list := &workloadv1alpha1.ApplicationList{}
	if err := r.List(ctx, list, &client.ListOptions{
		LabelSelector: labels.SelectorFromSet(map[string]string{tenantv1alpha1.ClusterLabel: cluster.Name}),
	}); err != nil {
		return err
	}

	for _, app := range list.Items {
		if clusterNetworkIsolationEnabled(cluster) {
			return r.isolateApplication(ctx, app)
		}

		return r.removeNamespaceIsolation(ctx, app)
	}

	cluster.Status.NetworkIsolated = *cluster.Spec.NetworkIsolation
	return nil
}
