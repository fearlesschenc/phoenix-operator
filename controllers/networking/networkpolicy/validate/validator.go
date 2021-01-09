package validate

import (
	networkingv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/networking/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
)

type validator struct {
	obj *networkingv1alpha1.NetworkPolicy
}

func (v *validator) EnsureValidated() (reconcile.Result, error) {
	return reconcile.Continue()
}
