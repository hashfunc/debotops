package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type (
	ApplicationEnvironment corev1.EnvVar
	ApplicationPort        corev1.ContainerPort
	ApplicationResource    corev1.ResourceRequirements
)

type ApplicationHealth struct {
	Startup *corev1.Probe `json:"startup"`
	Ready   *corev1.Probe `json:"ready"`
	Live    *corev1.Probe `json:"live"`
}

type ApplicationContainer struct {
	Name         string                   `json:"name"`
	Image        string                   `json:"image"`
	Ports        []ApplicationPort        `json:"ports"`
	Environments []ApplicationEnvironment `json:"environments"`
	Command      []string                 `json:"command"`
	Args         []string                 `json:"args"`
	Replicas     int32                    `json:"replicas"`
	Resource     ApplicationResource      `json:"resource"`
	Health       ApplicationHealth        `json:"health"`
}

type ApplicationOptionProxy struct {
	Enable bool `json:"enable"`
}

type ApplicationOption struct {
	Proxy ApplicationOptionProxy `json:"proxy"`
}

// ApplicationSpec defines the desired state of Application
type ApplicationSpec struct {
	Container ApplicationContainer `json:"container"`
	Option    ApplicationOption    `json:"option"`
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
