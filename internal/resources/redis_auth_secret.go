package resources

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"

	metadata "github.com/avekrivoy/redis-operator/internal/metadata"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type RedisAuthSecretBuilder struct {
	*RedisResourceBuilder
}

func (builder *RedisResourceBuilder) RedisAuthSecret() *RedisAuthSecretBuilder {
	return &RedisAuthSecretBuilder{builder}
}

func (builder *RedisAuthSecretBuilder) Build() (client.Object, error) {

	// secretLabels := metadata.Label{
	// 	"app.kubernetes.io/component": metadata.DefaultComponent,
	// }

	// labels := metadata.ResourceLabels(builder.Instance.Name, secretLabels)

	// redisPwd, err := randomEncodedString(24)
	// if err != nil {
	// 	return nil, err
	// }

	// return &corev1.Secret{
	// 	ObjectMeta: metav1.ObjectMeta{
	// 		Name:      metadata.RedisAuthSecretName(builder.Instance.Name),
	// 		Namespace: builder.Instance.Namespace,
	// 		Labels:    labels,
	// 	},
	// 	Type: corev1.SecretTypeOpaque,
	// 	Data: map[string][]byte{
	// 		"REDIS_PASSWORD": []byte(redisPwd),
	// 	},
	// }, nil

	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      metadata.RedisAuthSecretName(builder.Instance.Name),
			Namespace: builder.Instance.Namespace,
		},
		Type: corev1.SecretTypeOpaque,
		// Data: map[string][]byte{
		// 	"REDIS_PASSWORD":          []byte(),
		// },
	}, nil

}

func (builder *RedisAuthSecretBuilder) Update(object client.Object) error {
	secret := object.(*corev1.Secret)

	secretLabels := metadata.Label{
		"app.kubernetes.io/component": metadata.DefaultComponent,
	}
	labels := metadata.ResourceLabels(builder.Instance.Name, secretLabels)

	redisPwd, err := randomEncodedString(24)
	if err != nil {
		return err
	}

	secret.Labels = labels
	secret.Data = map[string][]byte{
		"REDIS_PASSWORD": []byte(redisPwd),
	}

	if err := controllerutil.SetControllerReference(builder.Instance, secret, builder.Scheme); err != nil {
		return fmt.Errorf("failed setting controller reference: %w", err)
	}

	return nil
}

func (builder *RedisAuthSecretBuilder) IsDeployed() bool {
	return builder.Instance.Spec.Common.Auth.ExistingSecret == ""
}

func randomEncodedString(dataLen int) (string, error) {
	randomBytes := make([]byte, dataLen)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(randomBytes), nil
}
