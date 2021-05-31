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
package release

import (
	"context"
	"testing"

	"github.com/sclevine/spec"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes/scheme"
	fake2 "sigs.k8s.io/controller-runtime/pkg/client/fake"

	"github.com/joaopapereira/atlas/api/v1alpha1"
	"github.com/joaopapereira/atlas/pkg/repository"
)

func TestRepo(t *testing.T) {
	v1alpha1.AddToScheme(scheme.Scheme)
	spec.Run(t, "Test Repo", testRepo)
}

func testRepo(t *testing.T, when spec.G, it spec.S) {
	var (
		k8sClient = fake2.NewFakeClientWithScheme(scheme.Scheme)
		subject   = NewRepo(k8sClient)
	)

	const (
		namespace = "my-namespace"
	)

	when("#Get", func() {
		it.Before(func() {
			err := k8sClient.Create(context.TODO(), &v1alpha1.ProductRelease{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespace,
					Name:      "slug-1-123",
				},
				Spec: v1alpha1.ProductReleaseSpec{
					Slug: "slug-1",
					Version: v1alpha1.Version{
						Major: "1",
						Minor: "2",
						Patch: "3",
					},
					Image: "repo/image:tag",
				},
			})
			require.NoError(t, err)

		})

		it("can retrieve an existing version", func() {
			slug1V2 := v1alpha1.ProductRelease{
				ObjectMeta: metav1.ObjectMeta{
					Namespace: namespace,
					Name:      "slug-1-200",
				},
				Spec: v1alpha1.ProductReleaseSpec{
					Slug: "slug-1",
					Version: v1alpha1.Version{
						Major: "2",
						Minor: "0",
						Patch: "0",
					},
					Image: "repo/image-1:tag",
				},
			}
			require.NoError(t, subject.Save(slug1V2))

			rel2, err := subject.Get(namespace, "slug-1", repository.NewVersion("2.0.0"))
			require.NoError(t, err)

			require.Equal(t, "repo/image-1:tag", rel2.Spec.Image)
		})

		it("returns empty when cannot find the version", func() {
			rel2, err := subject.Get(namespace, "slug-1", repository.NewVersion("1.0.0"))
			require.NoError(t, err)

			require.Equal(t, v1alpha1.ProductRelease{}, rel2)
		})
	})
}
