package test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/envtest"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	atlasv1alpha1 "github.com/joaopapereira/atlas/api/v1alpha1"
)

func Before(t *testing.T) (*rest.Config, client.Client, envtest.Environment, ctrl.Manager) {
	logf.SetLogger(zap.LoggerTo(os.Stdout, true))
	testEnv := envtest.Environment{
		CRDDirectoryPaths: []string{filepath.Join("..", "config", "crd", "bases")},
	}

	cfg, err := testEnv.Start()
	require.NoError(t, err)
	require.NotNil(t, cfg)

	err = atlasv1alpha1.AddToScheme(scheme.Scheme)
	require.NoError(t, err)

	k8sManager, err := ctrl.NewManager(cfg, ctrl.Options{
		Scheme: scheme.Scheme,
	})
	require.NoError(t, err)

	// +kubebuilder:scaffold:scheme

	go func() {
		require.NoError(t, k8sManager.Start(ctrl.SetupSignalHandler()))
	}()

	k8sClient := k8sManager.GetClient()
	require.NotNil(t, k8sClient)
	return cfg, k8sClient, testEnv, k8sManager
}

func After(t *testing.T, testEnv envtest.Environment) {
	err := testEnv.Stop()
	require.NoError(t, err)
}
