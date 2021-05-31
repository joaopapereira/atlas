package release

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/joaopapereira/atlas/api/v1alpha1"
	"github.com/joaopapereira/atlas/pkg/repository"
)

func NewRepo(k8sClient client.Client) *releaseRepo {
	return &releaseRepo{k8sClient: k8sClient}
}

type releaseRepo struct {
	k8sClient client.Client
}

func (r *releaseRepo) Get(namespace, slug string, version repository.Version) (v1alpha1.ProductRelease, error) {
	ctx := context.Background()

	var releases v1alpha1.ProductReleaseList
	err := r.k8sClient.List(ctx, &releases, client.InNamespace(namespace), client.MatchingFields{slug: slug})
	if err != nil {
		return v1alpha1.ProductRelease{}, err
	}

	for _, release := range releases.Items {
		if release.CompareVersion(version) {
			return release, nil
		}
	}

	return v1alpha1.ProductRelease{}, nil
}

func (r *releaseRepo) Save(release v1alpha1.ProductRelease) error {
	ctx := context.Background()

	t := metav1.Time{}
	if release.CreationTimestamp == t {
		return r.k8sClient.Create(ctx, &release)
	}
	return r.k8sClient.Update(ctx, &release)
}

func (r *releaseRepo) SaveStatus(release v1alpha1.ProductRelease) error {
	ctx := context.Background()

	return r.k8sClient.Status().Update(ctx, &release)
}
