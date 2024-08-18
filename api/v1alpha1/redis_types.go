/*
Copyright 2024.

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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// RedisSpec defines the desired state of Redis
type RedisSpec struct {
	// Common values for Redis deployment
	Common RedisCommonSpec `json:"common,omitempty"`
	// Redis master parameters
	Master RedisMasterSpec `json:"master,omitempty"`
	// Redis replica parameters
	Replica RedisReplicaSpec `json:"replica,omitempty"`
}

type RedisCommonSpec struct {
	// Redis image parameters
	// +kubebuilder:default={}
	Image RedisImageSpec `json:"image,omitempty"`
	// Storage class for Redis PVCs. Defaults to standard
	// +kubebuilder:default:="standard"
	StorageClass string `json:"storageClass,omitempty"`
	// Redis Authentication configuration
	Auth RedisAuthSpec `json:"auth,omitempty"`
}

type RedisImageSpec struct {
	// Docker image repository
	// +kubebuilder:default:="bitnami/redis"
	ImageRepository string `json:"imageRegistry,omitempty"`
	// Docker image tag
	// +kubebuilder:default:="7.2.5"
	ImageTag string `json:"imageTag,omitempty"`
	// List of ImagePullSecrets resources
	ImagePullSecrets []corev1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
	// Defaults to IfNotPresent
	// +kubebuilder:default:=IfNotPresent
	ImagePullPolicy string `json:"imagePullPolicy,omitempty"`
}

type RedisAuthSpec struct {
	// Enable password authentication
	// +kubebuilder:default:=true
	Enabled bool `json:"enabled,omitempty"`
	// The name of an existing secret with Redis credentials. Should contain REDIS_PASSWORD key
	ExistingSecret string `json:"existingSecret,omitempty"`
}

type RedisMasterSpec struct {
	// Number of Redis pods
	// +kubebuilder:default:=1
	Count int32 `json:"count,omitempty"`
	// Type of Redis deployment. Defaults to 'deployment'
	// +kubebuilder:default:=deployment
	Kind string `json:"kind,omitempty"`
	// Redis PVC configuration
	//Persistence RedisMasterPersistence `json:"persistence,omitempty"`
}

// type RedisMasterPersistence struct {
// 	// Enable persisten storage
// 	// +kubebuilder:default:=false
// 	Enabled bool `json:"enabled,omitempty"`
// 	// PVC size
// 	Size string `json:"size,omitempty"`
// }

type RedisReplicaSpec struct {
	// Number of Redis pods
	// +kubebuilder:default:=0
	Count int32 `json:"count,omitempty"`
	// Type of Redis deployment. Defaults to 'deployment'
	// +kubebuilder:default:=deployment
	Kind string `json:"kind,omitempty"`
	// Redis PVC configuration
	//Persistence RedisReplicaPersistence `json:"persistence,omitempty"`
}

// type RedisReplicaPersistence struct {
// 	// Enable persisten storage
// 	// +kubebuilder:default:=false
// 	Enabled bool `json:"enabled,omitempty"`
// 	// PVC size
// 	Size string `json:"size,omitempty"`
// }

// RedisStatus defines the observed state of Redis
type RedisStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Conditions []metav1.Condition `json:"conditions"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// Redis is the Schema for the redis API
type Redis struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RedisSpec   `json:"spec,omitempty"`
	Status RedisStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RedisList contains a list of Redis
type RedisList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Redis `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Redis{}, &RedisList{})
}
