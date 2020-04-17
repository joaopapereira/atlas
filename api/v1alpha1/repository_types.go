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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RepositorySpec defines the desired state of Repository
type RepositorySpec struct {
	Tag            string `json:"tag,omitempty"`
	ServiceAccount string `json:"serviceAccount,omitempty"`
}

const ImageUnreachable corev1.ConditionStatus = "ImageUnreachable"
const MetadataFormat corev1.ConditionStatus = "MetadataFormat"

// RepositoryStatus defines the observed state of Repository
type RepositoryStatus struct {
	Conditions         Conditions `json:"conditions,omitempty"`
	ObservedGeneration int64      `json:"observedGeneration,omitempty"`
	LatestImage        string     `json:"latestImage,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// Repository is the Schema for the repositories API
type Repository struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RepositorySpec   `json:"spec,omitempty"`
	Status RepositoryStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// RepositoryList contains a list of Repository
type RepositoryList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Repository `json:"items"`
}

func (in RepositoryStatus) Condition(condition ConditionType) corev1.ConditionStatus {
	for _, cond := range in.Conditions {
		if cond.Type == condition {
			return cond.Status
		}
	}
	return corev1.ConditionUnknown
}

func (in *RepositoryStatus) AddCondition(condition ConditionType, status corev1.ConditionStatus, message string) {
	in.Conditions = append(in.Conditions, Condition{
		Type:    condition,
		Status:  status,
		Message: message,
	})
}

func init() {
	SchemeBuilder.Register(&Repository{}, &RepositoryList{})
}
