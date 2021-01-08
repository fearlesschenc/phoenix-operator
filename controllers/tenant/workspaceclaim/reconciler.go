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

package workspaceclaim

import (
	"context"
	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"github.com/fearlesschenc/phoenix-operator/pkg/schedule"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/util/workqueue"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	ctrlreconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Reconciler reconciles a WorkspaceClaim object
type Reconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=tenant.phoenix.fearlesschenc.com,resources=workspaceclaims,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=tenant.phoenix.fearlesschenc.com,resources=workspaceclaims/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch;create;update;patch;delete

func (r *Reconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	logger := r.Log.WithValues("workspaceclaim", req.NamespacedName)

	// your logic here
	claim := &tenantv1alpha1.WorkspaceClaim{}
	if err := r.Get(ctx, req.NamespacedName, claim); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	reconciliation := newReconciliation(r.Client, logger, r.Scheme, claim)
	return reconcile.RunReconcileRoutine([]reconcile.SubroutineFunc{
		reconciliation.EnsureInitialized,
		reconciliation.EnsureValidated,
		reconciliation.UpdateStatus,
		reconciliation.EnsureFinalized,
		reconciliation.EnsurePossessionProcessed,
	})
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.TODO(), &corev1.Node{}, tenantv1alpha1.WorkspaceClaimOwnerKey, func(object runtime.Object) []string {
		node := object.(*corev1.Node)

		workspace := schedule.GetNodeWorkspace(node)
		if workspace != "" {
			return []string{workspace}
		}

		return nil
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&tenantv1alpha1.WorkspaceClaim{}).
		Watches(
			&source.Kind{Type: &corev1.Node{}},
			handler.Funcs{
				UpdateFunc: func(event event.UpdateEvent, limitingInterface workqueue.RateLimitingInterface) {
					oldNode := event.ObjectOld.(*corev1.Node)
					oldWorkspace := schedule.GetNodeWorkspace(oldNode)
					if oldWorkspace == "" {
						return
					}

					newNode := event.ObjectNew.(*corev1.Node)
					newWorkspace := schedule.GetNodeWorkspace(newNode)
					if oldWorkspace != newWorkspace {
						limitingInterface.Add(ctrlreconcile.Request{NamespacedName: types.NamespacedName{Name: oldWorkspace}})
						return
					}
				},
			}).
		Complete(r)
}
