/*
Copyright 2022 kubetrail.io Authors.

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

package v1beta1

import (
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var coinlog = logf.Log.WithName("coin-resource")

func (r *Coin) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!

//+kubebuilder:webhook:path=/mutate-crypto-kubetrail-io-v1beta1-coin,mutating=true,failurePolicy=fail,sideEffects=None,groups=crypto.kubetrail.io,resources=coins,verbs=create;update,versions=v1beta1,name=mcoin.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &Coin{}

// Default implements webhook.Defaulter so a webhook will be registered for the type
func (r *Coin) Default() {
	coinlog.Info("default", "name", r.Name)

	// TODO(user): fill in your defaulting logic.
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:path=/validate-crypto-kubetrail-io-v1beta1-coin,mutating=false,failurePolicy=fail,sideEffects=None,groups=crypto.kubetrail.io,resources=coins,verbs=create;update,versions=v1beta1,name=vcoin.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &Coin{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type
func (r *Coin) ValidateCreate() error {
	coinlog.Info("validate create", "name", r.Name)

	// TODO(user): fill in your validation logic upon object creation.
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type
func (r *Coin) ValidateUpdate(old runtime.Object) error {
	coinlog.Info("validate update", "name", r.Name)

	// TODO(user): fill in your validation logic upon object update.
	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type
func (r *Coin) ValidateDelete() error {
	coinlog.Info("validate delete", "name", r.Name)

	// TODO(user): fill in your validation logic upon object deletion.
	return nil
}
