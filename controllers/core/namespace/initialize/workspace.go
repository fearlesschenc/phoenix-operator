package initialize

import (
	"context"
	tenantv1alpha1 "github.com/fearlesschenc/kubesphere/pkg/apis/tenant/v1alpha1"
	"github.com/fearlesschenc/kubesphere/pkg/constants"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

func (init *initializer) isNamespaceFederated() bool {
	return init.obj.Labels[constants.KubefedManagedLabel] == "true"
}

func (init *initializer) isNamespaceControlledBy(workspace *tenantv1alpha1.Workspace) bool {
	return metav1.IsControlledBy(init.obj, workspace)
}

func (init *initializer) ensureWorkspaceBound() (reconcile.Result, error) {
	if init.isNamespaceFederated() {
		return reconcile.Continue()
	}

	workspaceName := init.obj.Labels[constants.WorkspaceLabelKey]
	workspace := &tenantv1alpha1.Workspace{}
	if err := init.client.Get(context.TODO(), types.NamespacedName{Name: workspaceName}, workspace); err != nil {
		return reconcile.RequeueWithError(err)
	}

	if !init.isNamespaceControlledBy(workspace) {
		init.obj.OwnerReferences = nil

		if err := controllerutil.SetControllerReference(workspace, init.obj, init.scheme); err != nil {
			return reconcile.RequeueWithError(err)
		}

		if err := init.client.Update(context.TODO(), init.obj); err != nil {
			return reconcile.RequeueWithError(err)
		}

		return reconcile.Stop()
	}

	return reconcile.Continue()
}
