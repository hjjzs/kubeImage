package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ImageSpec defines the desired state of Image
type ImageSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS -- desired state of cluster
	ImageType      string `json:"imageType"`
	ImageUrl       string `json:"imageUrl"`
	RegisterSecret string `json:"registerSecret"`
	ImageTag       string `json:"imageTag"`
}

// ImageStatus defines the observed state of Image.
// It should always be reconstructable from the state of the cluster and/or outside world.
type ImageStatus struct {
	// INSERT ADDITIONAL STATUS FIELDS -- observed state of cluster
	ImageSize     string `json:"imageSize"`
	ImagePullPath string `json:"imagePullPath"`
	State         string `json:"state"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Image is the Schema for the images API
// +genclient
// +genclient:nonNamespaced
// +k8s:openapi-gen=true
type Image struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ImageSpec   `json:"spec,omitempty"`
	Status ImageStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ImageList contains a list of Image
type ImageList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Image `json:"items"`
}
