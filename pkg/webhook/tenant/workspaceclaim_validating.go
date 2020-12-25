package tenant

import (
	"context"
	"fmt"
	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/schedule"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// * validate Node specified
//   * If node exist
//   * if node already occupied

type WorkspaceClaimValidator struct {
	Client  client.Client
	decoder *admission.Decoder
}

// +kubebuilder:rbac:groups="",resources=nodes,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:webhook:path=/mutate-tenant-phoenix-fearlesschenc-com-v1alpha1-workspaceclaim,mutating=false,failurePolicy=fail,groups=tenant.phoenix.fearlesschenc.com,versions=v1alpha1,resources=workspaceclaims,verbs=create,name=vworkspaceclaim.kubesphere.io

func (v *WorkspaceClaimValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	claim := &tenantv1alpha1.WorkspaceClaim{}
	if err := v.decoder.Decode(req, claim); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	// TODO: validate workspaceRef exist

	if claim.Name != claim.Spec.WorkspaceRef.Name {
		return admission.Errored(http.StatusBadRequest, fmt.Errorf("claim's name must be same with workspace"))
	}

	nodes := claim.Spec.Node
	if nodes != nil {
		for _, nodeName := range nodes {
			node := &v1.Node{}

			if err := v.Client.Get(ctx, types.NamespacedName{Name: nodeName}, node); err != nil {
				if errors.IsNotFound(err) {
					return admission.Errored(http.StatusBadRequest, err)
				}

				return admission.Errored(http.StatusInternalServerError, err)
			}

			// have been occupied
			if workspace := schedule.GetNodeWorkspace(node); workspace != "" && workspace != claim.Spec.WorkspaceRef.Name {
				return admission.Errored(http.StatusBadRequest, fmt.Errorf("node have been occupied"))
			}
		}
	}

	return admission.Allowed("validated")
}

func (v *WorkspaceClaimValidator) InjectDecoder(decoder *admission.Decoder) error {
	v.decoder = decoder
	return nil
}
