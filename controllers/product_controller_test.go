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
	ctrl "sigs.k8s.io/controller-runtime"

	"github.com/joaopapereira/atlas/api/v1alpha1"
	"github.com/joaopapereira/atlas/controllers"
	"github.com/joaopapereira/atlas/test"
)

func TestProductController(t *testing.T) {
	cfg, k8sClient, testEnv, k8sManager = test.Before(t)
	defer test.After(t, testEnv)

	spec.Run(t, "Test Product Controller", testProductController)
}

func testProductController(t *testing.T, when spec.G, it spec.S) {
	var (
		ctx = context.TODO()
	)
	const (
		namespace = "some-namespace"
	)

	it.Before(func() {
		require.NoError(t, (&controllers.ProductReconciler{
			Client:  k8sClient,
			Log:     ctrl.Log.WithName("controllers").WithName("Run"),
			Scheme:  nil,
			UseCase: fakeProductUseCase{},
		}).SetupWithManager(k8sManager),
		)
	})

	when("starts", func() {
		it("should process the product", func() {
			key := types.NamespacedName{
				Name:      "product1",
				Namespace: namespace,
			}
			repository := &v1alpha1.Product{
				ObjectMeta: metav1.ObjectMeta{
					Name:      key.Name,
					Namespace: key.Namespace,
				},
				Spec: v1alpha1.ProductSpec{
					Slug: "some-product",
					Name: "Some Product",
				},
			}
			require.NoError(t, k8sClient.Create(ctx, repository))

			var resultProduct v1alpha1.Product
			require.Eventually(t, func() bool {
				k8sClient.Get(ctx, key, &resultProduct)
				t.Logf("obs == %d && %d", resultProduct.Generation, resultProduct.Status.ObservedGeneration)
				return resultProduct.Status.ObservedGeneration == 2
			}, 15*time.Second, 3*time.Second)
		})
	})
}

type fakeProductUseCase struct {
}

func (f fakeProductUseCase) Execute(product v1alpha1.Product) (v1alpha1.Product, error) {
	return product, nil
}
