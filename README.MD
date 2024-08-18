Presented Redis operator creates two deployments for redis master and redis replicas respectively. Operator watches for the Redis custom resource definition.


Build and deploy redis operator:
- configure KUBECONFIG to connect to your kubernetes cluster. This operator was tested on Minikube running in docker container
- edit Makefile and replace IMG value with your docker repository
- build docker container
```
make docker-buildx
```
- generate kustomize manifests for controller deployment
```
make build-installer
```
- install operator manifests
```
make install
```
- deploy operator to cluster
```
make deploy
```

Alternatively:
-  run
```
make install
```
to deploy operator resources to kubernetes

- to run operator locally
```
make install run
```

To deploy redis run
```
kubectl apply -f test-redis.yaml
```

Run
```
kubectl explain redis
```
to see all availabe options for Redis deployment