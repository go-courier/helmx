package kubetypes

type ResourceUnit string

const (
	ResourceUnitRequestCpu    ResourceUnit = "10m"
	ResourceUnitRequestMemory ResourceUnit = "10Mi"
	ResourceUnitLimitCpu      ResourceUnit = "500m"
	ResourceUnitLimitMemory   ResourceUnit = "1024Mi"
)

type KubeResources struct {
	Requests Resource `yaml:"requests,omitempty"` // default  10m 10Mi
	Limits   Resource `yaml:"limits,omitempty"`   // default  500m 1024Mi
}

type Resource struct {
	Cpu    ResourceUnit `yaml:"cpu,omitempty"`
	Memory ResourceUnit `yaml:"memory,omitempty"`
}
