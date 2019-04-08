package tmpl

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"text/template"

	"github.com/go-courier/helmx/constants"
	"github.com/go-courier/helmx/kubetypes"
	"github.com/go-courier/helmx/spec"
)

var KubeFuncs = template.FuncMap{
	"toKubeIngressSpec":    ToKubeIngressSpec,
	"toKubeServiceSpec":    ToKubeServiceSpec,
	"toKubeDeploymentSpec": ToKubeDeploymentSpec,
	"toKubeJobSpec":        ToKubeJobSpec,
	"toKubeCronJobSpec":    ToKubeCronJobSpec,
}

func ToKubeServiceSpec(s spec.Spec) kubetypes.KubeServiceSpec {
	ss := kubetypes.KubeServiceSpec{
		Type: kubetypes.ServiceTypeClusterIP,
	}

	for _, port := range s.Service.Ports {
		p := kubetypes.KubeServicePort{
			Port:       port.Port,
			TargetPort: port.ContainerPort,
		}

		if port.IsNodePort {
			ss.Type = kubetypes.ServiceTypeNodePort
			p.NodePort = port.Port
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

func ToKubeDeploymentSpec(s spec.Spec) kubetypes.KubeDeploymentSpec {
	ds := kubetypes.KubeDeploymentSpec{}

	ds.Template.Metadata.Labels = map[string]string{
		"srv": s.Project.FullName(),
	}

	ds.DeploymentOpts = s.Service.DeploymentOpts
	ds.Template.Spec = ToKubePodSpec(s, s.Service.Pod)

	return ds
}

func ToKubeJobSpec(s spec.Spec, job spec.Job) kubetypes.KubeJobSpec {
	js := kubetypes.KubeJobSpec{}
	js.JobOpts = job.JobOpts
	js.Template.Spec = ToKubePodSpec(s, job.Pod)
	return js
}

func ToKubeCronJobSpec(s spec.Spec, job spec.Job) kubetypes.KubeCronJobSpec {
	js := kubetypes.KubeCronJobSpec{}
	if job.Cron != nil {
		js.CronJobOpts = *job.Cron
	}
	js.Template.Spec = ToKubeJobSpec(s, job)
	return js
}

func ToKubePodSpec(s spec.Spec, pod spec.Pod) kubetypes.KubePodSpec {
	ps := kubetypes.KubePodSpec{}

	ps.KubeVolumes = ToKubeVolumes(s)
	ps.KubeTolerations = ToKubeTolerations(s)

	ps.KubeInitContainers = ToKubeInitContainers(s, pod)
	ps.KubeContainers = ToKubeContainers(s, pod)
	ps.KubeImagePullSecrets = ToKubeImagePullSecrets(s, pod)
	ps.PodOpts = pod.PodOpts

	return ps
}

func ToKubeIngressSpec(s spec.Spec) kubetypes.KubeIngressSpec {
	is := kubetypes.KubeIngressSpec{}

	for _, r := range s.Service.Ingresses {
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

		is.Rules = append(is.Rules, rule)
	}

	return is
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

func ToKubeInitContainers(s spec.Spec, pod spec.Pod) kubetypes.KubeInitContainers {
	ss := kubetypes.KubeInitContainers{}

	for i, c := range pod.Initials {
		container := ToKubeContainer(s, c)
		container.Name = container.Name + "-init-" + strconv.FormatInt(int64(i), 10)

		ss.InitContainers = append(ss.InitContainers, container)
	}
	return ss
}

func ToKubeContainers(s spec.Spec, pod spec.Pod) kubetypes.KubeContainers {
	kc := kubetypes.KubeContainers{}

	c := ToKubeContainer(s, pod.Container)

	// only service can be ports
	if s.Service != nil {
		c.KubeContainerPorts = toKubeContainerPorts(s, s.Service.Ports)
	}
	kc.Containers = []kubetypes.KubeContainer{c}

	return kc
}

func ToKubeContainer(s spec.Spec, c spec.Container) kubetypes.KubeContainer {
	ss := kubetypes.KubeContainer{}

	ss.Name = s.Project.FullName()
	ss.KubeImage.Image = c.ImageTag(s.Project.DefaultImageTag())
	ss.KubeImage.ImagePullPolicy = c.ImagePullPolicy

	ss.WorkingDir = c.WorkingDir
	ss.Command = c.Command
	ss.Args = c.Args
	ss.TTY = c.TTY

	if c.LivenessProbe != nil {
		ss.LivenessProbe = &kubetypes.Probe{
			Handler:   c.LivenessProbe.Action.Handler,
			ProbeOpts: c.LivenessProbe.ProbeOpts,
		}
	}

	if c.ReadinessProbe != nil {
		ss.ReadinessProbe = &kubetypes.Probe{
			Handler:   c.ReadinessProbe.Action.Handler,
			ProbeOpts: c.ReadinessProbe.ProbeOpts,
		}
	}

	if c.Lifecycle != nil {
		ss.Lifecycle = &kubetypes.Lifecycle{
			PostStart: &c.Lifecycle.PostStart.Handler,
			PreStop:   &c.Lifecycle.PreStop.Handler,
		}
	}

	if s.Resources != nil {
		if s.Resources.Cpu.Request != 0 {
			ss.Resources.Requests.Cpu = fmt.Sprintf("%dm", s.Resources.Cpu.Request)
		}

		if s.Resources.Cpu.Limit != 0 {
			ss.Resources.Limits.Cpu = fmt.Sprintf("%dm", s.Resources.Cpu.Limit)
		}

		if s.Resources.Memory.Request != 0 {
			ss.Resources.Requests.Memory = fmt.Sprintf("%dMi", s.Resources.Memory.Request)
		}
		if s.Resources.Memory.Limit != 0 {
			ss.Resources.Limits.Memory = fmt.Sprintf("%dMi", s.Resources.Memory.Limit)
		}
	}

	if s.Envs != nil {
		if c.Envs == nil {
			c.Envs = spec.Envs{}
		}
		ss.KubeEnv = ToKubeEnv(c.Envs.Merge(s.Envs))
	}

	ss.KubeVolumeMounts = toKubeVolumeMounts(c)

	return ss
}

func toKubeVolumeMounts(container spec.Container) kubetypes.KubeVolumeMounts {
	ss := kubetypes.KubeVolumeMounts{}
	for _, volumeMount := range container.Mounts {
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

func ToKubeImagePullSecrets(s spec.Spec, pod spec.Pod) kubetypes.KubeImagePullSecrets {
	secretNames := map[string]bool{}
	name := (&spec.ImagePullSecret{}).SecretName()

	if name != "" {
		secretNames[name] = true
	}

	if pod.Image.ImagePullSecret != nil {
		secretNames[s.Service.Image.ImagePullSecret.SecretName()] = true
	}

	for _, v := range pod.Initials {
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

func toKubeContainerPorts(s spec.Spec, ports []spec.Port) kubetypes.KubeContainerPorts {
	ss := kubetypes.KubeContainerPorts{}

	for _, port := range ports {
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

func ToKubeTolerations(s spec.Spec) kubetypes.KubeTolerations {
	kt := kubetypes.KubeTolerations{}

	for _, value := range s.Tolerations {

		kv := strings.Split(value, "=")
		toleration := kubetypes.KubeToleration{
			Key:      kv[0],
			Value:    kv[1],
			Effect:   kubetypes.TolerationEffectNoExecute,
			Operator: kubetypes.TolerationOperatorEqual,
		}

		kt.Tolerations = append(kt.Tolerations, toleration)
	}
	return kt
}
