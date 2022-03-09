package tmpl

import (
	"fmt"
	"sort"
	"strconv"
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
	"toKubeRoleRules":      ToKubeRoleRoles,
}

func ToKubeServiceSpec(s spec.Spec) kubetypes.KubeServiceSpec {
	ss := kubetypes.KubeServiceSpec{
		Type: kubetypes.ServiceTypeClusterIP,
	}

	for _, port := range s.Service.Ports {
		appProtocol := port.AppProtocol

		if appProtocol == "" {
			appProtocol = "http"
		}

		p := kubetypes.KubeServicePort{
			Name:       fmt.Sprintf("%s-%d", appProtocol, port.Port),
			Port:       port.Port,
			TargetPort: port.ContainerPort,
		}

		if port.IsNodePort {
			ss.Type = kubetypes.ServiceTypeNodePort

			p.Name = "np-" + p.Name

			if port.Port >= 20000 && port.Port <= 40000 {
				p.NodePort = port.Port
			}
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

	if s.Labels != nil {
		for k, v := range s.Labels {
			ds.Template.Metadata.Labels[k] = v
		}
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
	ps.HostAliases = ToKubeHosts(s)
	ps.KubeTopologySpreadConstraints = ToKubeTopologySpreadConstraints(pod)
	ps.KubeAffinity = ToKubeAffinity(pod)

	return ps
}

func ToKubeRoleRoles(s spec.Spec) []kubetypes.KubeRoleRule {
	rules := make([]kubetypes.KubeRoleRule, 0)

	for _, r := range s.Service.ServiceAccountRoleRules {
		rule := kubetypes.KubeRoleRule{
			ApiGroups:     r.ApiGroups,
			Resources:     r.Resources,
			ResourceNames: r.ResourceNames,
			Verbs:         r.Verbs,
		}

		rules = append(rules, rule)
	}

	return rules
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

	if c.SecurityContext != nil {
		ss.SecurityContext = &kubetypes.SecurityContext{}
		if c.SecurityContext.Capabilities != nil {
			ss.SecurityContext.Capabilities = c.SecurityContext.Capabilities
		}

		if c.SecurityContext.RunAsUser != nil {
			ss.SecurityContext.RunAsUser = c.SecurityContext.RunAsUser
		}

		if c.SecurityContext.RunAsGroup != nil {
			ss.SecurityContext.RunAsGroup = c.SecurityContext.RunAsGroup
		}

		if c.SecurityContext.RunAsNonRoot != nil {
			ss.SecurityContext.RunAsNonRoot = c.SecurityContext.RunAsNonRoot
		}

		if c.SecurityContext.ReadOnlyRootFilesystem != nil {
			ss.SecurityContext.ReadOnlyRootFilesystem = c.SecurityContext.ReadOnlyRootFilesystem
		}

		if c.SecurityContext.AllowPrivilegeEscalation != nil {
			ss.SecurityContext.AllowPrivilegeEscalation = c.SecurityContext.AllowPrivilegeEscalation
		}

		if c.SecurityContext.ProcMount != nil {
			ss.SecurityContext.ProcMount = c.SecurityContext.ProcMount
		}

		if c.SecurityContext.Privileged != nil {
			ss.SecurityContext.Privileged = c.SecurityContext.Privileged
		}

		if c.SecurityContext.SELinuxOptions != nil {
			ss.SecurityContext.SELinuxOptions = c.SecurityContext.SELinuxOptions
		}
	}

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
		ss.Lifecycle = &kubetypes.Lifecycle{}
		if c.Lifecycle.PostStart != nil {
			ss.Lifecycle.PostStart = &c.Lifecycle.PostStart.Handler
		}
		if c.Lifecycle.PreStop != nil {
			ss.Lifecycle.PreStop = &c.Lifecycle.PreStop.Handler
		}
	}

	if s.Resources != nil {
		resources := &ss.Resources

		for resourceType, r := range s.Resources {
			resources.Add(resourceType, r.RequestString(), r.LimitString())
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

	secretNames[pod.Image.ResolveImagePullSecret().SecretName()] = true

	for _, v := range pod.Initials {
		secretNames[v.Image.ResolveImagePullSecret().SecretName()] = true
	}

	ss := kubetypes.KubeImagePullSecrets{}

	for name := range secretNames {
		if name == "" {
			continue
		}
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

	for _, toleration := range s.Tolerations {
		t := kubetypes.KubeToleration{
			Key:    toleration.Key,
			Value:  toleration.Value,
			Effect: toleration.Effect,
		}

		if t.Value == "" {
			t.Operator = "Exists"
		} else {
			t.Operator = "Equal"
		}

		if toleration.TolerationSeconds != nil {
			t.TolerationSeconds = toleration.TolerationSeconds
		}

		kt.Tolerations = append(kt.Tolerations, t)
	}
	return kt
}

func ToKubeHosts(s spec.Spec) []kubetypes.KubeHosts {

	var ss []kubetypes.KubeHosts
	if s.Service != nil {
		for _, h := range s.Service.Hosts {
			ss = append(ss, kubetypes.KubeHosts{
				Ip:        h.Ip,
				HostNames: h.HostNames,
			})
		}
	}
	return ss
}

func ToKubeTopologySpreadConstraints(pod spec.Pod) kubetypes.KubeTopologySpreadConstraints {
	if pod.KubeTopologySpreadConstraints != nil {
		return *pod.KubeTopologySpreadConstraints
	}
	return kubetypes.KubeTopologySpreadConstraints{}
}

func ToKubeAffinity(pod spec.Pod) kubetypes.KubeAffinity {
	if pod.KubeAffinity != nil {
		return *pod.KubeAffinity
	}
	return kubetypes.KubeAffinity{}
}