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

// ProductReleaseReconciler reconciles a ProductRelease object
type ProductReleaseReconciler struct {
	client.Client
	Log    logr.Logger
	Scheme *runtime.Scheme
}

// +kubebuilder:rbac:groups=atlas.jpereira.co.uk,resources=productreleases,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=atlas.jpereira.co.uk,resources=productreleases/status,verbs=get;update;patch

func (r *ProductReleaseReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	_ = context.Background()
	_ = r.Log.WithValues("productrelease", req.NamespacedName)

	// your logic here

	return ctrl.Result{}, nil
}

func (r *ProductReleaseReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&atlasv1alpha1.ProductRelease{}).
		Complete(r)
}
