package spec

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
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

	ReadinessProbe                *Probe                                   `json:"readinessProbe,omitempty" yaml:"readinessProbe,omitempty"`
	LivenessProbe                 *Probe                                   `json:"livenessProbe,omitempty" yaml:"livenessProbe,omitempty"`
	Lifecycle                     *Lifecycle                               `json:"lifecycle,omitempty" yaml:"lifecycle,omitempty"`
	SecurityContext               *kubetypes.SecurityContext               `json:"securityContext,omitempty" yaml:"securityContext,omitempty"` // securityContext
	kubetypes.KubeTopologySpreadConstraints `yaml:",inline" json:",inline"`
	kubetypes.KubeAffinity                  `yaml:",inline" json:",inline"`
}

type Lifecycle struct {
	PostStart *Action `json:"postStart,omitempty" yaml:"postStart,omitempty"`
	PreStop   *Action `json:"preStop,omitempty" yaml:"preStop,omitempty"`
}

type Probe struct {
	Action              Action `json:"action" yaml:"action"`
	kubetypes.ProbeOpts `yaml:",inline"`
}

type Image struct {
	// default as project.group/project.name:version
	Tag string `json:"image,omitempty" yaml:"image,omitempty"`
	// <schema_name>://<host>/[prefix-]
	ImagePullSecret *ImagePullSecret     `json:"imagePullSecret,omitempty" yaml:"imagePullSecret,omitempty"`
	ImagePullPolicy constants.PullPolicy `json:"imagePullPolicy,omitempty" yaml:"imagePullPolicy,omitempty"`
}

func (i Image) ImageTag(defaultTag string) string {
	if i.Tag == "" {
		i.Tag = defaultTag
	}

	return i.ResolveImagePullSecret().PrefixTag(i.Tag)
}

func (i Image) ResolveImagePullSecret() *ImagePullSecret {
	imagePullSecret := i.ImagePullSecret

	if imagePullSecret == nil {
		v := os.Getenv(EnvKeyImagePullSecret)
		if v != "" {
			imagePullSecret, _ = ParseImagePullSecret(v)
		} else {
			imagePullSecret = &ImagePullSecret{}
		}
	}

	return imagePullSecret
}

// openapi:strfmt image-pull-secret
type ImagePullSecret struct {
	Name     string
	Host     string
	Username string
	Password string
	Prefix   string
}

func (s ImagePullSecret) SecretName() string {
	return s.Name
}

func (s ImagePullSecret) PrefixTag(tag string) string {
	if len(tag) > 0 && tag[0] == '~' {
		return s.Host + s.Prefix + tag[1:]
	}

	return tag
}

func ParseImagePullSecret(uri string) (*ImagePullSecret, error) {
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return nil, err
	}

	endpoint := &ImagePullSecret{}

	endpoint.Name = u.Scheme
	endpoint.Host = u.Host
	endpoint.Prefix = u.Path

	if u.User != nil {
		endpoint.Username = u.User.Username()
		endpoint.Password, _ = u.User.Password()
	}

	return endpoint, nil
}

func (s ImagePullSecret) String() string {
	v := &url.URL{}
	v.Scheme = s.Name
	v.Host = s.Host
	v.Path = s.Prefix

	if s.Username != "" || s.Password != "" {
		v.User = url.UserPassword(s.Username, s.Password)
	}

	return v.String()
}

func (s ImagePullSecret) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *ImagePullSecret) UnmarshalText(data []byte) error {
	imagePullSecret, err := ParseImagePullSecret(string(data))
	if err != nil {
		return err
	}
	*s = *imagePullSecret
	return nil
}

func (s ImagePullSecret) RegistryAuth() string {
	authConfig := AuthConfig{Username: s.Username, Password: s.Password, ServerAddress: s.Host}
	b, _ := json.Marshal(authConfig)
	return base64.StdEncoding.EncodeToString(b)
}

func (s ImagePullSecret) DockerConfigJSON() []byte {
	v := struct {
		Auths map[string]AuthConfig `json:"auths"`
	}{
		Auths: map[string]AuthConfig{
			s.Host: {Username: s.Username, Password: s.Password},
		},
	}
	b, _ := json.Marshal(v)
	return b
}

func (s ImagePullSecret) Base64EncodedDockerConfigJSON() string {
	return base64.StdEncoding.EncodeToString(s.DockerConfigJSON())
}

type AuthConfig struct {
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	ServerAddress string `json:"serveraddress,omitempty"`
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

	return a, nil
}

// openapi:strfmt action
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

// key=value:NoExecute,3600
// key:NoExecute
func ParseToleration(s string) (*Toleration, error) {
	if s == "" {
		return nil, nil
	}

	t := &Toleration{}

	parts := strings.Split(s, ":")

	kv := strings.Split(parts[0], "=")

	t.Key = kv[0]

	if len(kv) > 1 {
		t.Value = kv[1]
	}

	if len(parts) > 1 {
		effectAndDuration := strings.Split(parts[1], ",")
		t.Effect = effectAndDuration[0]

		if len(effectAndDuration) > 1 {
			d, err := strconv.ParseInt(effectAndDuration[1], 10, 64)
			if err != nil {
				return nil, errors.New("invalid toleration seconds")
			}
			t.TolerationSeconds = &d
		}
	}

	return t, nil
}

// openapi:strfmt toleration
type Toleration struct {
	Key               string
	Value             string
	Effect            string
	TolerationSeconds *int64
}

func (t *Toleration) UnmarshalText(text []byte) error {
	to, err := ParseToleration(string(text))
	if err != nil {
		return err
	}
	*t = *to
	return nil
}

func (t Toleration) MarshalText() (text []byte, err error) {
	return []byte(t.String()), nil
}

func (t Toleration) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(t.Key)

	if t.Value != "" {
		buf.WriteRune('=')
		buf.WriteString(t.Value)
	}

	if t.Effect != "" {
		buf.WriteRune(':')
		buf.WriteString(t.Effect)
	}

	if t.TolerationSeconds != nil {
		buf.WriteRune(',')
		buf.WriteString(strconv.FormatInt(int64(*t.TolerationSeconds), 10))
	}

	return buf.String()
}

// openapi:strfmt hosts
type Hosts struct {
	Ip        string   `yaml:"ip" json:"ip"`
	HostNames []string `yaml:"hostnames" json:"hostNames"`
}

// 127.0.0.1:test1.com,test2.com
func ParseHosts(s string) (*Hosts, error) {
	if s == "" {
		return nil, nil
	}

	t := &Hosts{}

	parts := strings.Split(s, ":")

	if len(parts) < 2 {
		return nil, nil
	}
	t.Ip = parts[0]
	kv := strings.Split(parts[1], ",")

	if len(kv) > 0 {
		t.HostNames = append(t.HostNames, kv...)
	}

	return t, nil
}

func (t *Hosts) UnmarshalText(text []byte) error {
	to, err := ParseHosts(string(text))
	if err != nil {
		return err
	}
	*t = *to
	return nil
}

func (t Hosts) MarshalText() (text []byte, err error) {
	return []byte(t.String()), nil
}

func (t Hosts) String() string {
	buf := bytes.NewBuffer(nil)
	buf.WriteString(t.Ip)
	buf.WriteString(":")

	if len(t.HostNames) != 0 {
		for index, host := range t.HostNames {
			buf.WriteString(host)
			if index != len(t.HostNames)-1 {
				buf.WriteRune(',')
			}
		}
	}
	return buf.String()
}
