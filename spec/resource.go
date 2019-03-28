package spec

import (
	"fmt"
	"strconv"
	"strings"
)

type Resource struct {
	Cpu    RequestAndLimit `yaml:"cpu,omitempty"`
	Memory RequestAndLimit `yaml:"memory,omitempty"`
}

func ParseRequestAndLimit(s string) (*RequestAndLimit, error) {
	if s == "" {
		return nil, fmt.Errorf("missing request and limit")
	}
	parts := strings.Split(s, "/")

	rl := &RequestAndLimit{}

	i, err := strconv.ParseInt(parts[0], 10, 64)
	if err == nil {
		rl.Request = int(i)
	}

	if len(parts) == 2 {
		i, err := strconv.ParseInt(parts[1], 10, 64)
		if err == nil {
			rl.Limit = int(i)
		}
	}

	return rl, nil
}

type RequestAndLimit struct {
	Request int
	Limit   int
}

func (s RequestAndLimit) String() string {
	v := ""
	if s.Request != 0 {
		v = strconv.FormatInt(int64(s.Request), 10)
	}
	if s.Limit != 0 {
		v = "/" + strconv.FormatInt(int64(s.Limit), 10)
	}
	return v
}

func (s RequestAndLimit) MarshalText() ([]byte, error) {
	return []byte(s.String()), nil
}

func (s *RequestAndLimit) UnmarshalText(data []byte) error {
	servicePort, err := ParseRequestAndLimit(string(data))
	if err != nil {
		return err
	}
	*s = *servicePort
	return nil
}
