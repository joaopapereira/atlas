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

// ProductReleaseSpec defines the desired state of ProductRelease
type ProductReleaseSpec struct {
	Slug    string  `json:"slug,omitempty"`
	Version Version `json:"version,omitempty"`
	Image   string  `json:"image,omitempty"`
}

// ProductReleaseStatus defines the observed state of ProductRelease
type ProductReleaseStatus struct {
	AtlasStatus `json:",inline"`
}

// +kubebuilder:object:root=true

// ProductRelease is the Schema for the productreleases API
type ProductRelease struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ProductReleaseSpec   `json:"spec,omitempty"`
	Status ProductReleaseStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ProductReleaseList contains a list of ProductRelease
type ProductReleaseList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ProductRelease `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ProductRelease{}, &ProductReleaseList{})
}
