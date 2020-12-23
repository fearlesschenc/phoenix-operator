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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
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

	reconciliation := newReconciliation(r, logger, policy)
	return reconcile.Run([]reconcile.TaskFunc{
		reconciliation.EnsureNetworkPolicyValidated,
		reconciliation.EnsureInitialized,
		reconciliation.EnsureNetworkPolicyFinalized,
		reconciliation.UpdateStatus,
		reconciliation.EnsureNetworkPolicyProcessed,
	})
}
