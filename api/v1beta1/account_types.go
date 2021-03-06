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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AccountSpec defines the desired state of Account
type AccountSpec struct{}

// AccountStatus defines the observed state of Account
type AccountStatus struct {
	Phase      string             `json:"phase,omitempty"`
	Message    string             `json:"message,omitempty"`
	Reason     string             `json:"reason,omitempty"`
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	Meta       AccountStatusMeta  `json:"meta,omitempty"`
}

// AccountStatusMeta packages meta info on account
type AccountStatusMeta struct {
	NumCoins int      `json:"numCoins,omitempty"`
	CoinList []string `json:"coinList,omitempty"`
	Balance  string   `json:"balance,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase",description="Status of account"
//+kubebuilder:printcolumn:name="Coins",type="integer",JSONPath=".status.meta.numCoins",description="Number of coins"
//+kubebuilder:printcolumn:name="Balance",type="string",JSONPath=".status.meta.balance",description="Balance in account"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Account is the Schema for the accounts API
type Account struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AccountSpec   `json:"spec,omitempty"`
	Status AccountStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// AccountList contains a list of Account
type AccountList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Account `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Account{}, &AccountList{})
}
