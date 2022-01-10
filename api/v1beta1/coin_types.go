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

// CoinSpec defines the desired state of Coin
type CoinSpec struct {
	Ticker   string `json:"ticker,omitempty"`
	Currency string `json:"currency,omitempty"`
	NumCoins string `json:"numCoins,omitempty"`
}

// CoinStatus defines the observed state of Coin
type CoinStatus struct {
	Phase      string             `json:"phase,omitempty"`
	Message    string             `json:"message,omitempty"`
	Reason     string             `json:"reason,omitempty"`
	Conditions []metav1.Condition `json:"conditions,omitempty"`
	Meta       CoinStatusMeta     `json:"meta,omitempty"`
}

// CoinStatusMeta packages meta info on coin price
type CoinStatusMeta struct {
	Price   string `json:"price,omitempty"`
	Balance string `json:"balance,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:printcolumn:name="Status",type="string",JSONPath=".status.phase",description="Status of coin"
//+kubebuilder:printcolumn:name="Ticker",type="string",JSONPath=".spec.ticker",description="Ticker of coin"
//+kubebuilder:printcolumn:name="Price",type="string",JSONPath=".status.meta.price",description="Price of coin"
//+kubebuilder:printcolumn:name="NumCoins",type="string",JSONPath=".spec.numCoins",description="Number of coins"
//+kubebuilder:printcolumn:name="Balance",type="string",JSONPath=".status.meta.balance",description="Balance in coin"
//+kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp"

// Coin is the Schema for the coins API
type Coin struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   CoinSpec   `json:"spec,omitempty"`
	Status CoinStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// CoinList contains a list of Coin
type CoinList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Coin `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Coin{}, &CoinList{})
}
