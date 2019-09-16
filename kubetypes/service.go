package kubetypes

import (
	"github.com/go-courier/helmx/constants"
)

type ServiceType string

const (
	ServiceTypeClusterIP ServiceType = "ClusterIP"
	ServiceTypeNodePort  ServiceType = "NodePort"
)

type KubeServiceSpec struct {
	Type  ServiceType       `yaml:"type,omitempty"`
	Ports []KubeServicePort `yaml:"ports,omitempty"`
}

type KubeServicePort struct {
	Name       string             `yaml:"name,omitempty"`
	NodePort   uint16             `yaml:"nodePort,omitempty"`
	Port       uint16             `yaml:"port"`
	TargetPort uint16             `yaml:"targetPort"`
	Protocol   constants.Protocol `yaml:"protocol"`
}
