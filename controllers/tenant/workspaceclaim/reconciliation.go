package workspaceclaim

import (
	"context"
	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/controllers/tenant/workspaceclaim/finalize"
	"github.com/fearlesschenc/phoenix-operator/controllers/tenant/workspaceclaim/initialize"
	"github.com/fearlesschenc/phoenix-operator/controllers/tenant/workspaceclaim/status"
	"github.com/fearlesschenc/phoenix-operator/controllers/tenant/workspaceclaim/validate"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"github.com/fearlesschenc/phoenix-operator/pkg/schedule"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"
)

// Reconciliation contains all of information that's needed to do
// one time Reconciliation
type Reconciliation struct {
	initialize.Initializer
	finalize.Finalizer
	status.Updater
	validate.Validator

	client client.Client
	log    logr.Logger
	scheme *runtime.Scheme
	obj    *tenantv1alpha1.WorkspaceClaim
}

func newReconciliation(client client.Client, logger logr.Logger, scheme *runtime.Scheme, obj *tenantv1alpha1.WorkspaceClaim) *Reconciliation {
	r := &Reconciliation{
		client: client,
		log:    logger,
		scheme: scheme,
		obj:    obj,
	}
	r.Initializer = initialize.NewInitializer(client, obj)
	r.Finalizer = finalize.NewFinalizer(client, obj)
	r.Updater = status.NewUpdater(client, obj)
	r.Validator = validate.NewValidator(obj)

	return r
}

type NodePossessionStatus struct {
	claimed   bool
	possessed bool
}

func (r *Reconciliation) EnsurePossessionProcessed() (reconcile.Result, error) {
	nodeList := &corev1.NodeList{}
	if err := r.client.List(context.TODO(), nodeList); err != nil {
		return reconcile.RequeueWithError(err)
	}

	possessionStatus := make(map[string]*NodePossessionStatus)
	for _, node := range nodeList.Items {
		possessionStatus[node.Name] = &NodePossessionStatus{claimed: false, possessed: false}
	}
	for _, node := range r.obj.Spec.Node {
		possessionStatus[node].claimed = true
	}
	for _, node := range r.obj.Status.Node {
		possessionStatus[node].possessed = true
	}

	changed := false
	for nodeName, nodePossessionStatus := range possessionStatus {
		if nodePossessionStatus.claimed == nodePossessionStatus.possessed {
			continue
		}

		node := &corev1.Node{}
		if err := r.client.Get(context.TODO(), types.NamespacedName{Name: nodeName}, node); err != nil {
			return reconcile.RequeueWithError(err)
		}

		if !nodePossessionStatus.claimed {
			schedule.RemoveWorkspacePossessionOfNode(node, r.obj.Spec.WorkspaceRef.Name)
		} else {
			schedule.AddWorkspacePossessionOfNode(node, r.obj.Spec.WorkspaceRef.Name)
		}

		if err := r.client.Update(context.TODO(), node); err != nil {
			return reconcile.RequeueWithError(err)
		}

		changed = true
	}

	if changed {
		return reconcile.RequeueAfter(2*time.Second, nil)
	}

	return reconcile.Continue()
}
