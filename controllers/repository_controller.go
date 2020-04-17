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

package controllers

import (
	"context"

	"github.com/go-logr/logr"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	atlasv1alpha1 "github.com/joaopapereira/atlas/api/v1alpha1"
)

type RepositoryUseCase interface {
	Execute(repository atlasv1alpha1.Repository) (atlasv1alpha1.Repository, error)
}

// RepositoryReconciler reconciles a Repository object
type RepositoryReconciler struct {
	client.Client
	Log     logr.Logger
	Scheme  *runtime.Scheme
	UseCase RepositoryUseCase
}

// +kubebuilder:rbac:groups=atlas.jpereira.co.uk,resources=repositories,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.jpereira.co.uk,resources=repositories/status,verbs=get;update;patch

func (r *RepositoryReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("repository", req.NamespacedName)

	var repository atlasv1alpha1.Repository
	err := r.Get(ctx, req.NamespacedName, &repository)
	if err != nil {
		log.Error(err, "unable to get repository")
		return ctrl.Result{}, err
	}
	if repository.Generation == repository.Status.ObservedGeneration {
		log.Info("no need to reconcile")
		return ctrl.Result{}, nil
	}

	repository, err = r.UseCase.Execute(repository)
	if err != nil {
		log.Error(err, "executing use case")
		return ctrl.Result{}, err
	}

	repository.Status.ObservedGeneration = repository.Generation + 1

	if err := r.Status().Update(ctx, &repository); err != nil {
		log.Error(err, "updating repository")
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

func (r *RepositoryReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&atlasv1alpha1.Repository{}).
		Complete(r)
}
