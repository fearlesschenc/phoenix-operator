package finalize

import (
	"context"
	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"github.com/fearlesschenc/phoenix-operator/pkg/schedule"
	"github.com/fearlesschenc/phoenix-operator/pkg/utils"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"time"
)

type finalizer struct {
	obj    *tenantv1alpha1.WorkspaceClaim
	client client.Client
}

func (f *finalizer) isBeingDeleted() bool {
	return !f.obj.ObjectMeta.DeletionTimestamp.IsZero()
}

func (f *finalizer) finalizeWorkspaceClaim() (bool, error) {
	changed := false

	for _, name := range f.obj.Status.Node {
		node := &corev1.Node{}
		if err := f.client.Get(context.TODO(), types.NamespacedName{Name: name}, node); err != nil {
			return false, err
		}

		schedule.RemoveWorkspacePossessionOfNode(node, f.obj.Spec.WorkspaceRef.Name)
		if err := f.client.Update(context.TODO(), node); err != nil {
			return false, err
		}
		changed = true
	}

	return changed, nil
}

func (f *finalizer) EnsureFinalized() (reconcile.Result, error) {
	if !f.isBeingDeleted() {
		return reconcile.Continue()
	}

	if len(f.obj.Status.Node) != 0 {
		changed, err := f.finalizeWorkspaceClaim()
		if err != nil {
			return reconcile.RequeueWithError(err)
		} else if changed {
			return reconcile.RequeueAfter(2*time.Second, nil)
		}
	}

	controllerutil.RemoveFinalizer(f.obj, utils.TenantFinalizer)
	if err := f.client.Update(context.TODO(), f.obj); err != nil {
		return reconcile.RequeueWithError(err)
	}

	return reconcile.Stop()
}
