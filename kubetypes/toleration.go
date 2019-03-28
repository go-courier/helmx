package kubetypes

type TolerationEffect string
type TolerationOperator string

const (
	TolerationEffectNoExecute TolerationEffect   = "NoExecute"
	TolerationOperatorEqual   TolerationOperator = "Equal"
)

type KubeTolerations struct {
	Tolerations []KubeToleration `yaml:"tolerations,omitempty"`
}

type KubeToleration struct {
	Key      string             `yaml:"key"`
	Value    string             `yaml:"value"`
	Effect   TolerationEffect   `yaml:"effect"`
	Operator TolerationOperator `yaml:"operator"`
}
