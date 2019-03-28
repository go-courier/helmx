package helmx

import (
	"fmt"
	"github.com/go-courier/helmx/spec"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func init() {
	os.Chdir("./testdata")
	os.Setenv(spec.EnvKeyImagePullSecret, "qcloud-registry://ccr.ccs.tencentyun.com/prefix-")
}

func TestHelmXFromWorkingDir(t *testing.T) {
	hx := NewHelmX()

	hlemXData := MustReadOrFetch("./helmx.yml")

	err := hx.FromYAML(hlemXData)
	require.NoError(t, err)

	fmt.Println(string(MustBytes(hx.ToYAML())))
}

func TestTemplates(t *testing.T) {
	hx := NewHelmX()

	hlemXData := MustReadOrFetch("./helmx.yml")

	err := hx.FromYAML(hlemXData)
	require.NoError(t, err)

	hx.AddTemplate(
		"ingress",
		`
apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: {{ ( .Project.FullName ) }}
  annotations:
    kubernetes.io/ingress.class: "nginx"
spec:
{{ spaces 2 | toYamlIndent ( toKubeIngressRules . )}}
`,
		func(s *spec.Spec) bool {
			return len(s.Service.Ingress) > 0
		},
	)

	hx.AddTemplate("service", `
apiVersion: v1
kind: Service
metadata:
  name: {{ ( .Project.FullName ) }}
  annotations: 
    helmx/project: {{ toJson .Project  | quote }}
    helmx/upstreams: {{ join .Upstreams "," | quote }}
spec:
  selector:
    srv: {{ ( .Project.FullName ) }}
{{ spaces 2 | toYamlIndent ( toKubeServicePorts . )  }}
`)

	hx.AddTemplate("deployment", `
apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ ( .Project.FullName ) }}
  annotations: 
    helmx: {{ toJson . | quote }}
spec:
  replicas: {{ .Service.Replicas | default 1 }}
  selector:
    matchLabels:
      srv: {{ ( .Project.FullName ) }}
  template:
    metadata:
      labels:
        srv: {{ ( .Project.FullName ) }}
    spec:
{{ spaces 6 | toYamlIndent ( toKubeImagePullSecrets . ) }}
{{ spaces 6 | toYamlIndent ( toKubeTolerations . ) }}
{{ spaces 6 | toYamlIndent ( toKubeVolumes . ) }}
{{ spaces 6 | toYamlIndent ( toKubeInitContainers . ) }}
      containers:
      - 
{{ spaces 8 | toYamlIndent ( toKubeMainContainer . ) }}
`)

	hx.ExecuteAll(os.Stdout, &hx.Spec)
}
