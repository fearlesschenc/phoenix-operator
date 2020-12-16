package webhook

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

type podAnnotator struct {
	Client  client.Client
	decoder *admission.Decoder
}

func (a *podAnnotator) MutatePod(pod *corev1.Pod, workspace string) {
	//pod.Spec.NodeSelector[constants.WorkspaceLabel] = workspace
	pod.Spec.Affinity.NodeAffinity = &corev1.NodeAffinity{
		PreferredDuringSchedulingIgnoredDuringExecution: []corev1.PreferredSchedulingTerm{
			{
				Weight: 100,
				Preference: corev1.NodeSelectorTerm{
					MatchExpressions: []corev1.NodeSelectorRequirement{
						{
							Key:      constants.WorkspaceLabel,
							Operator: corev1.NodeSelectorOpIn,
							Values:   []string{workspace},
						},
					},
				},
			},
		},
	}
	pod.Spec.Tolerations = append(pod.Spec.Tolerations, corev1.Toleration{
		Key:      constants.WorkspaceLabel,
		Operator: corev1.TolerationOpEqual,
		Value:    workspace,
		Effect:   corev1.TaintEffectNoExecute,
	}, corev1.Toleration{
		Key:      constants.WorkspaceLabel,
		Operator: corev1.TolerationOpEqual,
		Value:    workspace,
		Effect:   corev1.TaintEffectNoSchedule,
	})
}

func (a *podAnnotator) Handle(ctx context.Context, req admission.Request) admission.Response {
	pod := &corev1.Pod{}
	err := a.decoder.Decode(req, pod)
	if err != nil {
		return admission.Errored(http.StatusBadRequest, err)
	}

	namespace := &corev1.Namespace{}
	if err := a.Client.Get(ctx, types.NamespacedName{Name: pod.Namespace}, namespace); err != nil {
		return admission.Errored(http.StatusBadGateway, err)
	}

	workspace, ok := namespace.Labels[constants.WorkspaceLabel]
	if !ok {
		return admission.Allowed("ok")
	}

	a.MutatePod(pod, workspace)
	marshaledPod, err := json.Marshal(pod)
	if err != nil {
		return admission.Errored(http.StatusInternalServerError, err)
	}

	return admission.PatchResponseFromRaw(req.Object.Raw, marshaledPod)
}

func (a *podAnnotator) InjectDecoder(d *admission.Decoder) error {
	a.decoder = d
	return nil
}
