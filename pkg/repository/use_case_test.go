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
package repository

import (
	"errors"
	"testing"

	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/sclevine/spec"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	atlasv1alpha1 "github.com/joaopapereira/atlas/api/v1alpha1"
)

func TestUseCase(t *testing.T) {
	spec.Run(t, "Test UseCase", testUseCase)
}

func testUseCase(t *testing.T, when spec.G, it spec.S) {
	var (
		readerFake         = &repositoryReaderFake{}
		productCreatorFake = &productCreatorFake{}
		subject            = NewUseCaseWithReader(productCreatorFake, readerFake)
	)

	when("#Execute", func() {
		when("cannot read image", func() {
			it.Before(func() {
				readerFake.returnError = errors.New("unable to find image")
			})

			it("returns status ImageUnreachable", func() {
				repo := atlasv1alpha1.Repository{
					Spec: atlasv1alpha1.RepositorySpec{
						Tag:            "some-image/place:tag",
						ServiceAccount: "service-account",
					},
					Status: atlasv1alpha1.RepositoryStatus{
						LatestImage: "some-image/place:123abc",
					},
				}
				result, err := subject.Execute(repo)
				require.NoError(t, err)
				require.Len(t, result.Status.Conditions, 1)
				assert.Equal(t, atlasv1alpha1.ImageUnreachable, result.Status.Condition(atlasv1alpha1.ConditionSucceeded))
				assert.Equal(t, "some-image/place:123abc", result.Status.LatestImage)
				assert.Equal(t, 0, productCreatorFake.numberOfCalls)
			})
		})

		it("returns status MetadataFormat when cannot unmarshal image metadata", func() {
			readerFake.returnLabel = "{asdasd}"
			repo := atlasv1alpha1.Repository{
				Spec: atlasv1alpha1.RepositorySpec{
					Tag:            "some-image/place:tag",
					ServiceAccount: "service-account",
				},
				Status: atlasv1alpha1.RepositoryStatus{
					LatestImage: "some-image/place:123abc",
				},
			}
			result, err := subject.Execute(repo)
			require.NoError(t, err)
			require.Len(t, result.Status.Conditions, 1)
			assert.Equal(t, atlasv1alpha1.MetadataFormat, result.Status.Condition(atlasv1alpha1.ConditionSucceeded))
			assert.Equal(t, "some-image/place:123abc", result.Status.LatestImage)
			assert.Equal(t, 0, productCreatorFake.numberOfCalls)
		})

		it("returns status ConditionTrue and updates repository image sha when new repo image is present", func() {
			readerFake.returnLabel = `[{
	"name": "some product",
	"slug": "some-product",
    "version": {
		"major": "3",
		"minor": "2",
		"patch": "1"
    },
	"imageRepo": "some-other/image",
    "imageSha": "de434baf69e0d821eba847e1187e68dff4c27bdf99caf1e3417e6ab50e8533a7"
}]`
			readerFake.returnImageRef = "some-image/place:1231232523452345"
			repo := atlasv1alpha1.Repository{
				Spec: atlasv1alpha1.RepositorySpec{
					Tag:            "some-image/place:tag",
					ServiceAccount: "service-account",
				},
				Status: atlasv1alpha1.RepositoryStatus{
					LatestImage: "some-image/place:123abc",
				},
			}
			result, err := subject.Execute(repo)
			require.NoError(t, err)
			require.Len(t, result.Status.Conditions, 1)
			assert.Equal(t, v1.ConditionTrue, result.Status.Condition(atlasv1alpha1.ConditionSucceeded))
			assert.Equal(t, "some-image/place:1231232523452345", result.Status.LatestImage)
			require.Equal(t, 1, productCreatorFake.numberOfCalls)
			assert.Equal(t, ImageJSON{
				Product: Product{
					Name: "some product",
					Slug: "some-product",
					Version: Version{
						Major: "3",
						Minor: "2",
						Patch: "1",
					},
				},
				ImageRepo: "some-other/image",
				ImageSHA:  "de434baf69e0d821eba847e1187e68dff4c27bdf99caf1e3417e6ab50e8533a7",
			}, productCreatorFake.calledWith[0])
		})
	})
}

type repositoryReaderFake struct {
	returnImageRef string
	returnLabel    string
	returnError    error
}

func (r repositoryReaderFake) GetImagesFromRepository(_ string, _ ...remote.Option) (string, string, error) {
	return r.returnImageRef, r.returnLabel, r.returnError
}

type productCreatorFake struct {
	returnError   error
	calledWith    []ImageJSON
	numberOfCalls int
}

func (r *productCreatorFake) CreateOrUpdate(_ *metav1.OwnerReference, _ string, imageJson ImageJSON) error {
	r.numberOfCalls++
	r.calledWith = append(r.calledWith, imageJson)
	return r.returnError
}
