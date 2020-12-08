package webhook

import (
	"context"
	"fmt"
	"github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	"net"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:path=/validate-tenant-v1alpha1-workspace,mutating=false,failurePolicy=fail,groups=""
type WorkspaceValidator struct {
	client.Client
	decoder *admission.Decoder
}

func (v *WorkspaceValidator) validateHosts(ctx context.Context, hosts []string) error {
	for _, host := range hosts {
		if ip := net.ParseIP(host); ip == nil || ip.To4() == nil {
			return fmt.Errorf("invalid ip address")
		}
	}

	return nil
}

func (v *WorkspaceValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	workspace := &v1alpha1.Workspace{}

	err := v.decoder.Decode(req, workspace)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	if err := v.validateHosts(ctx, workspace.Spec.Hosts); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	return admission.Allowed("")
}

func (v *WorkspaceValidator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}
