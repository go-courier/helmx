package kubetypes

type KubeAffinity struct {
	Affinity *Affinity `json:"affinity,omitempty" yaml:"affinity,omitempty"`
}

type Affinity struct {
	NodeAffinity    *NodeAffinity    `json:"nodeAffinity,omitempty" yaml:"nodeAffinity,omitempty"`
	PodAffinity     *PodAffinity     `json:"podAffinity,omitempty" yaml:"podAffinity,omitempty"`
	PodAntiAffinity *PodAntiAffinity `json:"podAntiAffinity,omitempty" yaml:"podAntiAffinity,omitempty"`
}

type NodeAffinity struct {
	RequiredDuringSchedulingIgnoredDuringExecution *NodeSelector `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty" yaml:"requiredDuringSchedulingIgnoredDuringExecution,omitempty"`

	PreferredDuringSchedulingIgnoredDuringExecution []PreferredSchedulingTerm `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty" yaml:"preferredDuringSchedulingIgnoredDuringExecution,omitempty"`
}

type PodAffinity struct {
	RequiredDuringSchedulingIgnoredDuringExecution  []PodAffinityTerm         `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,1,rep,name=requiredDuringSchedulingIgnoredDuringExecution"`
	PreferredDuringSchedulingIgnoredDuringExecution []WeightedPodAffinityTerm `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,2,rep,name=preferredDuringSchedulingIgnoredDuringExecution"`
}

type PodAntiAffinity struct {
	RequiredDuringSchedulingIgnoredDuringExecution  []PodAffinityTerm         `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty" yaml:"requiredDuringSchedulingIgnoredDuringExecution,omitempty"`
	PreferredDuringSchedulingIgnoredDuringExecution []WeightedPodAffinityTerm `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty" yaml:"preferredDuringSchedulingIgnoredDuringExecution,omitempty"`
}

type WeightedPodAffinityTerm struct {
	Weight          int32           `json:"weight" yaml:"weight"`
	PodAffinityTerm PodAffinityTerm `json:"podAffinityTerm" yaml:"podAffinityTerm"`
}

type PodAffinityTerm struct {
	LabelSelector *LabelSelector `json:"labelSelector,omitempty" yaml:"labelSelector,omitempty"`
	Namespaces    []string       `json:"namespaces,omitempty" yaml:"namespaces,omitempty"`
	TopologyKey   string         `json:"topologyKey" yaml:"topologyKey"`
}

type PreferredSchedulingTerm struct {
	Weight     int32            `json:"weight" yaml:"weight"`
	Preference NodeSelectorTerm `json:"preference" yaml:"preference"`
}

type NodeSelector struct {
	NodeSelectorTerms []NodeSelectorTerm `json:"nodeSelectorTerms" yaml:"nodeSelectorTerms"`
}

type NodeSelectorTerm struct {
	MatchExpressions []NodeSelectorRequirement `json:"matchExpressions,omitempty" yaml:"matchExpressions,omitempty"`
	MatchFields      []NodeSelectorRequirement `json:"matchFields,omitempty" yaml:"matchFields,omitempty"`
}

type NodeSelectorRequirement struct {
	Key      string               `json:"key" yaml:"key"`
	Operator NodeSelectorOperator `json:"operator" yaml:"operator"`
	Values   []string             `json:"values,omitempty" yaml:"values,omitempty"`
}

type NodeSelectorOperator string

const (
	NodeSelectorOpIn           NodeSelectorOperator = "In"
	NodeSelectorOpNotIn        NodeSelectorOperator = "NotIn"
	NodeSelectorOpExists       NodeSelectorOperator = "Exists"
	NodeSelectorOpDoesNotExist NodeSelectorOperator = "DoesNotExist"
	NodeSelectorOpGt           NodeSelectorOperator = "Gt"
	NodeSelectorOpLt           NodeSelectorOperator = "Lt"
)
