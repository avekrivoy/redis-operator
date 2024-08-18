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

package controller

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/log"

	cachev1alpha1 "github.com/avekrivoy/redis-operator/api/v1alpha1"
	"github.com/avekrivoy/redis-operator/internal/metadata"
	resources "github.com/avekrivoy/redis-operator/internal/resources"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8serrors "k8s.io/apimachinery/pkg/api/errors"
)

// RedisReconciler reconciles a Redis object
type RedisReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=cache.assignment.yazio.com,resources=redis,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=cache.assignment.yazio.com,resources=redis/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=cache.assignment.yazio.com,resources=redis/finalizers,verbs=update
//+kubebuilder:rbac:groups=core,resources=services;secrets,verbs=create;update;delete;get;list;watch
//+kubebuilder:rbac:groups=apps,resources=deployments,verbs=create;update;delete;get;list;watch

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Redis object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *RedisReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)
	redis, err := r.getRedisInstance(ctx, req.NamespacedName)

	if client.IgnoreNotFound(err) != nil {
		return ctrl.Result{}, err
	} else if k8serrors.IsNotFound(err) {
		// No need to requeue if the resource no longer exists
		return ctrl.Result{}, nil
	}

	resourceBuilder := resources.RedisResourceBuilder{
		Instance: redis,
		Scheme:   r.Scheme,
	}

	builders := resourceBuilder.ResourceBuilders()

	for _, builder := range builders {
		if builder.IsDeployed() {
			resource, err := builder.Build()
			if err != nil {
				return ctrl.Result{}, err
			}

			// Do not recreate Redis auth secret if already exists
			if resource.GetName() == metadata.RedisAuthSecretName(redis.Name) {
				secret := &corev1.Secret{}
				err = r.Client.Get(context.TODO(), types.NamespacedName{
					Name:      metadata.RedisAuthSecretName(redis.Name),
					Namespace: redis.Namespace,
				}, secret)
				if err == nil {
					logger.Info("Redis auth secret already exists. Skipping resource creation")
					continue
				}
			}

			_, apiError := controllerutil.CreateOrUpdate(ctx, r.Client, resource, func() error {
				return builder.Update(resource)
			})

			if apiError != nil {
				return ctrl.Result{}, err
			}
		}
	}

	logger.Info("Finished reconciling")
	return ctrl.Result{}, nil
}

func (r *RedisReconciler) getRedisInstance(ctx context.Context, namespacedName types.NamespacedName) (*cachev1alpha1.Redis, error) {
	redisInstance := &cachev1alpha1.Redis{}
	err := r.Get(ctx, namespacedName, redisInstance)
	return redisInstance, err
}

func (r *RedisReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&cachev1alpha1.Redis{}).
		Owns(&appsv1.Deployment{}).
		Owns(&corev1.Service{}).
		Owns(&corev1.Secret{}).
		Complete(r)
}
