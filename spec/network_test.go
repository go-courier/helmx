package spec

import (
    "testing"

    "github.com/go-courier/helmx/constants"
    "github.com/stretchr/testify/require"
    "gopkg.in/yaml.v2"
)

func TestIngressTLS(t *testing.T) {
    t.Run("parse & string", func(t *testing.T) {
        r, _ := ParseIngressTLS("secretName:host1,host2,host3")
        require.Equal(t, "secretName", r.SecretName)
        require.Equal(t, "host1", r.Hosts[0])
        require.Equal(t, "host2", r.Hosts[1])
        require.Equal(t, "host3", r.Hosts[2])
    })

    t.Run("yaml marshal & unmarshal", func(t *testing.T) {
        data, err := yaml.Marshal(struct {
            IngressTLS IngressTLS `yaml:"ingress"`
        }{
            IngressTLS: IngressTLS{
                SecretName: "secretName",
                Hosts:      []string{"host1", "host2", "host3"},
            },
        })
        require.NoError(t, err)
        require.Equal(t, "ingress: secretName:host1,host2,host3\n", string(data))

        v := struct {
            IngressTLS IngressTLS `yaml:"ingress"`
        }{}

        err = yaml.Unmarshal(data, &v)
        require.NoError(t, err)
        require.Equal(t, "secretName:host1,host2,host3", v.IngressTLS.String())
    })
}

func TestIngressRule(t *testing.T) {
    t.Run("parse & string", func(t *testing.T) {
        r, _ := ParseIngressRule("http://*.helmx/helmx")

        require.Equal(t, uint16(80), r.Port)
        require.Equal(t, "*.helmx", r.Host)
        require.Equal(t, "http", r.Scheme)

        require.Equal(t, "http://*.helmx:80/helmx", r.String())
    })

    t.Run("yaml marshal & unmarshal", func(t *testing.T) {
        data, err := yaml.Marshal(struct {
            IngressRule IngressRule `yaml:"ingress"`
        }{
            IngressRule: IngressRule{
                Port: 80,
                Host: "*.helmx",
                Path: "/helmx",
            },
        })
        require.NoError(t, err)
        require.Equal(t, "ingress: http://*.helmx:80/helmx\n", string(data))

        v := struct {
            IngressRule IngressRule `yaml:"ingress"`
        }{}

        err = yaml.Unmarshal(data, &v)
        require.NoError(t, err)
        require.Equal(t, "http://*.helmx:80/helmx", v.IngressRule.String())
    })
}

func TestPort(t *testing.T) {
    t.Run("parse & string", func(t *testing.T) {
        sp, _ := ParsePort("http-80:8080/TCP")

        require.Equal(t, "http", sp.AppProtocol)
        require.Equal(t, uint16(80), sp.Port)
        require.Equal(t, uint16(8080), sp.ContainerPort)
        require.Equal(t, constants.ProtocolTCP, sp.Protocol)

        require.Equal(t, "http-80:8080/tcp", sp.String())
    })

    t.Run("parse & string without target port ", func(t *testing.T) {
        sp, _ := ParsePort("80/TCP")
        require.Equal(t, uint16(80), sp.Port)
        require.Equal(t, uint16(80), sp.ContainerPort)
        require.Equal(t, constants.ProtocolTCP, sp.Protocol)

        require.Equal(t, "80/tcp", sp.String())
    })

    t.Run("parse & string without protocol ", func(t *testing.T) {
        sp, _ := ParsePort("80")
        require.Equal(t, uint16(80), sp.Port)
        require.Equal(t, uint16(80), sp.ContainerPort)
    })

    t.Run("parse & string without node port", func(t *testing.T) {
        sp, _ := ParsePort("!20000:8080")
        require.Equal(t, true, sp.IsNodePort)
        require.Equal(t, uint16(20000), sp.NodePort)
        require.Equal(t, uint16(8080), sp.Port)
        require.Equal(t, uint16(8080), sp.ContainerPort)
        require.Equal(t, "!20000:8080", sp.String())
    })

    t.Run("parse & string without node port prefix", func(t *testing.T) {
        sp, err := ParsePort("!8080")
        require.NoError(t, err)
        require.Equal(t, uint16(8080), sp.Port)
        require.Equal(t, uint16(8080), sp.ContainerPort)
        require.Equal(t, "!8080", sp.String())
    })

    t.Run("parse & string without protocol", func(t *testing.T) {
        sp, _ := ParsePort("80:8080")
        require.Equal(t, uint16(80), sp.Port)
        require.Equal(t, uint16(8080), sp.ContainerPort)

        require.Equal(t, "80:8080", sp.String())
    })

    t.Run("yaml marshal & unmarshal", func(t *testing.T) {
        data, err := yaml.Marshal(struct {
            Port Port `yaml:"port"`
        }{
            Port: Port{
                Port:          80,
                ContainerPort: 8080,
                Protocol:      "TCP",
            },
        })
        require.NoError(t, err)
        require.Equal(t, "port: 80:8080/tcp\n", string(data))

        v := struct {
            Port Port `yaml:"port"`
        }{}

        err = yaml.Unmarshal(data, &v)
        require.NoError(t, err)
        require.Equal(t, "80:8080/tcp", v.Port.String())
    })

    t.Run("node port range in 20000-40000", func(t *testing.T) {
        _, ltErr := ParsePort("!19999:80")
        _, noLtErr := ParsePort("!20000:80")
        _, gtErr := ParsePort("!40001:80")
        _, noGtErr := ParsePort("!40000:80")

        require.Error(t, ltErr)
        require.NoError(t, noLtErr)
        require.Error(t, gtErr)
        require.NoError(t, noGtErr)
    })

}
