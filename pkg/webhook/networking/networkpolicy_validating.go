package networking

import (
	"context"
	networkingv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/networking/v1alpha1"
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

	// TODO: validate workspace exist
	// TODO: validate From peer workspace exist

	return admission.Allowed("")
}

func (v *NetworkPolicyValidator) InjectDecoder(d *admission.Decoder) error {
	v.decoder = d
	return nil
}
