package initialize

import (
	"context"
	networkingv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/networking/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
)

func (init *initializer) initializeFields(obj *networkingv1alpha1.NetworkPolicy) reconcile.ObjectState {
	objState := reconcile.ObjectUnchanged

	if init.obj.Status.NetworkPolicyRefs == nil {
		init.obj.Status.NetworkPolicyRefs = make([]networkingv1alpha1.NetworkPolicyRef, 0)
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
