package spec

import (
	"net/url"
	"os"
	"strconv"
	"strings"

	"github.com/go-courier/helmx/constants"
	"github.com/go-courier/helmx/kubetypes"
)

var (
	EnvKeyImagePullSecret = "IMAGE_PULL_SECRET"
)

type Container struct {
	Image      `yaml:",inline"`
	WorkingDir string        `json:"workingDir,omitempty" yaml:"workingDir,omitempty"`
	Command    []string      `json:"command,omitempty" yaml:"command,omitempty"`
	Args       []string      `json:"args,omitempty" yaml:"args,omitempty"`
	Mounts     []VolumeMount `json:"mounts,omitempty" yaml:"mounts,omitempty"`
	Envs       Envs          `json:"envs,omitempty" yaml:"envs,omitempty"`
	TTY        bool          `json:"tty,omitempty" yaml:"tty,omitempty"`

	ReadinessProbe *Probe     `json:"readinessProbe,omitempty" yaml:"readinessProbe,omitempty"`
	LivenessProbe  *Probe     `json:"livenessProbe,omitempty" yaml:"livenessProbe,omitempty"`
	Lifecycle      *Lifecycle `json:"lifecycle,omitempty" yaml:"lifecycle,omitempty"`
}

type Lifecycle struct {
	PostStart Action `json:"postStart,omitempty" yaml:"postStart,omitempty"`
	PreStop   Action `json:"preStop,omitempty" yaml:"preStop,omitempty"`
}

type Probe struct {
	Action              Action `json:"action" yaml:"action"`
	kubetypes.ProbeOpts `yaml:",inline"`
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

// http://:80
// tcp://:80
// exec
func ParseAction(s string) (*Action, error) {
	if s == "" {
		return nil, nil
	}

	a := &Action{}

	if strings.HasPrefix(s, "http") || strings.HasPrefix(s, "tcp") {
		u, err := url.Parse(s)
		if err != nil {
			return nil, err
		}

		port, _ := strconv.ParseUint(u.Port(), 10, 64)

		if u.Scheme == "tcp" {
			a.TCPSocket = &kubetypes.TCPSocketAction{}
			a.TCPSocket.Host = u.Hostname()
			a.TCPSocket.Port = uint16(port)
			return a, nil
		}

		a.HTTPGet = &kubetypes.HTTPGetAction{}
		a.HTTPGet.Port = uint16(port)
		a.HTTPGet.Host = u.Hostname()
		a.HTTPGet.Path = u.Path
		a.HTTPGet.Scheme = strings.ToUpper(u.Scheme)

		return a, nil
	}

	a.Exec = &kubetypes.ExecAction{
		Command: []string{"sh", "-c", s},
	}

	return nil, nil
}

type Action struct {
	kubetypes.Handler
}

func (a Action) String() string {
	if a.Exec != nil {
		return a.Exec.Command[2]
	}

	if a.HTTPGet != nil {
		u := &url.URL{}
		u.Scheme = strings.ToLower(a.HTTPGet.Scheme)
		u.Path = a.HTTPGet.Path
		u.Host = a.HTTPGet.Host + ":" + strconv.FormatUint(uint64(a.HTTPGet.Port), 10)

		if u.Scheme != "" {
			u.Scheme = "http"
		}
		return u.String()
	}

	if a.TCPSocket != nil {
		u := &url.URL{}
		u.Scheme = "tcp"
		u.Host = a.TCPSocket.Host + ":" + strconv.FormatUint(uint64(a.TCPSocket.Port), 10)

		return u.String()
	}

	return ""
}

func (a Action) MarshalText() ([]byte, error) {
	return []byte(a.String()), nil
}

func (a *Action) UnmarshalText(data []byte) error {
	action, err := ParseAction(string(data))
	if err != nil {
		return err
	}
	if action != nil {
		*a = *action
	}
	return nil
}
