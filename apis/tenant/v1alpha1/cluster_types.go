/*
Authored by fearlesschenc@gmail.com

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	"github.com/fearlesschenc/phoenix-operator/apis/workload/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	ClusterReady           = "ClusterReady"
	ClusterNetworkIsolated = "ClusterNetworkIsolated"
	ClusterProvisioned     = "ClusterProvisioned"
)

const ClusterLabel = "cluster.phoenix.fearlesschenc.com"

type OccupyPolicy string

const (
	None      OccupyPolicy = "None"
	Exclusive OccupyPolicy = "exclusive"
)

type NodeOccupy struct {
	// +required
	NodeName string `json:"host"`

	// +optional
	Policy OccupyPolicy `json:"policy"`
}

// ClusterSpec defines the desired state of Cluster
type ClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// NetworkIsolation identifies whether the cluster should be isolated
	// from other cluster
	// +optional
	NetworkIsolation *bool `json:"networkIsolation,omitempty"`

	// NodeOccupies describe the node that this cluster claim to occupy.
	// +optional
	NodeOccupies []NodeOccupy `json:"nodeOccupies,omitempty" yaml:"nodeOccupies,omitempty"`
}

// ClusterStatus defines the observed state of Cluster
type ClusterStatus struct {
	// +optional
	//Condition []metav1.Condition `json:"condition"`

	Applications v1alpha1.ApplicationReference `json:"applications"`

	// +optional
	NodeOccupied []NodeOccupy `json:"occupiedNodes"`

	// +optional
	NetworkIsolated bool `json:"networkIsolated"`
}

// Cluster metaphor a set of machines that join the
// kubernetes cluster, on which developer will deploy
// bunch of applications on it.
// TODO: printcolumn not working, value empty
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:printcolumn:name="Isolated",type=boolean,JSONPath=".status.NetworkIsolated"
// +kubebuilder:printcolumn:name="NodeOccupied",type=string,JSONPath=`.status.NodeOccupied[*].NodeName`

// Cluster is the Schema for the clusters API
type Cluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterSpec   `json:"spec,omitempty"`
	Status ClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterList contains a list of Cluster
type ClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Cluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Cluster{}, &ClusterList{})
}
