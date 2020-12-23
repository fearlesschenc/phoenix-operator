package status

import (
	"context"
	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	corev1 "k8s.io/api/core/v1"
	"reflect"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sort"
)

type updater struct {
	client client.Client
	obj    *tenantv1alpha1.WorkspaceClaim
}

func (u *updater) getLatestWorkspaceNodeStatus(status *tenantv1alpha1.WorkspaceClaimStatus) error {
	// Update Nodes
	nodeList := &corev1.NodeList{}
	if err := u.client.List(context.TODO(), nodeList, client.MatchingFields{tenantv1alpha1.WorkspaceClaimOwnerKey: u.obj.Spec.WorkspaceRef.Name}); err != nil {
		return err
	}

	nodes := make([]string, 0)
	for _, node := range nodeList.Items {
		nodes = append(nodes, node.Name)
	}
	sort.Strings(nodes)
	status.Node = nodes

	return nil
}

// UpdateWorkspaceClaimStatus initialize status of workspaceClaim
func (u *updater) UpdateStatus() (reconcile.Result, error) {
	status := &tenantv1alpha1.WorkspaceClaimStatus{}

	if err := u.getLatestWorkspaceNodeStatus(status); err != nil {
		return reconcile.RequeueWithError(err)
	}

	if !reflect.DeepEqual(status, &u.obj.Status) {
		u.obj.Status = *status
		if err := u.client.Status().Update(context.TODO(), u.obj); err != nil {
			return reconcile.RequeueWithError(err)
		}

		return reconcile.Stop()
	}

	return reconcile.Continue()
}
