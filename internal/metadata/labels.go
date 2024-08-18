package metadata

type Label map[string]string

func CommonLabels(instanceName string) Label {
	return Label{
		"app.kubernetes.io/name":    instanceName,
		"app.kubernetes.io/part-of": "redis",
	}
}

func ResourceLabels(instanceName string, instanceLabels Label) Label {
	common := CommonLabels(instanceName)
	for k, v := range common {
		instanceLabels[k] = v
	}
	return instanceLabels
}

func LabelSelector(instanceName string, resource string) Label {
	return Label{
		"app.kubernetes.io/name":      instanceName,
		"app.kubernetes.io/component": resource,
	}
}
