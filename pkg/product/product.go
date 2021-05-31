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
package product

import (
	"context"

	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/joaopapereira/atlas/api/v1alpha1"
	"github.com/joaopapereira/atlas/pkg/repository"
)

func NewRepo(k8sClient client.Client) *productRepo {
	return &productRepo{k8sClient: k8sClient}
}

type productRepo struct {
	k8sClient client.Client
}

func (p *productRepo) CreateOrUpdate(ref *metav1.OwnerReference, namespace string, imageJson repository.ImageJSON) error {
	ctx := context.Background()
	key := client.ObjectKey{
		Namespace: namespace,
		Name:      imageJson.Slug,
	}
	var product v1alpha1.Product
	err := p.k8sClient.Get(ctx, key, &product)
	if errors.IsNotFound(err) {
		product = v1alpha1.Product{
			ObjectMeta: metav1.ObjectMeta{
				Name:            key.Name,
				Namespace:       key.Namespace,
				OwnerReferences: []metav1.OwnerReference{*ref},
			},
			Spec: v1alpha1.ProductSpec{
				Name:        imageJson.Name,
				Slug:        imageJson.Slug,
				Description: imageJson.Description,
			},
		}
	} else if err != nil {
		return err
	}

	versionExists := false
	for _, version := range product.Status.Versions {
		if compareVersions(version.Version, imageJson.Version) {
			versionExists = true
		}
	}

	if versionExists {
		return nil
	}
	product.Spec.Versions = append(product.Status.Versions, v1alpha1.ProductLocation{
		ImageRepo: imageJson.ImageRepo,
		ImageSHA:  imageJson.ImageSHA,
		Version: v1alpha1.Version{
			Major: imageJson.Version.Major(),
			Minor: imageJson.Version.Minor(),
			Patch: imageJson.Version.Patch(),
		},
	})

	return p.createOrUpdate(ctx, product)
}

func (p *productRepo) createOrUpdate(ctx context.Context, product v1alpha1.Product) error {
	a := metav1.Time{}
	if product.CreationTimestamp == a {
		return p.k8sClient.Create(ctx, &product)
	}
	return p.k8sClient.Update(ctx, &product)
}

func compareVersions(v1 v1alpha1.Version, v2 repository.Version) bool {
	return v1.Major == v2.Major() && v1.Minor == v2.Minor() && v1.Patch == v2.Patch()
}
