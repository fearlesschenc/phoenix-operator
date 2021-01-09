package validate

import (
	networkingv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/networking/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
)

type Validator interface {
	EnsureValidated() (reconcile.Result, error)
}

func NewValidator(obj *networkingv1alpha1.NetworkPolicy) Validator {
	return &validator{obj: obj}
}
