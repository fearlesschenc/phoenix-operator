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

type WorkspaceRef struct {
	Name string `json:"name"`
}

// WorkspaceClaimSpec defines the desired state of WorkspaceClaim
type WorkspaceClaimSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// +required
	WorkspaceRef WorkspaceRef `json:"workspaceRef"`

	// +optional
	Node []string `json:"node,omitempty"`
}

// WorkspaceClaimStatus defines the observed state of WorkspaceClaim
type WorkspaceClaimStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Node []string `json:"node"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:shortName=wksc,singular=workspaceclaim,scope=Cluster

// WorkspaceClaim is the Schema for the workspaceclaims API
type WorkspaceClaim struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   WorkspaceClaimSpec   `json:"spec,omitempty"`
	Status WorkspaceClaimStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// WorkspaceClaimList contains a list of WorkspaceClaim
type WorkspaceClaimList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []WorkspaceClaim `json:"items"`
}

func init() {
	SchemeBuilder.Register(&WorkspaceClaim{}, &WorkspaceClaimList{})
}

const WorkspaceClaimOwnerKey = ".meta.controller"
