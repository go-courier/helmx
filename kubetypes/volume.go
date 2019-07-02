package kubetypes

type KubeVolumes struct {
	Volumes []KubeVolume `yaml:"volumes,omitempty"`
}

type KubeVolume struct {
	Name             string `yaml:"name"`
	KubeVolumeSource `yaml:",inline"`
}

type KubeVolumeSource struct {
	EmptyDir              *EmptyDirVolumeSource              `toml:"emptyDir,omitempty" json:"emptyDir,omitempty" yaml:"emptyDir,omitempty"`
	PersistentVolumeClaim *PersistentVolumeClaimVolumeSource `toml:"persistentVolumeClaim,omitempty" json:"persistentVolumeClaim,omitempty" yaml:"persistentVolumeClaim,omitempty"`
	Secret                *SecretVolumeSource                `toml:"secret,omitempty" json:"secret,omitempty" yaml:"secret,omitempty"`
	ConfigMap             *ConfigMapVolumeSource             `toml:"configMap,omitempty" json:"configMap,omitempty" yaml:"configMap,omitempty"`
	HostPath              *HostPathVolumeSource              `toml:"hostPath,omitempty" json:"hostPath,omitempty" yaml:"hostPath,omitempty"`
}
