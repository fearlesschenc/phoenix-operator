package initialize

import (
	"context"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
	corev1 "k8s.io/api/core/v1"
)

func (init *initializer) initializeFields(obj *corev1.Namespace) reconcile.ObjectState {
	objState := reconcile.ObjectUnchanged

	if obj.Labels == nil {
		obj.Labels = make(map[string]string)
		objState = reconcile.ObjectChanged
	}

	return objState
}

func (init *initializer) ensureFieldsInitialized() (reconcile.Result, error) {
	objState := init.initializeFields(init.obj)
	if objState == reconcile.ObjectUnchanged {
		return reconcile.Continue()
	}

	if err := init.client.Update(context.TODO(), init.obj); err != nil {
		return reconcile.RequeueWithError(err)
	}

	return reconcile.Stop()
}
