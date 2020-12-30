package core

import (
	"context"
	"encoding/json"
	"github.com/fearlesschenc/kubesphere/pkg/constants"
	"github.com/fearlesschenc/phoenix-operator/pkg/schedule"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

type PodMutator struct {
	Client  client.Client
	decoder *admission.Decoder
}

// +kubebuilder:rbac:groups="",resources=namespaces,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:webhook:path=/mutate-v1-pod,mutating=true,failurePolicy=Ignore,groups="",versions=v1,resources=pods,verbs=create,name=mpod.kubesphere.io

func (m *PodMutator) Handle(ctx context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}
	if err := m.decoder.Decode(req, pod); err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	// get pod's namespace
	namespace := &corev1.Namespace{}
	if err := m.Client.Get(ctx, types.NamespacedName{Name: req.Namespace}, namespace); err != nil {
		return admission.Errored(http.StatusBadGateway, err)
	}

	// get pod's workspace
	workspace, ok := namespace.Labels[constants.WorkspaceLabelKey]
	if !ok {
		return admission.Allowed("skip: pod have no workspace")
	}

	// TODO:
	// * check affinity existence, affinity merge
	// * check toleration existence, toleration merge
	pod.Spec.Affinity = schedule.GetPodAffinityOf(workspace)
	pod.Spec.Tolerations = append(pod.Spec.Tolerations, schedule.GetPodTolerationOf(workspace)...)
	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

func (m *PodMutator) InjectDecoder(d *admission.Decoder) error {
	m.decoder = d
	return nil
}
