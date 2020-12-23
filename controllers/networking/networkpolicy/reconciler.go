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

package networkpolicy

import (
	"context"
	networkingv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/networking/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	ctrlreconcile "sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// Reconciler reconciles a NetworkPolicyHandler object
type Reconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=networking.phoenix.fearlesschenc.com,resources=networkpolicies,verbs=get;list;watch;create;update;patch;delete

func (r *Reconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	logger := r.Log.WithValues("networkpolicy", req.NamespacedName)

	policy := &networkingv1alpha1.NetworkPolicy{}
	if err := r.Get(ctx, req.NamespacedName, policy); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	reconciliation := newReconciliation(r.Client, logger, r.Scheme, policy)
	return reconcile.Run([]reconcile.TaskFunc{
		reconciliation.EnsureNetworkPolicyValidated,
		reconciliation.EnsureInitialized,
		reconciliation.EnsureNetworkPolicyFinalized,
		reconciliation.UpdateStatus,
		reconciliation.EnsureNetworkPolicyProcessed,
	})
}

func (r *Reconciler) handleNetworkPolicy(object handler.MapObject) []ctrlreconcile.Request {
	return []ctrlreconcile.Request{{types.NamespacedName{Name: object.Meta.GetName()}}}
}

func (r *Reconciler) filterNetworkPolicy(meta metav1.Object, object runtime.Object) bool {
	np := &networkingv1alpha1.NetworkPolicy{}
	if err := r.Client.Get(context.TODO(), types.NamespacedName{Name: meta.GetName()}, np); err != nil {
		return false
	}

	return true
}

func (r *Reconciler) handleNamespace(object handler.MapObject) []ctrlreconcile.Request {
	ret := []ctrlreconcile.Request{}
	npList := &networkingv1alpha1.NetworkPolicyList{}
	if err := r.Client.List(context.TODO(), npList); err != nil {
		return nil
	}

	for _, np := range npList.Items {
		ret = append(ret, ctrlreconcile.Request{NamespacedName: types.NamespacedName{Name: np.Name}})
	}

	return ret
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	if err := mgr.GetFieldIndexer().IndexField(context.TODO(), &networkingv1.NetworkPolicy{}, networkingv1alpha1.NetworkPolicyOwnerKey, func(object runtime.Object) []string {
		policy := object.(*networkingv1.NetworkPolicy)
		return []string{policy.Name}
	}); err != nil {
		return err
	}

	return ctrl.NewControllerManagedBy(mgr).
		For(&networkingv1alpha1.NetworkPolicy{}).
		Watches(
			&source.Kind{Type: &networkingv1.NetworkPolicy{}},
			&handler.EnqueueRequestsFromMapFunc{ToRequests: handler.ToRequestsFunc(r.handleNetworkPolicy)},
			builder.WithPredicates(predicate.NewPredicateFuncs(r.filterNetworkPolicy))).
		Watches(
			&source.Kind{Type: &corev1.Namespace{}},
			&handler.EnqueueRequestsFromMapFunc{ToRequests: handler.ToRequestsFunc(r.handleNamespace)}).
		Complete(r)
}
