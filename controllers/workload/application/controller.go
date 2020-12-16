package application

import (
	"context"
	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	workloadv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/workload/v1alpha1"
	"github.com/go-logr/logr"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/apiutil"
)

type Reconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

func (r *Reconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("application", req.Name)

	app := &workloadv1alpha1.Application{}
	if err := r.Get(ctx, req.NamespacedName, app); err != nil {
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	if app.Spec.Cluster == "" {
		log.Info("invalid application")
		return ctrl.Result{}, nil
	}

	cluster := &tenantv1alpha1.Cluster{}
	if err := r.Get(ctx, types.NamespacedName{Name: app.Spec.Cluster}, cluster); err != nil {
		return ctrl.Result{}, err
	}

	if !metav1.IsControlledBy(app, cluster) {
		gvk, err := apiutil.GVKForObject(cluster, r.Scheme)
		if err != nil {
			return ctrl.Result{}, err
		}

		ref := metav1.NewControllerRef(cluster, gvk)
		refs := app.GetOwnerReferences()
		refs = append(refs, *ref)

		app.SetOwnerReferences(refs)
	}

	return ctrl.Result{}, nil
}

func (r *Reconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&workloadv1alpha1.Application{}).
		Complete(r)
}
