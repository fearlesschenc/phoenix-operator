package util

import (
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
)

func HandleNetworkPolicy(object handler.MapObject) []reconcile.Request {
	return []reconcile.Request{{types.NamespacedName{Name: object.Meta.GetName()}}}
}
