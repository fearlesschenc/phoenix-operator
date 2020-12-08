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
	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
)

const finalizer = "finalizers.phoenix.fearlesschenc.com"

// WorkspaceReconciler reconciles a Workspace object
type WorkspaceReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=tenant.phoenix.fearlesschenc.com,resources=workspaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tenant.phoenix.fearlesschenc.com,resources=workspaces/status,verbs=get;update;patch

func (r *WorkspaceReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	_ = r.Log.WithValues("workspace", req.NamespacedName)

	workspace := &tenantv1alpha1.Workspace{}
	if err := r.Get(ctx, req.NamespacedName, workspace); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if !controllerutil.ContainsFinalizer(workspace, finalizer) {
		controllerutil.AddFinalizer(workspace, finalizer)
		if err := r.Update(ctx, workspace); err != nil {
			return ctrl.Result{}, err
		}
	}

	// being deleted
	if !workspace.DeletionTimestamp.IsZero() {
		return r.reconcileDelete(ctx, workspace)
	}

	// TODO: handle network isolation
	if *workspace.Spec.NetworkIsolationEnabled {

	}

	// TODO: occupy hosts
	if len(workspace.Spec.Hosts) > 0 {

	}

	return ctrl.Result{}, nil
}

func (r *WorkspaceReconciler) reconcileDelete(ctx context.Context, workspace *tenantv1alpha1.Workspace) (ctrl.Result, error) {
	if controllerutil.ContainsFinalizer(workspace, finalizer) {
		controllerutil.RemoveFinalizer(workspace, finalizer)
		if err := r.Update(ctx, workspace); err != nil {
			return ctrl.Result{}, err
		}
	}

	return ctrl.Result{}, nil
}

func (r *WorkspaceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tenantv1alpha1.Workspace{}).
		Complete(r)
}
