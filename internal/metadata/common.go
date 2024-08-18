package metadata

import (
	"fmt"
)

const (
	AuthSecretSuffix = "auth-secret"
	DefaultComponent = "redis"
)

func RedisAuthSecretName(name string) string {
	secretName := fmt.Sprintf("%s-%s", name, AuthSecretSuffix)
	return secretName
}

func RedisMasterComponent() string {
	return "redis-master"
}

func RedisReplicaComponent() string {
	return "redis-replica"
}

func RedisServiceName(name string, component string) string {
	return fmt.Sprintf("%s-%s", name, component)
}
