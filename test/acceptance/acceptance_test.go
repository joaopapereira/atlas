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
package acceptance

import (
	"testing"

	"github.com/google/go-containerregistry/pkg/name"
	v1 "github.com/google/go-containerregistry/pkg/v1"
	"github.com/google/go-containerregistry/pkg/v1/empty"
	"github.com/google/go-containerregistry/pkg/v1/mutate"
	"github.com/google/go-containerregistry/pkg/v1/remote"
	"github.com/sclevine/spec"
)

func createRepositoryImage() error {
	ref, err := name.ParseReference("registry.default.svc.cluster.local:5001/repository", name.WeakValidation)
	if err != nil {
		return err
	}
	image, err := emptyImage()
	if err != nil {
		return err
	}
	cfg, err := image.ConfigFile()
	if err != nil {
		return err
	}
	config := cfg.Config.DeepCopy()
	if config.Labels == nil {
		config.Labels = map[string]string{}
	}
	config.Labels["repository.atlas.jpereira.co.uk"] = `[{
	"name": "some product",
	"slug": "some-product",
    "version": {
		"major": "3",
		"minor": "2",
		"patch": "1"
    },
	"imageRepo": "registry.default.svc.cluster.local:5001/product-1",
    "imageSha": "sha256:03cbce912ef1a8a658f73c660ab9c539d67188622f00b15c4f15b89b884f0e10"
}]`

	image, err = mutate.Config(image, *config)

	return remote.Write(ref, image)
}

func emptyImage() (v1.Image, error) {
	cfg := &v1.ConfigFile{
		Architecture: "amd64",
		OS:           "linux",
		RootFS: v1.RootFS{
			Type:    "layers",
			DiffIDs: []v1.Hash{},
		},
	}
	return mutate.ConfigFile(empty.Image, cfg)
}

func createProductImage() error {
	ref, err := name.ParseReference("registry.default.svc.cluster.local:5001/product-1", name.WeakValidation)
	if err != nil {
		return err
	}

	image, err := emptyImage()
	if err != nil {
		return err
	}

	return remote.Write(ref, image)
}

func Test(t *testing.T) {
	spec.Run(t, "Test ", test)
}

func test(t *testing.T, when spec.G, it spec.S) {
	when("a", func() {
		it("b", func() {
			//require.NoError(t, createProductImage())
			//require.NoError(t, createRepositoryImage())
		})
	})
}
