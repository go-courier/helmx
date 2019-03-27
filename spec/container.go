package spec

import (
	"github.com/go-courier/helmx/constants"
	"github.com/go-courier/helmx/kubetypes"
	"net/url"
	"os"
	"strings"
)

var (
	EnvKeyImagePullSecret = "IMAGE_PULL_SECRET"
)

type Container struct {
	Image      `yaml:",inline"`
	WorkingDir string                  `json:"workingDir,omitempty" yaml:"workingDir,omitempty"`
	Command    []string                `json:"command,omitempty" yaml:"command,omitempty"`
	Args       []string                `json:"args,omitempty" yaml:"args,omitempty"`
	Mounts     []VolumeMount           `json:"mounts,omitempty" yaml:"mounts,omitempty"`
	Envs       Envs                    `json:"envs,omitempty" yaml:"envs,omitempty"`
	TTY        bool                    `json:"tty,omitempty" yaml:"tty,omitempty"`
	Resources  kubetypes.KubeResources `json:"resources,omitempty" yaml:"resources,omitempty"`
}

type Image struct {
	// default as project.group/project.name:version
	Tag             string               `json:"image,omitempty" yaml:"image,omitempty"`
	ImagePullPolicy constants.PullPolicy `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
	ImagePullSecret *ImagePullSecret     `json:"imagePullSecret,omitempty" yaml:"imagePullSecret,omitempty"`
}

func (i Image) ImageTag(defaultTag string) string {
	if i.Tag == "" {
		i.Tag = defaultTag
	}

	if len(i.Tag) > 0 && i.Tag[0] == '~' {
		if i.ImagePullSecret == nil {
			i.ImagePullSecret = &ImagePullSecret{}
			i.ImagePullSecret.init()
		}

		i.Tag = i.Tag[1:]
	}

	if i.ImagePullSecret != nil {
		return i.ImagePullSecret.PrefixTag(i.Tag)
	}

	return i.Tag
}

type ImagePullSecret struct {
	Name   string `json:"name" yaml:"name"`
	Host   string `json:"host" yaml:"host"`
	Prefix string `json:"prefix,omitempty" yaml:"prefix,omitempty"`
}

func (s ImagePullSecret) SecretName() string {
	s.init()
	return s.Name
}

func (s ImagePullSecret) PrefixTag(tag string) string {
	s.init()
	return s.Host + "/" + s.Prefix + tag
}

func (s *ImagePullSecret) init() {
	if s.Host == "" {
		imagePullSecret := os.Getenv(EnvKeyImagePullSecret)
		if imagePullSecret != "" {
			u, err := url.Parse(imagePullSecret)
			if err != nil {
				panic(err)
			}
			s.Host = u.Host
			s.Name = u.Scheme
			s.Prefix = strings.TrimLeft(u.Path, "/")
		}
	}
}
