/*
Copyright 2021.

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

package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// DnsConfigSpec defines the desired state of DnsConfig
type DnsConfigSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of DnsConfig. Edit dnsconfig_types.go to remove/update
	IPHosts []IPHosts `json:"ipHosts,omitempty"`
}

type IPHosts struct {
	Option   string `json:"option,omitempty"`
	ConfIp   string `json:"confIp,omitempty"`
	ConfHost string `json:"confHost,omitempty"`
}

// DnsConfigStatus defines the observed state of DnsConfig
//type DnsConfigStatus struct {
//	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
//	// Important: Run "make" to regenerate code after modifying this file
//}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// DnsConfig is the Schema for the dnsconfigs API
type DnsConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec DnsConfigSpec `json:"spec,omitempty"`
}

//+kubebuilder:object:root=true

// DnsConfigList contains a list of DnsConfig
type DnsConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []DnsConfig `json:"items"`
}

func init() {
	SchemeBuilder.Register(&DnsConfig{}, &DnsConfigList{})
}
