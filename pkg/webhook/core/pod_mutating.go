package core

import (
	"context"
	"encoding/json"
	"github.com/fearlesschenc/phoenix-operator/pkg/constants"
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

func getPodAffinity(workspace string) *corev1.Affinity {
	return &corev1.Affinity{
		NodeAffinity: &corev1.NodeAffinity{
			PreferredDuringSchedulingIgnoredDuringExecution: []corev1.PreferredSchedulingTerm{
				{
					Weight: 100,
					Preference: corev1.NodeSelectorTerm{
						MatchExpressions: []corev1.NodeSelectorRequirement{
							{
								Key:      constants.WorkspaceLabelKey,
								Operator: corev1.NodeSelectorOpIn,
								Values:   []string{workspace},
							},
						},
					},
				},
			},
		},
	}
}

func getWorkspaceToleration(workspace string) []corev1.Toleration {
	return []corev1.Toleration{
		{
			Key:      constants.WorkspaceLabelKey,
			Operator: corev1.TolerationOpEqual,
			Value:    workspace,
			Effect:   corev1.TaintEffectNoExecute,
		},
		{
			Key:      constants.WorkspaceLabelKey,
			Operator: corev1.TolerationOpEqual,
			Value:    workspace,
			Effect:   corev1.TaintEffectNoSchedule,
		},
	}
}

func (m *PodMutator) InjectPodAffinityAndTolerations(pod *corev1.Pod, workspace string) {
	// pod.Spec.NodeSelector[constants.WorkspaceLabelKey] = workspace
	// prefer to schedule on workspace claim node.
	pod.Spec.Affinity = getPodAffinity(workspace)
	// TODO: toleration merge
	pod.Spec.Tolerations = append(pod.Spec.Tolerations, getWorkspaceToleration(workspace)...)
}

//func Injected(pod *corev1.Pod) bool {
//	// TODO: check affinity existence
//	// TODO: check toleration existence
//	//for _, toleration := range getWorkspaceToleration("") {
//	//	for _, tlr := range pod.Spec.Tolerations {
//	//		tolerations.AreEqual()
//	//	}
//	//}

//	return false
//}

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

	m.InjectPodAffinityAndTolerations(pod, workspace)
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
