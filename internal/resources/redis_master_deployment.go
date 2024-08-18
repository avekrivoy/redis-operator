package resources

import (
	"fmt"

	metadata "github.com/avekrivoy/redis-operator/internal/metadata"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type RedisMasterDeploymentBuilder struct {
	*RedisResourceBuilder
}

func (builder *RedisResourceBuilder) RedisMasterDeployment() *RedisMasterDeploymentBuilder {
	return &RedisMasterDeploymentBuilder{builder}
}

func (builder *RedisMasterDeploymentBuilder) Build() (client.Object, error) {
	component := metadata.RedisMasterComponent()
	deploymentName := fmt.Sprintf("%s-%s", builder.Instance.Name, component)

	return &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName,
			Namespace: builder.Instance.Namespace,
		},
	}, nil
}

func (builder *RedisMasterDeploymentBuilder) Update(object client.Object) error {
	component := metadata.RedisMasterComponent()

	deploymentLabels := metadata.Label{
		"app.kubernetes.io/component": component,
	}

	labels := metadata.ResourceLabels(builder.Instance.Name, deploymentLabels)
	redisImage := fmt.Sprintf("%s:%s", builder.Instance.Spec.Common.Image.ImageRepository, builder.Instance.Spec.Common.Image.ImageTag)

	var redisAuthSecretName string

	// Check if existing auth secret was specified
	if builder.Instance.Spec.Common.Auth.ExistingSecret == "" {
		redisAuthSecretName = metadata.RedisAuthSecretName(builder.Instance.Name)
	} else {
		redisAuthSecretName = builder.Instance.Spec.Common.Auth.ExistingSecret
	}

	deployment := object.(*appsv1.Deployment)
	deployment.ObjectMeta.Labels = labels
	deployment.Spec.Replicas = &builder.Instance.Spec.Master.Count
	deployment.Spec.Template.ObjectMeta.Labels = labels
	deployment.Spec.Selector = &metav1.LabelSelector{
		MatchLabels: metadata.LabelSelector(builder.Instance.Name, component),
	}

	deployment.Spec.Template.Spec = corev1.PodSpec{
		Containers: []corev1.Container{{
			Image:           redisImage,
			ImagePullPolicy: corev1.PullPolicy(builder.Instance.Spec.Common.Image.ImagePullPolicy),
			Name:            "redis",
			Ports: []corev1.ContainerPort{{
				ContainerPort: 6379,
				Name:          "redis",
			}},
			EnvFrom: []corev1.EnvFromSource{{
				SecretRef: &corev1.SecretEnvSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: redisAuthSecretName,
					},
				},
			}},
			Env: []corev1.EnvVar{
				{
					Name:  "REDIS_REPLICATION_MODE",
					Value: "master",
				},
			},
		}},
	}

	if err := controllerutil.SetControllerReference(builder.Instance, deployment, builder.Scheme); err != nil {
		return fmt.Errorf("failed setting controller reference: %w", err)
	}
	return nil
}

func (builder *RedisMasterDeploymentBuilder) IsDeployed() bool {
	return builder.Instance.Spec.Master.Count >= 1
}
