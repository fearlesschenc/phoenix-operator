package util

import (
	"context"
	networkingv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/networking/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type NetworkPolicyPredicator struct {
	Client client.Client
}

func (n *NetworkPolicyPredicator) Predicate(meta metav1.Object, object runtime.Object) bool {
	np := &networkingv1alpha1.NetworkPolicy{}
	if err := n.Client.Get(context.TODO(), types.NamespacedName{Name: meta.GetName()}, np); err != nil {
		return false
	}

	return true
}
