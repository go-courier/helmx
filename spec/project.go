package spec

type Project struct {
	Name        string  `env:"name" yaml:"name" json:"name"`
	Feature     string  `env:"feature" yaml:"feature,omitempty" json:"feature,omitempty"`
	Version     Version `env:"version" yaml:"version" json:"version"`
	Group       string  `env:"group" yaml:"group,omitempty" json:"group,omitempty"`
	Description string  `env:"description" yaml:"description,omitempty" json:"description,omitempty"`
}

func (p Project) FullName() string {
	if p.Feature != "" {
		return p.Name + "--" + p.Feature
	}
	return p.Name
}

func (p Project) DefaultImageTag() string {
	return "~" + p.Group + "/" + p.Name + ":" + p.Version.String()
}
