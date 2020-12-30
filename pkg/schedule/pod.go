package schedule

import (
	"github.com/fearlesschenc/kubesphere/pkg/constants"
	corev1 "k8s.io/api/core/v1"
)

func GetPodTolerationOf(workspace string) []corev1.Toleration {
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

func GetPodAffinityOf(workspace string) *corev1.Affinity {
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
