package node

import (
	tenantv1alpha1 "github.com/fearlesschenc/phoenix-operator/apis/tenant/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	"reflect"
)

type OccupyStatus struct {
	Occupied bool
	Policy   tenantv1alpha1.OccupyPolicy
}

type OccupyOperation struct {
	Cluster      string
	Node         *corev1.Node
	ExpectStatus *OccupyStatus
}

type Occupier interface {
	Occupy(cluster string, node *corev1.Node) error
	DeOccupy(cluster string, node *corev1.Node) error
}

func NewOccupier(policy tenantv1alpha1.OccupyPolicy) Occupier {
	switch policy {
	case tenantv1alpha1.Exclusive:
		return &ExclusiveOccupier{}
	}

	return &ExclusiveOccupier{}
}

func (op *OccupyOperation) ShouldSkipOn(status *OccupyStatus) bool {
	if status.Occupied == false && op.ExpectStatus.Occupied == false {
		return true
	}

	return reflect.DeepEqual(status, op.ExpectStatus)
}

func (op *OccupyOperation) OperateOn(currentStatus *OccupyStatus) error {
	currentOccupier := NewOccupier(currentStatus.Policy)
	expectOccupier := NewOccupier(op.ExpectStatus.Policy)

	if !currentStatus.Occupied {
		// expect status must be true or this operation has been skipped
		return expectOccupier.Occupy(op.Cluster, op.Node)
	}

	// then, left with current status is occupied.
	if !op.ExpectStatus.Occupied {
		return currentOccupier.DeOccupy(op.Cluster, op.Node)
	}

	var err error
	err = currentOccupier.DeOccupy(op.Cluster, op.Node)
	if err != nil {
		return err
	}
	err = expectOccupier.Occupy(op.Cluster, op.Node)
	if err != nil {
		return err
	}

	return nil
}
