/*
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

// ProductSpec defines the desired state of Product
type ProductSpec struct {
	// +kubebuilder:validation:Required
	Name string `json:"name,omitempty"`
	// +kubebuilder:validation:Required
	Slug                 string            `json:"slug,omitempty"`
	Import               bool              `json:"import,omitempty"`
	BaseRegistryLocation string            `json:"baseRegistryLocation,omitempty"`
	Description          string            `json:"description,omitempty"`
	Versions             []ProductLocation `json:"versions,omitempty"`
}

type Version struct {
	Major string `json:"major,omitempty"`
	Minor string `json:"minor,omitempty"`
	Patch string `json:"patch,omitempty"`
}

// ProductStatus defines the observed state of Product
type ProductStatus struct {
	AtlasStatus `json:",inline"`
	Versions    []ProductLocation `json:"versions,omitempty"`
}

type ProductLocation struct {
	// +kubebuilder:validation:Required
	ImageRepo string `json:"imageRepo,omitempty"`
	// +kubebuilder:validation:Required
	ImageSHA string `json:"imageSha,omitempty"`
	// +kubebuilder:validation:Required
	Version Version `json:"version,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// Product is the Schema for the products API
type Product struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProductSpec   `json:"spec,omitempty"`
	Status ProductStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ProductList contains a list of Product
type ProductList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Product `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Product{}, &ProductList{})
}
