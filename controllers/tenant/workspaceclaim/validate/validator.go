package validate

import (
	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/reconcile"
)

type validator struct {
	obj *tenantv1alpha1.WorkspaceClaim
}

func (v *validator) EnsureWorkspaceClaimValidated() (reconcile.Result, error) {
	if v.obj.Name != v.obj.Spec.WorkspaceRef.Name {
		return reconcile.Stop()
	}

	return reconcile.Continue()
}
