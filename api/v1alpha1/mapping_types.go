package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type MappingListener struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type MappingApplication struct {
	Name string `json:"name"`
	Port uint32 `json:"port"`
}

// MappingSpec defines the desired state of Mapping
type MappingSpec struct {
	Hosts       []string           `json:"hosts"`
	Listener    MappingListener    `json:"listener"`
	Application MappingApplication `json:"application"`
}

// MappingStatus defines the observed state of Mapping
type MappingStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Mapping is the Schema for the mappings API
type Mapping struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   MappingSpec   `json:"spec,omitempty"`
	Status MappingStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// MappingList contains a list of Mapping
type MappingList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Mapping `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Mapping{}, &MappingList{})
}
