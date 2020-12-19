package networking

import (
	"context"
	"fmt"
	networkingv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/networking/v1alpha1"
	"github.com/fearlesschenc/phoenix-operator/pkg/constants"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type NetworkPolicyValidator struct {
	Client  client.Client
	decoder *admission.Decoder
}

// +kubebuilder:webhook:path=/validate-tenant-phoenix-fearlesschenc-com-v1alpha1-networkpolicy,mutating=false,failurePolicy=fail,groups=networking.phoenix.fearlesschenc.com,versions=v1alpha1,resources=networkpolicys,verbs=create,name=vnetworkpolicy.kubesphere.io

func (v *NetworkPolicyValidator) Handle(ctx context.Context, req admission.Request) admission.Response {
	np := &networkingv1alpha1.NetworkPolicy{}
	v.decoder.Decode(req, np)

	workspace, ok := np.Labels[constants.WorkspaceLabel]
	if !ok {
		return admission.Errored(http.StatusBadRequest, fmt.Errorf("network policy must have a workspace"))
	}

	// TODO: validate workspace exist
	_ = ctx
	_ = workspace
	// TODO: validate Ingress peer workspace exist

	return admission.Allowed("")
}

func (v *NetworkPolicyValidator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}
