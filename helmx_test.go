package helmx

import (
	"fmt"
	"github.com/go-courier/helmx/spec"
	"github.com/stretchr/testify/require"
	"os"
	"strconv"
	"strings"
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
   helmx/upstreams: {{ join .Upstreams ","}}
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
{{ spaces 6 | toYamlIndent ( toKubeTolerations . ) }}
{{ spaces 6 | toYamlIndent ( toKubeVolumes . ) }}
{{ spaces 6 | toYamlIndent ( toKubeInitContainers . ) }}
      containers:
      - 
{{ spaces 8 | toYamlIndent ( toKubeMainContainer . ) }}
`, func(s *spec.Spec) bool {

		if s.Resources.Cpu != "" {
			cpuResources := strings.Split(s.Resources.Cpu, "/")
			if len(cpuResources) >= 2 {
				cpuRequest, err := strconv.Atoi(cpuResources[0])
				if err != nil {
					return false
				}
				cpuLimit, err := strconv.Atoi(cpuResources[1])
				if err != nil {
					return false
				}
				if cpuRequest > cpuLimit {
					return false
				}
			} else {
				cpuRequest, err := strconv.Atoi(cpuResources[0])
				if err != nil {
					return false
				}
				if cpuRequest > 500 {
					return false
				}
				s.Resources.Cpu = fmt.Sprintf("%d/500", cpuRequest)
			}
		} else {
			s.Resources.Cpu = "10/500"
		}

		if s.Resources.Memory != "" {
			memoryResources := strings.Split(s.Resources.Memory, "/")
			if len(memoryResources) >= 2 {
				cpuRequest, err := strconv.Atoi(memoryResources[0])
				if err != nil {
					return false
				}
				cpuLimit, err := strconv.Atoi(memoryResources[1])
				if err != nil {
					return false
				}
				if cpuRequest > cpuLimit {
					return false
				}
			} else {
				memoryRequest, err := strconv.Atoi(memoryResources[0])
				if err != nil {
					return false
				}
				if memoryRequest > 500 {
					return false
				}
				s.Resources.Memory = fmt.Sprintf("%d/1024", memoryRequest)
			}
		} else {
			s.Resources.Memory = "10/500"
		}

		if s.Tolerations != nil {
			for _, t := range s.Tolerations {
				if len(strings.Split(t, "=")) != 2 {
					return false
				}
			}
		}

		return true
	})

	hx.ExecuteAll(os.Stdout, &hx.Spec)
}
