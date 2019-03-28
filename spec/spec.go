package spec

type Spec struct {
	Project     Project                `json:"project" yaml:"project"`
	Service     Service                `json:"service" yaml:"service"`
	Volumes     Volumes                `json:"volumes,omitempty" yaml:"volumes,omitempty"`
	Envs        Envs                   `json:"envs,omitempty" yaml:"envs,omitempty"`
	Values      map[string]interface{} `json:"values,omitempty" yaml:"values,omitempty"`
	Tolerations []string               `json:"tolerations,omitempty" yaml:"tolerations,omitempty"`
	Resources   Resource               `json:"resources,omitempty",yaml:"resources,omitempty"`
	// just host or service name list
	Upstreams []string `json:"upstreams,omitempty" yaml:"upstreams,omitempty"`
}
