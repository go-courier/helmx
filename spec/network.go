package spec

import (
    "fmt"
    "net/url"
    "strconv"
    "strings"

    "github.com/go-courier/helmx/constants"
)

func ParsePort(s string) (*Port, error) {
    if s == "" {
        return nil, fmt.Errorf("missing port value")
    }

    appProtocol := ""
    port := uint16(0)
    targetPort := uint16(0)
    protocol := ""
    isNodePort := false
    nodePort := uint16(0)

    parts := strings.Split(s, "/")
    s = parts[0]
    if len(parts) == 2 {
        protocol = strings.ToLower(parts[1])
    }

    if s[0] == '!' {
        isNodePort = true
        s = s[1:]
    }

    ports := strings.Split(s, ":")
    portStr := ports[0]
    appProtocolAndPort := strings.Split(portStr, "-")
    if len(appProtocolAndPort) == 2 {
        portStr = appProtocolAndPort[1]
        appProtocol = strings.ToLower(appProtocolAndPort[0])
    }

    p, err := strconv.ParseUint(portStr, 10, 16)
    if err != nil {
        return nil, fmt.Errorf("invalid port %v", ports[0])
    }
    if len(ports) == 2 {
        tp, err := strconv.ParseUint(ports[1], 10, 16)
        if err != nil {
            return nil, fmt.Errorf("invalid target port %v", ports[1])
        }
        targetPort = uint16(tp)

        if isNodePort {
            nodePort = uint16(p)
            port = uint16(tp)
            if nodePort < 20000 || nodePort > 40000 {
                return nil, fmt.Errorf("invalid value: %d: provided port is not in the valid range. The range of valid ports is 20000-40000", port)
            }
        } else {
            port = uint16(p)
        }

    } else {
        targetPort = uint16(p)
        port = targetPort
    }

    return &Port{
        AppProtocol:   appProtocol,
        NodePort:      nodePort,
        Port:          port,
        IsNodePort:    isNodePort,
        ContainerPort: targetPort,
        Protocol:      constants.Protocol(strings.ToUpper(protocol)),
    }, nil
}

// openapi:strfmt port
type Port struct {
    AppProtocol   string
    NodePort      uint16
    Port          uint16
    IsNodePort    bool
    ContainerPort uint16
    Protocol      constants.Protocol
}

func (s Port) String() string {
    v := ""
    if s.IsNodePort {
        v = "!"
    }

    if s.AppProtocol != "" {
        v += s.AppProtocol + "-"
    }

    if s.IsNodePort {
        if s.NodePort != 0 {
            v += strconv.FormatUint(uint64(s.NodePort), 10)
            if s.ContainerPort != 0 {
                v += ":" + strconv.FormatUint(uint64(s.ContainerPort), 10)
            }
        } else {
            if s.ContainerPort != 0 {
                v += strconv.FormatUint(uint64(s.ContainerPort), 10)
            }
        }

    } else {
        if s.Port != 0 {
            v += strconv.FormatUint(uint64(s.Port), 10)
        }
        if s.ContainerPort != 0 && s.ContainerPort != s.Port {
            v += ":" + strconv.FormatUint(uint64(s.ContainerPort), 10)
        }
    }

    if s.Protocol != "" {
        v += "/" + strings.ToLower(string(s.Protocol))
    }

    return v
}

func (s Port) MarshalText() ([]byte, error) {
    return []byte(s.String()), nil
}

func (s *Port) UnmarshalText(data []byte) error {
    servicePort, err := ParsePort(string(data))
    if err != nil {
        return err
    }
    *s = *servicePort
    return nil
}

func ParseIngressRule(s string) (*IngressRule, error) {
    if s == "" {
        return nil, fmt.Errorf("invalid ingress rule")
    }

    u, err := url.Parse(s)
    if err != nil {
        return nil, err
    }

    r := &IngressRule{
        Scheme: u.Scheme,
        Host:   u.Hostname(),
        Path:   u.Path,
    }

    if r.Scheme == "" {
        r.Scheme = "http"
    }

    p := u.Port()
    if p == "" {
        r.Port = 80
    } else {
        port, _ := strconv.ParseUint(p, 10, 16)
        r.Port = uint16(port)
    }

    return r, nil
}

// openapi:strfmt ingress-rule
type IngressRule struct {
    Scheme string
    Host   string
    Path   string
    Port   uint16
}

func (r IngressRule) String() string {
    if r.Scheme == "" {
        r.Scheme = "http"
    }
    if r.Port == 0 {
        r.Port = 80
    }

    return (&url.URL{
        Scheme: r.Scheme,
        Host:   r.Host + ":" + strconv.FormatUint(uint64(r.Port), 10),
        Path:   r.Path,
    }).String()
}

func (r IngressRule) MarshalText() ([]byte, error) {
    return []byte(r.String()), nil
}

func (r *IngressRule) UnmarshalText(data []byte) error {
    ir, err := ParseIngressRule(string(data))
    if err != nil {
        return err
    }
    *r = *ir
    return nil
}

// SecretName:host1,host2,host3
func ParseIngressTLS(s string) (*IngressTLS, error) {
    if s == "" {
        return nil, fmt.Errorf("invalid ingress tls")
    }

    tlsStr := strings.Split(s, ":")
    secretName := tlsStr[0]
    r := &IngressTLS{
        SecretName: secretName,
    }
    if len(tlsStr) > 1 {
        r.Hosts = strings.Split(tlsStr[1], ",")
    }
    return r, nil
}

// openapi:strfmt ingress-rule
type IngressTLS struct {
    Hosts      []string
    SecretName string
}

func (r IngressTLS) String() string {
    return fmt.Sprintf("%s:%s", r.SecretName, strings.Join(r.Hosts, ","))
}

func (r IngressTLS) MarshalText() ([]byte, error) {
    return []byte(r.String()), nil
}

func (r *IngressTLS) UnmarshalText(data []byte) error {
    ir, err := ParseIngressTLS(string(data))
    if err != nil {
        return err
    }
    *r = *ir
    return nil
}
