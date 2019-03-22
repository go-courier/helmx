package kubetypes

type KubeLocalObjectReference struct {
	Name string `yaml:"name" json:"name" toml:"name"`
}
