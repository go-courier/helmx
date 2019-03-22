package spec

type Service struct {
	Container `yaml:",inline"`

	Replicas int           `json:"replicas,omitempty" yaml:"replicas,omitempty"`
	Initials []Container   `json:"initials,omitempty" yaml:"initials,omitempty"`
	Ports    []Port        `json:"ports,omitempty" yaml:"ports,omitempty"`
	Ingress  []IngressRule `json:"ingress,omitempty" yaml:"ingress,omitempty"`
}
