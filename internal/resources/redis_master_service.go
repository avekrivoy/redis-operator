package resources

import (
	"fmt"

	metadata "github.com/avekrivoy/redis-operator/internal/metadata"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

type RedisMasterServiceBuilder struct {
	*RedisResourceBuilder
}

func (builder *RedisResourceBuilder) RedisMasterService() *RedisMasterServiceBuilder {
	return &RedisMasterServiceBuilder{builder}
}

func (builder *RedisMasterServiceBuilder) Build() (client.Object, error) {
	component := metadata.RedisMasterComponent()
	svcName := metadata.RedisServiceName(builder.Instance.Name, component)

	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      svcName,
			Namespace: builder.Instance.Namespace,
		},
	}, nil
}

func (builder *RedisMasterServiceBuilder) Update(object client.Object) error {
	component := metadata.RedisMasterComponent()
	svcLabels := metadata.Label{
		"app.kubernetes.io/component": component,
	}

	svc := object.(*corev1.Service)
	svc.Labels = metadata.ResourceLabels(builder.Instance.Name, svcLabels)

	svc.ObjectMeta.Labels = svcLabels
	svc.Spec = corev1.ServiceSpec{
		Selector: metadata.LabelSelector(builder.Instance.Name, component),
		Ports: []corev1.ServicePort{
			{
				Port:     6379,
				Protocol: corev1.ProtocolTCP,
			},
		},
		Type: corev1.ServiceTypeClusterIP,
	}

	if err := controllerutil.SetControllerReference(builder.Instance, svc, builder.Scheme); err != nil {
		return fmt.Errorf("failed setting controller reference: %w", err)
	}

	return nil
}

func (builder *RedisMasterServiceBuilder) IsDeployed() bool {
	return builder.Instance.Spec.Master.Count > 0
}
