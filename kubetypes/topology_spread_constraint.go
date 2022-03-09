package kubetypes

// TopologySpreadConstraints
type KubeTopologySpreadConstraints struct {
	TopologySpreadConstraints []KubeTopologySpreadConstraint `json:"topologySpreadConstraints,omitempty" yaml:"topologySpreadConstraints,omitempty"`
}

type KubeTopologySpreadConstraint struct {
	MaxSkew           int32                         `json:"maxSkew,omitempty" yaml:"maxSkew,omitempty"`
	TopologyKey       string                        `json:"topologyKey,omitempty" yaml:"topologyKey,omitempty"`
	WhenUnsatisfiable UnsatisfiableConstraintAction `json:"whenUnsatisfiable,omitempty" yaml:"whenUnsatisfiable,omitempty"`
	LabelSelector     *LabelSelector                `json:"labelSelector,omitempty" yaml:"labelSelector,omitempty"`
}

type UnsatisfiableConstraintAction string

const (
	// DoNotSchedule instructs the scheduler not to schedule the pod
	// when constraints are not satisfied.
	DoNotSchedule UnsatisfiableConstraintAction = "DoNotSchedule"
	// ScheduleAnyway instructs the scheduler to schedule the pod
	// even if constraints are not satisfied.
	ScheduleAnyway UnsatisfiableConstraintAction = "ScheduleAnyway"
)

type LabelSelector struct {
	MatchLabels      map[string]string          `json:"matchLabels,omitempty" yaml:"matchLabels,omitempty"`
	MatchExpressions []LabelSelectorRequirement `json:"matchExpressions,omitempty" yaml:"matchExpressions,omitempty"`
}

type LabelSelectorRequirement struct {
	Key      string                `json:"key" yaml:"key"`
	Operator LabelSelectorOperator `json:"operator" yaml:"operator"`
	Values   []string              `json:"values,omitempty" yaml:"values,omitempty"`
}

type LabelSelectorOperator string

const (
	LabelSelectorOpIn           LabelSelectorOperator = "In"
	LabelSelectorOpNotIn        LabelSelectorOperator = "NotIn"
	LabelSelectorOpExists       LabelSelectorOperator = "Exists"
	LabelSelectorOpDoesNotExist LabelSelectorOperator = "DoesNotExist"
)
