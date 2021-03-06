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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// NetworkPolicySpec defines the desired state of NetworkPolicy
type NetworkPolicySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// +required
	Workspace string `json:"workspace"`

	// +optional
	NamespaceSelector metav1.LabelSelector `json:"namespaceSelector,omitempty"`

	// +required
	From []NetworkPolicyPeer `json:"from"`
}

type NetworkPolicyPeer struct {
	// +required
	Workspace string `json:"workspace"`

	// NamespaceSelector match namespace in networkpolicy's workspace. An empty
	// NamespaceSelector select all namespace in workspace.
	// +optional
	NamespaceSelector metav1.LabelSelector `json:"namespaceSelector,omitempty"`
}

type NetworkPolicyRef struct {
	// +required
	Namespace string `json:"namespace"`
}

type NetworkPolicyStatus struct {
	NetworkPolicyRefs []NetworkPolicyRef `json:"networkPolicyRefs"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:scope=Cluster

// NetworkPolicy is the Schema for the networkpolicies API
type NetworkPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NetworkPolicySpec   `json:"spec,omitempty"`
	Status NetworkPolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// NetworkPolicyList contains a list of NetworkPolicy
type NetworkPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NetworkPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NetworkPolicy{}, &NetworkPolicyList{})
}

const NetworkPolicyFinalizer = "networkpolicy.finalizer.kubesphere.io"
const NetworkPolicyOwnerKey = ".meta.controller"
