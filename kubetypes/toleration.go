package kubetypes

type KubeTolerations struct {
	Tolerations []KubeToleration `yaml:"tolerations,omitempty"`
}

type KubeToleration struct {
	Key               string `yaml:"key"`
	Value             string `yaml:"value,omitempty"`
	Effect            string `yaml:"effect,omitempty"`
	Operator          string `yaml:"operator,omitempty"`
	TolerationSeconds *int64 `yaml:"tolerationSeconds,omitempty"`
}
