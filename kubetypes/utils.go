package kubetypes

type KubeLocalObjectReference struct {
	Name string `yaml:"name" json:"name" toml:"name"`
}

type KubeOption struct {
	Name  string `yaml:"name" json:"name" toml:"name"`
	Value string `yaml:"value,omitempty" json:"value,omitempty" toml:"value,omitempty"`
}

type KubeMetadata struct {
	Metadata struct {
		Labels map[string]string `yaml:"labels,omitempty"`
	} `yaml:"metadata,omitempty"`
}

type KubeSelector struct {
	Selector struct {
		MatchLabels map[string]string `yaml:"matchLabels,omitempty"`
	} `yaml:"selector"`
}
