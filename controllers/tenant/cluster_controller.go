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
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
)

// ClusterReconciler reconciles a Cluster object
type ClusterReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *ClusterReconciler) reconcileWorkspace(ctx context.Context, templates []tenantv1alpha1.WorkspaceTemplate) error {
	workspaces := &tenantv1alpha1.WorkspaceList{}
	if err := r.List(ctx, workspaces); err != nil {
		return err
	}

	workspaceToPrune := make(map[int]bool)
	for i := 0; i < len(workspaces.Items); {
		workspaceToPrune[i] = true
	}

templateLoop:
	for _, template := range templates {
		for listIndex, workspace := range workspaces.Items {
			if workspace.Name == template.Name {
				workspaceToPrune[listIndex] = false

				if reflect.DeepEqual(workspace.Spec, template.Template) {
					continue templateLoop
				}

				// update workspace
				workspace.Spec = template.Template
				if err := r.Update(ctx, &workspace); err != nil {
					return err
				}

				continue templateLoop
			}
		}

		workspace := &tenantv1alpha1.Workspace{}
		workspace.Spec = template.Template
		// create workspace
		if err := r.Create(ctx, workspace); err != nil {
			return err
		}
	}

	// prune workspace
	for listIndex, prune := range workspaceToPrune {
		if prune {
			if err := r.Delete(ctx, &workspaces.Items[listIndex]); err != nil {
				return err
			}
		}
	}

	return nil
}

// +kubebuilder:rbac:groups=tenant.phoenix.fearlesschenc.com,resources=clusters,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tenant.phoenix.fearlesschenc.com,resources=clusters/status,verbs=get;update;patch

func (r *ClusterReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("cluster", req.NamespacedName)

	cluster := &tenantv1alpha1.Cluster{}
	if err := r.Get(ctx, req.NamespacedName, cluster); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if err := r.reconcileWorkspace(ctx, cluster.Spec.Workspace); err != nil {
		log.Error(err, "reconcile workspace error")
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

func (r *ClusterReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&tenantv1alpha1.Cluster{}).
		Complete(r)
}
