package validate

import (
	networkingv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/networking/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/constants"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
)

type validator struct {
	obj *networkingv1alpha1.NetworkPolicy
}

func (v *validator) EnsureNetworkPolicyValidated() (reconcile.Result, error) {
	_, ok := v.obj.GetLabels()[constants.WorkspaceLabelKey]
	if !ok {
		return reconcile.Stop()
	}

	return reconcile.Continue()
}
