/*
Copyright 2022.

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

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	gatewaysv1 "f5.com/bigip-k8s-gateway/api/v1"
)

// AdcReconciler reconciles a Adc object
type AdcReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=gateways.f5.com,resources=adcs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=gateways.f5.com,resources=adcs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=gateways.f5.com,resources=adcs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Adc object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.13.0/pkg/reconcile
func (r *AdcReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	// TODO(user): your logic here
	ctrl.Log.Info("reconcile is called to do action on " + req.NamespacedName.String())
	var obj gatewaysv1.Adc
	err := r.Get(ctx, req.NamespacedName, &obj)
	if err != nil {
		ctrl.Log.Info("Failed to get object: " + err.Error())
	}
	ctrl.Log.Info("spec contains " + obj.Spec.Foo)
	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AdcReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&gatewaysv1.Adc{}).
		Complete(r)
}