package validate

import (
	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
)

type Validator interface {
	EnsureWorkspaceClaimValidated() (reconcile.Result, error)
}

func NewValidator(obj *tenantv1alpha1.WorkspaceClaim) Validator {
	return &validator{obj: obj}
}
