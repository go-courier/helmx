package kubetypes

import "github.com/go-courier/helmx/constants"

type KubeInitContainers struct {
	InitContainers []KubeContainer `yaml:"initContainers,omitempty"`
}

type KubeContainers struct {
	Containers []KubeContainer `yaml:"containers,omitempty"`
}

type KubeImagePullSecrets struct {
	ImagePullSecrets []KubeLocalObjectReference `yaml:"imagePullSecrets,omitempty"`
}

type KubeContainer struct {
	Name               string   `yaml:"name"`
	Command            []string `yaml:"command,omitempty"`
	Args               []string `yaml:"args,omitempty"`
	WorkingDir         string   `yaml:"workingDir,omitempty"`
	TTY                bool     `yaml:"tty,omitempty"`
	KubeImage          `yaml:",inline"`
	KubeContainerPorts `yaml:",inline"`
	KubeVolumeMounts   `yaml:",inline"`
	KubeEnv            `yaml:",inline"`
}

type KubeImage struct {
	Image           string               `yaml:"image,omitempty"`
	ImagePullPolicy constants.PullPolicy `yaml:"imagePullPolicy,omitempty"`
}

type KubeEnv struct {
	Env []KubeEnvVar `yaml:"env,omitempty"`
}

type KubeEnvVar struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type KubeContainerPorts struct {
	Ports []KubeContainerPort `yaml:"ports,omitempty"`
}

type KubeContainerPort struct {
	ContainerPort uint16             `yaml:"containerPort"`
	Protocol      constants.Protocol `yaml:"protocol,omitempty"`
}

type KubeVolumeMounts struct {
	VolumeMounts []KubeVolumeMount `yaml:"volumeMounts,omitempty"`
}

type KubeVolumeMount struct {
	Name      string `yaml:"name"`
	MountPath string `yaml:"mountPath"`
	SubPath   string `yaml:"subPath,omitempty"`
	ReadOnly  bool   `yaml:"readOnly,omitempty"`
}
