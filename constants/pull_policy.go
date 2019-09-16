package constants

import (
	"fmt"
)

type PullPolicy string

const (
	PullAlways       PullPolicy = "Always"
	PullNever        PullPolicy = "Never"
	PullIfNotPresent PullPolicy = "IfNotPresent"
)

func (p *PullPolicy) UnmarshalText(text []byte) error {
	pp := PullPolicy(text)
	switch pp {
	case "":
		return nil
	case PullAlways, PullNever, PullIfNotPresent:
		*p = pp
		return nil
	}
	return fmt.Errorf("unsupported pull policy %s", pp)
}
