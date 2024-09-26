package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// BuilderSpec defines the desired state of Builder
type BuilderSpec struct {
	DockerFileBase64 string        `json:"dockerFileBase64"`
	DockerFileString string        `json:"dockerFileString"`
	RemoteContext    RemoteContext `json:"remoteContext"`

	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10
	BuildTimeout int `json:"buildTimeout"`

	BuildName string `json:"buildName"`
}

// BuilderStatus defines the observed state of Builder.
// It should always be reconstructable from the state of the cluster and/or outside world.
type BuilderStatus struct {
	State string `json:"state"`
	// INSERT ADDITIONAL STATUS FIELDS -- observed state of cluster
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Builder is the Schema for the builders API
// +genclient
// +genclient:nonNamespaced
// +k8s:openapi-gen=true
// +kubebuilder:resource:scope=Cluster
// +kubebuilder:subresource:status
type Builder struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   BuilderSpec   `json:"spec,omitempty"`
	Status BuilderStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// BuilderList contains a list of Builder
type BuilderList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Builder `json:"items"`
}

type RemoteContext struct {

	// +kubebuilder:validation:MaxLength=200
	ContentUrl string `json:"contentUrl"`
	Type       string `json:"type"`

	// +kubebuilder:validation:MaxLength=20
	// +kubebuilder:validation:MinLength=1
	DockerFileName string `json:"dockerFileName"`

	// +kubebuilder:validation:MaxLength=20
	// +kubebuilder:validation:MinLength=1
	AuthConfigMap string `json:"authConfigMap"`
}
