package spec

type NameSpace struct {
	Name         string `env:"name" yaml:"name" json:"name"`
	NodeSelector string `env:"nodeSelector" yaml:"nodeSelector" json:"nodeSelector"`
}
