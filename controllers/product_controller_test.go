package controllers_test

import (
	"testing"

	"github.com/sclevine/spec"
	"github.com/stretchr/testify/require"
	"sigs.k8s.io/controller-runtime/pkg/envtest"

	"github.com/joaopapereira/atlas/test"
)

func TestProductController(t *testing.T) {
	spec.Run(t, "Test Describe Image", testProductController)
}

func testProductController(t *testing.T, when spec.G, it spec.S) {
	var (
		//cfg       *rest.Config
		//k8sClient client.Client
		testEnv envtest.Environment
	)

	it.Before(func() {
		_, _, testEnv, _ = test.Before(t)
	})

	it.After(func() {
		test.After(t, testEnv)
	})

	when("test", func() {
		it("should work", func() {
			require.Equal(t, true, false)
		})
	})
}
