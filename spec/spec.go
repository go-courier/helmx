package spec

type Spec struct {
	Project *Project `json:"project,omitempty" yaml:"project,omitempty"`

	Service *Service       `json:"service,omitempty" yaml:"service,omitempty"`
	Jobs    map[string]Job `json:"jobs,omitempty" yaml:"jobs,omitempty"`

	Volumes     Volumes   `json:"volumes,omitempty" yaml:"volumes,omitempty"`
	Envs        Envs      `json:"envs,omitempty" yaml:"envs,omitempty"`
	Tolerations []string  `json:"tolerations,omitempty" yaml:"tolerations,omitempty"`
	Resources   *Resource `json:"resources,omitempty" yaml:"resources,omitempty"`
	// just host or service name list
	Upstreams []string `json:"upstreams,omitempty" yaml:"upstreams,omitempty"`
}
