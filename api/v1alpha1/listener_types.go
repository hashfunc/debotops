package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ListenerGateway struct {
	Selector map[string]string `json:"selector"`
}

type ListenerTLS struct {
	Mode       string `json:"mode"`
	Credential string `json:"credential"`
}

type ListenerBind struct {
	Hosts    []string     `json:"hosts"`
	Name     string       `json:"name"`
	Port     uint32       `json:"port"`
	Protocol string       `json:"protocol"`
	TLS      *ListenerTLS `json:"tls,omitempty"`
}

// ListenerSpec defines the desired state of Listener
type ListenerSpec struct {
	Gateway ListenerGateway `json:"gateway"`
	Bind    []ListenerBind  `json:"bind"`
}

// ListenerStatus defines the observed state of Listener
type ListenerStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Listener is the Schema for the listeners API
type Listener struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ListenerSpec   `json:"spec,omitempty"`
	Status ListenerStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ListenerList contains a list of Listener
type ListenerList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Listener `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Listener{}, &ListenerList{})
}
