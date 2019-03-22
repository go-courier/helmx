package tmpl

import (
	"github.com/go-courier/helmx/constants"
	"github.com/go-courier/helmx/kubetypes"
	"github.com/go-courier/helmx/spec"
	"sort"
	"strconv"
	"text/template"
)

var KubeFuncs = template.FuncMap{
	"toKubeEnv":              ToKubeEnv,
	"toKubeInitContainers":   ToKubeInitContainers,
	"toKubeMainContainer":    ToKubeMainContainer,
	"toKubeContainerImage":   ToKubeContainerImage,
	"toKubeVolumeMounts":     ToKubeVolumeMounts,
	"toKubeVolumes":          ToKubeVolumes,
	"toKubeImagePullSecrets": ToKubeImagePullSecrets,
	"toKubeContainerPorts":   ToKubeContainerPorts,
	"toKubeIngressRules":     ToKubeIngressRules,
	"toKubeServicePorts":     ToKubeServicePorts,
}

func ToKubeEnv(envs spec.Envs) kubetypes.KubeEnv {
	keys := make([]string, 0)

	for k := range envs {
		keys = append(keys, k)
	}

	sort.Strings(keys)

	e := kubetypes.KubeEnv{}

	for _, k := range keys {
		e.Env = append(e.Env, kubetypes.KubeEnvVar{
			Name:  k,
			Value: envs[k],
		})
	}

	return e
}

func ToKubeInitContainers(s spec.Spec) kubetypes.KubeInitContainers {
	ss := kubetypes.KubeInitContainers{}

	for i, c := range s.Service.Initials {
		container := kubeContainer(s, c)
		container.Name = container.Name + "-init-" + strconv.FormatInt(int64(i), 10)

		ss.InitContainers = append(ss.InitContainers, container)
	}
	return ss
}

func ToKubeMainContainer(s spec.Spec) kubetypes.KubeContainer {
	ss := kubeContainer(s, s.Service.Container)
	ss.KubeContainerPorts = ToKubeContainerPorts(s)
	return ss
}

func ToKubeContainerImage(s spec.Spec) kubetypes.KubeImage {
	return kubetypes.KubeImage{
		Image:           s.Service.Image.ImageTag(s.Project.DefaultImageTag()),
		ImagePullPolicy: s.Service.ImagePullPolicy,
	}
}

func kubeContainer(s spec.Spec, c spec.Container) kubetypes.KubeContainer {
	ss := kubetypes.KubeContainer{}

	ss.Name = s.Project.FullName()
	ss.KubeImage = ToKubeContainerImage(s)

	ss.WorkingDir = c.WorkingDir
	ss.Command = c.Command
	ss.Args = c.Args
	ss.TTY = c.TTY

	if s.Envs != nil {
		if c.Envs == nil {
			c.Envs = spec.Envs{}
		}
		ss.KubeEnv = ToKubeEnv(c.Envs.Merge(s.Envs))
	}

	ss.KubeVolumeMounts = ToKubeVolumeMounts(s)

	return ss
}

func ToKubeVolumeMounts(s spec.Spec) kubetypes.KubeVolumeMounts {
	ss := kubetypes.KubeVolumeMounts{}
	for _, volumeMount := range s.Service.Container.Mounts {
		ss.VolumeMounts = append(ss.VolumeMounts, toKubeVolumeMount(volumeMount))
	}
	return ss
}

func ToKubeVolumes(s spec.Spec) kubetypes.KubeVolumes {
	ss := kubetypes.KubeVolumes{}
	for name, v := range s.Volumes {
		ss.Volumes = append(ss.Volumes, toKubeVolume(name, v))
	}
	return ss
}

func toKubeVolumeMount(volumeMount spec.VolumeMount) kubetypes.KubeVolumeMount {
	return kubetypes.KubeVolumeMount{
		MountPath: volumeMount.MountPath,
		Name:      volumeMount.Name,
		SubPath:   volumeMount.SubPath,
		ReadOnly:  volumeMount.ReadOnly,
	}
}

func toKubeVolume(name string, v spec.Volume) kubetypes.KubeVolume {
	return kubetypes.KubeVolume{
		Name:             name,
		KubeVolumeSource: v.KubeVolumeSource,
	}
}

func ToKubeImagePullSecrets(s spec.Spec) kubetypes.KubeImagePullSecrets {
	secretNames := map[string]bool{}

	if s.Service.Image.ImagePullSecret != nil {
		secretNames[s.Service.Image.ImagePullSecret.SecretName()] = true
	}

	for _, v := range s.Service.Initials {
		if v.ImagePullSecret != nil {
			secretNames[v.ImagePullSecret.SecretName()] = true
		}
	}

	ss := kubetypes.KubeImagePullSecrets{}

	for name := range secretNames {
		ss.ImagePullSecrets = append(ss.ImagePullSecrets, kubetypes.KubeLocalObjectReference{Name: name})
	}

	return ss
}

func ToKubeContainerPorts(s spec.Spec) kubetypes.KubeContainerPorts {
	ss := kubetypes.KubeContainerPorts{}

	for _, port := range s.Service.Ports {
		p := kubetypes.KubeContainerPort{
			ContainerPort: port.ContainerPort,
		}

		if p.ContainerPort == 0 {
			port.ContainerPort = port.Port
		}

		if port.Protocol == "" {
			p.Protocol = constants.ProtocolTCP
		} else {
			p.Protocol = port.Protocol
		}

		ss.Ports = append(ss.Ports, p)
	}

	return ss
}

func ToKubeServicePorts(s spec.Spec) kubetypes.KubeServicePorts {
	ss := kubetypes.KubeServicePorts{}
	for _, port := range s.Service.Ports {
		p := kubetypes.KubeServicePort{
			Port:       port.Port,
			TargetPort: port.ContainerPort,
		}

		if port.Protocol == "" {
			p.Protocol = constants.ProtocolTCP
		} else {
			p.Protocol = port.Protocol
		}

		ss.Ports = append(ss.Ports, p)
	}
	return ss
}

func ToKubeIngressRules(s spec.Spec) kubetypes.KubeIngressRules {
	ss := kubetypes.KubeIngressRules{}

	for _, r := range s.Service.Ingress {
		rule := kubetypes.IngressRule{
			Host: r.Host,
			HTTP: &kubetypes.HTTPIngressRuleValue{
				Paths: []kubetypes.HTTPIngressPath{
					{
						Path: r.Path,
						Backend: kubetypes.IngressBackend{
							ServiceName: s.Project.FullName(),
							ServicePort: r.Port,
						},
					},
				},
			},
		}

		ss.Rules = append(ss.Rules, rule)
	}

	return ss
}
