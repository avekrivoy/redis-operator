package resources

import (
	cachev1alpha1 "github.com/avekrivoy/redis-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type RedisResourceBuilder struct {
	Instance *cachev1alpha1.Redis
	Scheme   *runtime.Scheme
}

type ResourceBuilder interface {
	Build() (client.Object, error)
	Update(client.Object) error
	IsDeployed() bool
}

func (builder *RedisResourceBuilder) ResourceBuilders() []ResourceBuilder {

	builders := []ResourceBuilder{
		builder.RedisAuthSecret(),
		builder.RedisMasterService(),
		builder.RedisMasterDeployment(),
		builder.RedisReplicaService(),
		builder.RedisReplicaDeployment(),
	}
	return builders
}
