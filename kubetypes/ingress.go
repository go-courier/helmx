package kubetypes

type KubeIngressSpec struct {
	Rules []IngressRule `yaml:"rules"`
}

type IngressRule struct {
	Host string                `yaml:"host,omitempty"`
	HTTP *HTTPIngressRuleValue `yaml:"http,omitempty"`
}

type HTTPIngressRuleValue struct {
	Paths []HTTPIngressPath `yaml:"paths"`
}

type HTTPIngressPath struct {
	Path    string         `yaml:"path,omitempty"`
	Backend IngressBackend `yaml:"backend"`
}

type IngressBackend struct {
	ServiceName string `yaml:"serviceName"`
	ServicePort uint16 `yaml:"servicePort"`
}
