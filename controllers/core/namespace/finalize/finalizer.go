package finalize

import (
	"context"
	"github.com/fearlesschenc/kubesphere/pkg/constants"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	"github.com/fearlesschenc/phoenix-operator/pkg/utils"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type finalizer struct {
	obj    *corev1.Namespace
	client client.Client
}

func (f *finalizer) isObjBeingDeleted() bool {
	return !f.obj.ObjectMeta.DeletionTimestamp.IsZero()
}

func (f *finalizer) ensureRouteRemoved() (reconcile.Result, error) {
	name := constants.IngressControllerPrefix + f.obj.Name
	namespace := constants.IngressControllerNamespace

	if err := f.client.Delete(context.TODO(), &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
	}); err != nil && !errors.IsNotFound(err) {
		return reconcile.RequeueWithError(err)
	}

	if err := f.client.Delete(context.TODO(), &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Namespace: namespace,
			Name:      name,
		},
	}); err != nil && !errors.IsNotFound(err) {
		return reconcile.RequeueWithError(err)
	}

	return reconcile.Continue()
}

func (f *finalizer) ensureFinalizersRemoved() (reconcile.Result, error) {
	controllerutil.RemoveFinalizer(f.obj, utils.NamespaceFinalizer)
	if err := f.client.Update(context.TODO(), f.obj); err != nil {
		return reconcile.RequeueWithError(err)
	}

	return reconcile.Stop()
}

func (f *finalizer) EnsureFinalized() (reconcile.Result, error) {
	if !f.isObjBeingDeleted() {
		return reconcile.Continue()
	}

	return reconcile.RunSubRoutine([]reconcile.SubroutineFunc{
		f.ensureRouteRemoved,
		f.ensureFinalizersRemoved,
	})
}
