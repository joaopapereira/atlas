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

// ProductReconciler reconciles a Product object
type ProductReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=atlas.jpereira.co.uk,resources=products,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.jpereira.co.uk,resources=products/status,verbs=get;update;patch

func (r *ProductReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("product", req.NamespacedName)

	// your logic here
	var product atlasv1alpha1.Product
	if err := r.Get(ctx, req.NamespacedName, &product); err != nil {
		log.Error(err, "unable to fetch CronJob")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	return ctrl.Result{}, nil
}

func (r *ProductReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&atlasv1alpha1.Product{}).
		Complete(r)
}
