/*
Copyright The Kubernetes Authors.

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

// Code generated by client-gen. DO NOT EDIT.

package v1

import (
	v1 "builder/pkg/apis/image/v1"
	scheme "builder/pkg/client/generated/clientset/versioned/scheme"
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	types "k8s.io/apimachinery/pkg/types"
	watch "k8s.io/apimachinery/pkg/watch"
	gentype "k8s.io/client-go/gentype"
)

// ImagesGetter has a method to return a ImageInterface.
// A group's client should implement this interface.
type ImagesGetter interface {
	Images() ImageInterface
}

// ImageInterface has methods to work with Image resources.
type ImageInterface interface {
	Create(ctx context.Context, image *v1.Image, opts metav1.CreateOptions) (*v1.Image, error)
	Update(ctx context.Context, image *v1.Image, opts metav1.UpdateOptions) (*v1.Image, error)
	// Add a +genclient:noStatus comment above the type to avoid generating UpdateStatus().
	UpdateStatus(ctx context.Context, image *v1.Image, opts metav1.UpdateOptions) (*v1.Image, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Image, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.ImageList, error)
	Watch(ctx context.Context, opts metav1.ListOptions) (watch.Interface, error)
	Patch(ctx context.Context, name string, pt types.PatchType, data []byte, opts metav1.PatchOptions, subresources ...string) (result *v1.Image, err error)
	ImageExpansion
}

// images implements ImageInterface
type images struct {
	*gentype.ClientWithList[*v1.Image, *v1.ImageList]
}

// newImages returns a Images
func newImages(c *ImageV1Client) *images {
	return &images{
		gentype.NewClientWithList[*v1.Image, *v1.ImageList](
			"images",
			c.RESTClient(),
			scheme.ParameterCodec,
			"",
			func() *v1.Image { return &v1.Image{} },
			func() *v1.ImageList { return &v1.ImageList{} }),
	}
}
