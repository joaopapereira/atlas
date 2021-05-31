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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	atlasv1alpha1 "github.com/joaopapereira/atlas/api/v1alpha1"
	pkg "github.com/joaopapereira/atlas/pkg"
)

type releaseRepository interface {
	Save(release atlasv1alpha1.ProductRelease) error
	Get(namespace, slug string, version atlasv1alpha1.Version) (atlasv1alpha1.ProductRelease, error)
}

func NewUseCase(releaseRepo releaseRepository) *useCase {
	return &useCase{
		releaseRepo: releaseRepo,
	}
}

type useCase struct {
	releaseRepo releaseRepository
}

func (u useCase) Execute(product atlasv1alpha1.Product) (atlasv1alpha1.Product, error) {
	var updateVersions []atlasv1alpha1.ProductLocation
	var newVersions []atlasv1alpha1.ProductLocation
	for _, version := range product.Spec.Versions {
		exists := false
		for _, productLocation := range product.Status.Versions {
			if version.Version.Equals(productLocation.Version) {
				if !version.Equals(productLocation) {
					updateVersions = append(updateVersions, version)
					productLocation.ImageRepo = version.ImageRepo
					productLocation.ImageSHA = version.ImageSHA
				}
				exists = true
				break
			}
		}

		if !exists {
			newVersion := atlasv1alpha1.ProductLocation{
				ImageRepo: version.ImageRepo,
				ImageSHA:  version.ImageSHA,
				Version:   version.Version,
			}
			newVersions = append(newVersions, newVersion)
			product.Status.Versions = append(product.Status.Versions, newVersion)
		}
	}

	for _, version := range newVersions {
		u.releaseRepo.Save(atlasv1alpha1.ProductRelease{
			ObjectMeta: metav1.ObjectMeta{
				OwnerReferences: []metav1.OwnerReference{*pkg.RepoOwnerRef(&product)},
				Name: product.Name + "",
			},
			Spec:       atlasv1alpha1.ProductReleaseSpec{
				Slug:    product.Spec.Slug,
				Version: version.Version,
				Image:   version.ImageRepo + "@sha256:" + version.ImageSHA,
			},
		})
	}
	return product, nil
}
