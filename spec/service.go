package spec

import "github.com/go-courier/helmx/kubetypes"

type Service struct {
	Pod                      `yaml:",inline"`
	Ingress                  `yaml:",inline"`
	kubetypes.DeploymentOpts `yaml:",inline"`

	Ports []Port `json:"ports,omitempty" yaml:"ports,omitempty"`
}

type Ingress struct {
	IngressClass string        `json:"ingressClass,omitempty" yaml:"ingressClass,omitempty"`
	Ingresses    []IngressRule `json:"ingresses,omitempty" yaml:"ingresses,omitempty"`
}

type Pod struct {
	Initials          []Container `json:"initials,omitempty" yaml:"initials,omitempty"`
	Container         `yaml:",inline"`
	kubetypes.PodOpts `yaml:",inline"`
	Hosts             []Hosts `yaml:"hosts,omitempty" json:"hosts,omitempty"`
}
