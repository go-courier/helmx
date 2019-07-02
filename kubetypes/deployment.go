package kubetypes

type KubeDeploymentSpec struct {
	DeploymentOpts `yaml:",inline"`
	Template       struct {
		KubeMetadata `yaml:",inline"`
		Spec         KubePodSpec `yaml:"spec"`
	} `yaml:"template"`
}

type DeploymentOpts struct {
	Replicas *int32 `yaml:"replicas,omitempty" json:"replicas,omitempty"`
}
