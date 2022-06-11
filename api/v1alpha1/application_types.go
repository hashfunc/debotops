package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ApplicationHealth struct {
	Startup *corev1.Probe `json:"startup"`
	Ready   *corev1.Probe `json:"ready"`
	Live    *corev1.Probe `json:"live"`
}

type ApplicationContainer struct {
	Name         string                       `json:"name"`
	Image        string                       `json:"image"`
	Ports        []corev1.ContainerPort       `json:"ports,omitempty"`
	Environments []corev1.EnvVar              `json:"environments,omitempty"`
	Command      []string                     `json:"command,omitempty"`
	Args         []string                     `json:"args,omitempty"`
	Resource     *corev1.ResourceRequirements `json:"resource,omitempty"`
	Health       *ApplicationHealth           `json:"health,omitempty"`
}

type ApplicationOptionProxy struct {
	Enable bool `json:"enable"`
}

type ApplicationOption struct {
	Proxy *ApplicationOptionProxy `json:"proxy,omitempty"`
}

// ApplicationSpec defines the desired state of Application
type ApplicationSpec struct {
	Scale     int32                `json:"scale"`
	Container ApplicationContainer `json:"container"`
	Option    *ApplicationOption   `json:"option,omitempty"`
}

// ApplicationStatus defines the observed state of Application
type ApplicationStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Application is the Schema for the applications API
type Application struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApplicationSpec   `json:"spec,omitempty"`
	Status ApplicationStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ApplicationList contains a list of Application
type ApplicationList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`

	Items []Application `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Application{}, &ApplicationList{})
}
