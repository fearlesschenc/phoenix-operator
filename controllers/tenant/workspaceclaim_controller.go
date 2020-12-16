/*
Authored by fearlesschenc@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controllers

import (
	"context"
	"github.com/fearlesschenc/phoenix-operator/pkg/util"
	adaptor2 "github.com/fearlesschenc/phoenix-operator/pkg/workspace/adaptor"
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"time"

	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
)

const WorkspaceClaimFinalizer = "workspaceclaim.finalizer.kubesphere.io"

// WorkspaceClaimReconciler reconciles a WorkspaceClaim object
type WorkspaceClaimReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=tenant.phoenix.fearlesschenc.com,resources=workspaceclaims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tenant.phoenix.fearlesschenc.com,resources=workspaceclaims/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch;create;update;patch;delete

func (r *WorkspaceClaimReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	logger := r.Log.WithValues("workspaceclaim", req.NamespacedName)

	// your logic here
	claim := &tenantv1alpha1.WorkspaceClaim{}
	if err := r.Get(ctx, req.NamespacedName, claim); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	adaptor := adaptor2.NewAdaptor(ctx, r.Client, logger, r.Scheme, claim)
	operations := []util.ReconcileOperation{
		adaptor.EnsureFinalizerAppended,
		adaptor.UpdateWorkspaceClaimStatus,
		adaptor.EnsureWorkspaceClaimDeletionProcessed,
		adaptor.EnsureWorkspaceClaimPossessionProcessed,
	}
	for _, operation := range operations {
		result, err := operation()
		if err != nil || result.RequeueRequest {
			return r.requeueAfter(result.RequeueDelay, err)
		}

		if result.CancelRequest {
			return r.doNotRequeue()
		}
	}

	return r.doNotRequeue()
}

func (r *WorkspaceClaimReconciler) doNotRequeue() (ctrl.Result, error) {
	return ctrl.Result{}, nil
}

func (r *WorkspaceClaimReconciler) requeueOnErr(err error) (ctrl.Result, error) {
	return ctrl.Result{}, err
}

func (r *WorkspaceClaimReconciler) requeueAfter(duration time.Duration, err error) (ctrl.Result, error) {
	return ctrl.Result{RequeueAfter: duration}, err
}

func (r *WorkspaceClaimReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tenantv1alpha1.WorkspaceClaim{}).
		// TODO
		// watch node label change, removal
		//Watches(&source.Kind{Type: &v1.Node{}}, &handler.EnqueueRequestsFromMapFunc{ToRequests: handler.ToRequestsFunc(r.mapNode)}).
		Complete(r)
}
