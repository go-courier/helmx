package kubetypes

import "github.com/go-courier/helmx/constants"

type ServiceType string

const (
	ServiceTypeClusterIP    ServiceType = "ClusterIP"
	ServiceTypeNodePort     ServiceType = "NodePort"
	ServiceTypeLoadBalancer ServiceType = "LoadBalancer"
	ServiceTypeExternalName ServiceType = "ExternalName"
)

type KubeServiceType struct {
	Type ServiceType `yaml:"type,omitempty"`
}

type KubeServicePorts struct {
	Ports []KubeServicePort `yaml:"ports,omitempty"`
}

type KubeServicePort struct {
	Port       uint16             `yaml:"port"`
	TargetPort uint16             `yaml:"targetPort"`
	Protocol   constants.Protocol `yaml:"protocol"`
}
