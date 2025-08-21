/*
Copyright 2025.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json:"-" or json:"fieldname,omitempty".

// ServiceConfig defines the configuration for a single service
type ServiceConfig struct {
	// Enabled indicates whether this service should be deployed
	Enabled bool `json:"enabled,omitempty"`

	// Replicas specifies the number of replicas for this service
	Replicas *int32 `json:"replicas,omitempty"`

	// Image specifies the container image to use
	Image string `json:"image,omitempty"`

	// Tag specifies the image tag
	Tag string `json:"tag,omitempty"`

	// Resources specifies resource requirements
	Resources *ResourceRequirements `json:"resources,omitempty"`
}

// ResourceRequirements defines resource requirements for a service
type ResourceRequirements struct {
	// Limits describes the maximum amount of compute resources allowed
	Limits map[string]string `json:"limits,omitempty"`

	// Requests describes the minimum amount of compute resources required
	Requests map[string]string `json:"requests,omitempty"`
}

// DatabaseConfig defines database configuration
type DatabaseConfig struct {
	// Enabled indicates whether to deploy the database
	Enabled bool `json:"enabled,omitempty"`

	// Type specifies the database type (mysql, postgres, etc.)
	Type string `json:"type,omitempty"`

	// Image specifies the database container image
	Image string `json:"image,omitempty"`

	// Tag specifies the database image tag
	Tag string `json:"tag,omitempty"`

	// StorageSize specifies the storage size for the database
	StorageSize string `json:"storageSize,omitempty"`

	// StorageClass specifies the storage class
	StorageClass string `json:"storageClass,omitempty"`
}

// ClusterTesterSpec defines the desired state of ClusterTester
type ClusterTesterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// CoffeeShop service configuration
	CoffeeShop ServiceConfig `json:"coffeeShop,omitempty"`

	// PetStore service configuration
	PetStore ServiceConfig `json:"petStore,omitempty"`

	// Restaurant service configuration
	Restaurant ServiceConfig `json:"restaurant,omitempty"`

	// CollegeAdmission service configuration
	CollegeAdmission ServiceConfig `json:"collegeAdmission,omitempty"`

	// ElectronicsStore service configuration
	ElectronicsStore ServiceConfig `json:"electronicsStore,omitempty"`

	// ElectronicsStoreTracing service configuration
	ElectronicsStoreTracing ServiceConfig `json:"electronicsStoreTracing,omitempty"`

	// Database configuration for services that need it
	Database DatabaseConfig `json:"database,omitempty"`

	// Global configuration
	Global GlobalConfig `json:"global,omitempty"`
}

// GlobalConfig defines global configuration options
type GlobalConfig struct {
	// Namespace specifies the target namespace for deployments
	Namespace string `json:"namespace,omitempty"`

	// ImagePullPolicy specifies the image pull policy
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`

	// ServiceType specifies the default service type (ClusterIP, NodePort, LoadBalancer)
	ServiceType string `json:"serviceType,omitempty"`

	// IngressEnabled indicates whether to create ingress resources
	IngressEnabled bool `json:"ingressEnabled,omitempty"`

	// IngressHost specifies the base host for ingress
	IngressHost string `json:"ingressHost,omitempty"`
}

// ServiceStatus defines the status of a deployed service
type ServiceStatus struct {
	// Name of the service
	Name string `json:"name"`

	// Ready indicates if the service is ready
	Ready bool `json:"ready"`

	// Replicas indicates the current number of replicas
	Replicas int32 `json:"replicas"`

	// ReadyReplicas indicates the number of ready replicas
	ReadyReplicas int32 `json:"readyReplicas"`

	// Endpoint indicates the service endpoint
	Endpoint string `json:"endpoint,omitempty"`
}

// ClusterTesterStatus defines the observed state of ClusterTester
type ClusterTesterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Phase indicates the current phase of the ClusterTester deployment
	Phase string `json:"phase,omitempty"`

	// Conditions represents the latest available observations of the ClusterTester's state
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// Services contains the status of individual services
	Services []ServiceStatus `json:"services,omitempty"`

	// ObservedGeneration reflects the generation observed by the controller
	ObservedGeneration int64 `json:"observedGeneration,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Namespaced
//+kubebuilder:printcolumn:name="Phase",type=string,JSONPath=`.status.phase`
//+kubebuilder:printcolumn:name="Ready Services",type=string,JSONPath=`.status.services[?(@.ready==true)].name`
//+kubebuilder:printcolumn:name="Age",type=date,JSONPath=`.metadata.creationTimestamp`

// ClusterTester is the Schema for the clustertesters API
type ClusterTester struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ClusterTesterSpec   `json:"spec,omitempty"`
	Status ClusterTesterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// ClusterTesterList contains a list of ClusterTester
type ClusterTesterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterTester `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterTester{}, &ClusterTesterList{})
}
