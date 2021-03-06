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
package controllers_test

import (
	"context"
	"testing"
	"time"

	"github.com/sclevine/spec"
	"github.com/stretchr/testify/require"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	"github.com/joaopapereira/atlas/api/v1alpha1"
	"github.com/joaopapereira/atlas/controllers"
	"github.com/joaopapereira/atlas/test"
)

var (
	cfg        *rest.Config
	k8sClient  client.Client
	testEnv    envtest.Environment
	k8sManager ctrl.Manager
)

func TestRepositoryController(t *testing.T) {
	cfg, k8sClient, testEnv, k8sManager = test.Before(t)
	defer test.After(t, testEnv)

	spec.Run(t, "Test Repository Controller", testRepositoryController)
}

type fakeRepositoryUseCase struct {
}

func (f fakeRepositoryUseCase) Execute(repository v1alpha1.Repository) (v1alpha1.Repository, error) {
	return repository, nil
}

func testRepositoryController(t *testing.T, when spec.G, it spec.S) {
	var (
		ctx = context.TODO()
	)
	const (
		namespace = "some-namespace"
	)

	it.Before(func() {
		require.NoError(t, (&controllers.RepositoryReconciler{
			Client:  k8sClient,
			Log:     ctrl.Log.WithName("controllers").WithName("Run"),
			Scheme:  nil,
			UseCase: fakeRepositoryUseCase{},
		}).SetupWithManager(k8sManager),
		)
	})

	when("test", func() {
		it("should work", func() {
			key := types.NamespacedName{
				Name:      "repo1",
				Namespace: namespace,
			}
			repository := &v1alpha1.Repository{
				ObjectMeta: metav1.ObjectMeta{
					Name:      key.Name,
					Namespace: key.Namespace,
				},
				Spec: v1alpha1.RepositorySpec{
					Tag:            "some.image/tag",
					ServiceAccount: "service-account",
				},
			}
			require.NoError(t, k8sClient.Create(ctx, repository))

			var resultRepository v1alpha1.Repository
			require.Eventually(t, func() bool {
				k8sClient.Get(ctx, key, &resultRepository)
				return resultRepository.Status.ObservedGeneration == 2
			}, 15*time.Second, 3*time.Second)
		})
	})
}
