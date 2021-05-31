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
	"encoding/json"
	"log"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	atlasv1alpha1 "github.com/joaopapereira/atlas/api/v1alpha1"
	pkg "github.com/joaopapereira/atlas/pkg"
)

const RepositoryLabel = "repository.atlas.jpereira.co.uk"

type RegistryReader interface {
	GetImagesFromRepository(imageRepository string, option ...remote.Option) (string, string, error)
}

type ProductCreator interface {
	CreateOrUpdate(ref *metav1.OwnerReference, namespace string, imageJson ImageJSON) error
}

type ReleaseRepo interface {
	Save(release atlasv1alpha1.ProductRelease) error
}

func NewUseCaseWithReader(productCreator ProductCreator, releaseRepo ReleaseRepo, reader RegistryReader) *useCase {
	return &useCase{
		RepositoryReader: reader,
		ProductCreator:   productCreator,
		releaseRepo:      releaseRepo,
	}
}

func NewUseCase(productCreator ProductCreator, releaseRepo ReleaseRepo) *useCase {
	return &useCase{
		RepositoryReader: &repositoryReader{},
		ProductCreator:   productCreator,
		releaseRepo:      releaseRepo,
	}
}

type useCase struct {
	RepositoryReader RegistryReader
	ProductCreator   ProductCreator
	releaseRepo      ReleaseRepo
}

func (u useCase) Execute(repository atlasv1alpha1.Repository) (atlasv1alpha1.Repository, error) {
	repo := *repository.DeepCopy()
	repo.Status.Conditions = atlasv1alpha1.Conditions{}
	imageRef, images, err := u.RepositoryReader.GetImagesFromRepository(repo.Spec.Tag)
	if err != nil {
		repo.Status.AddCondition(atlasv1alpha1.ConditionSucceeded, atlasv1alpha1.ImageUnreachable, err.Error())
		return repo, nil
	}

	if imageRef == repo.Status.LatestImage {
		return repo, nil
	}

	var imagesJson ImagesJSON
	err = json.Unmarshal([]byte(images), &imagesJson)
	if err != nil {
		repo.Status.AddCondition(atlasv1alpha1.ConditionSucceeded, atlasv1alpha1.MetadataFormat, err.Error())
		return repo, nil
	}

	for _, imageJSON := range imagesJson {
		if err := u.ProductCreator.CreateOrUpdate(pkg.RepoOwnerRef(&repo), repo.Namespace, imageJSON); err != nil {
			return atlasv1alpha1.Repository{}, err
		}
		//
		//u.releaseRepo.Save(atlasv1alpha1.ProductRelease{
		//	ObjectMeta: metav1.ObjectMeta{
		//		Namespace:    repo.Namespace,
		//		GenerateName: imageJSON.Slug,
		//	},
		//	Spec: atlasv1alpha1.ProductReleaseSpec{
		//		Slug: imageJSON.Slug,
		//		Version: atlasv1alpha1.Version{
		//			Major: imageJSON.Version.Major(),
		//			Minor: imageJSON.Version.Minor(),
		//			Patch: imageJSON.Version.Patch(),
		//		},
		//		Image: fmt.Sprintf("%s@sha256:%s", imageJSON.ImageRepo, imageJSON.ImageSHA),
		//	},
		//})
	}

	repo.Status.AddCondition(atlasv1alpha1.ConditionSucceeded, corev1.ConditionTrue, "")
	repo.Status.LatestImage = imageRef

	return repo, nil
}

type repositoryReader struct {
}

func (r *repositoryReader) GetImagesFromRepository(imageRepository string, options ...remote.Option) (string, string, error) {
	ref, err := name.ParseReference(imageRepository, name.WeakValidation)
	if err != nil {
		return "", "", err
	}
	img, err := remote.Image(ref, options...) // TODO: missing service account creds
	if err != nil {
		return "", "", err
	}

	imgId, err := getIdentifier(img, ref)
	if err != nil {
		return "", "", err
	}
	log.Printf("repo: %v", imageRepository)
	images, err := getStringLabel(img, RepositoryLabel)

	if err != nil {
		return "", "", err
	}
	return imgId, images, nil
}

func getStringLabel(image v1.Image, key string) (string, error) {
	configFile, err := configFile(image)
	if err != nil {
		return "", err
	}

	config := configFile.Config.DeepCopy()
	log.Printf("all labels: %+v", config.Labels)
	stringValue, ok := config.Labels[key]
	if !ok {
		return "", errors.Errorf("could not find label %s", key)
	}

	return stringValue, nil
}

func getIdentifier(image v1.Image, ref name.Reference) (string, error) {
	digest, err := image.Digest()
	if err != nil {
		return "", errors.Wrapf(err, "failed to get digest for image '%s'", ref.Context().Name())
	}
	return ref.Context().Name() + "@" + digest.String(), nil
}

func configFile(image v1.Image) (*v1.ConfigFile, error) {
	cfg, err := image.ConfigFile()
	if err != nil {
		return nil, errors.Wrap(err, "failed to get image config")
	} else if cfg == nil {
		return nil, errors.Errorf("got nil image config")
	}
	return cfg, nil
}
