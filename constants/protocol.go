package constants

import (
	"fmt"
)

type Protocol string

const (
	ProtocolTCP  Protocol = "TCP"
	ProtocolUDP  Protocol = "UDP"
	ProtocolSCTP Protocol = "SCTP"
)

func (p *Protocol) UnmarshalText(text []byte) error {
	pp := Protocol(text)
	switch pp {
	case "":
		return nil
	case ProtocolTCP, ProtocolUDP, ProtocolSCTP:
		*p = pp
		return nil
	}
	return fmt.Errorf("unsupported protocol %s", pp)
}
