package tmpl_test

import (
	"encoding/json"
	"fmt"
	"github.com/go-courier/helmx/kubetypes"
	"gopkg.in/yaml.v2"
	"testing"

	"github.com/davecgh/go-spew/spew"
	"github.com/go-courier/helmx/spec"

	"github.com/go-courier/helmx/tmpl"
)

func TestToKubeJobSpec(t *testing.T) {
	s := spec.Spec{}
	s.Project = &spec.Project{Name: "test"}

	job := spec.Job{}
	job.Image.Tag = "busybox"

	s.Jobs = map[string]spec.Job{
		"do-once": job,
	}

	js := tmpl.ToKubeJobSpec(s, s.Jobs["do-once"])
	spew.Dump(js)
}

func TestKubeType(t *testing.T) {
	t.Run("#KubeTopologySpreadConstraints", func(t *testing.T) {
		ymlStr := `
topologySpreadConstraints:
  - maxSkew: 1
    topologyKey: zone
    whenUnsatisfiable: DoNotSchedule
    labelSelector:
      matchLabels:
        foo: bar
`
		obj := new(kubetypes.KubeTopologySpreadConstraints)
		if err := yaml.Unmarshal([]byte(ymlStr), obj); err != nil {
			t.Fatal(err)
		}

		pod := spec.Pod{}
		pod.KubeTopologySpreadConstraints = *obj

		js := tmpl.ToKubeTopologySpreadConstraints(pod)
		spew.Dump(js)

		fmt.Println(obj)
	})

	t.Run("#KubeAffinity", func(t *testing.T) {
		ymlStr := `
affinity:
    nodeAffinity:
      requiredDuringSchedulingIgnoredDuringExecution:
        nodeSelectorTerms:
        - matchExpressions:
          - key: zone
            operator: NotIn
            values:
            - zoneC
`
		obj := new(kubetypes.KubeAffinity)
		if err := yaml.Unmarshal([]byte(ymlStr), obj); err != nil {
			t.Fatal(err)
		}

		pod := spec.Pod{}
		pod.KubeAffinity = *obj

		js := tmpl.ToKubeAffinity(pod)
		spew.Dump(js)
	})
}

func LogJSON(v interface{}) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}