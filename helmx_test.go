package helmx

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/go-courier/helmx/spec"
	"github.com/stretchr/testify/require"
)

func init() {
	os.Setenv(spec.EnvKeyImagePullSecret, "qcloud-registry://ccr.ccs.tencentyun.com/prefix-")
}

func Example() {
	hx := NewHelmX()
	hx.AddTemplate("ingress", ingress)
	hx.AddTemplate("service", service)
	hx.AddTemplate("deployment", deployment)
	hx.AddTemplate("job", job)
	hx.AddTemplate("cronJob", cronJob)

	hx.FromYAML([]byte(
		`
project:
  name: helmx
  feature: test
  group: helmx
  version: 0.0.0
  description: helmx

service:
  mounts:
    - "data:/usr/share/nginx:ro"
  ports:
    - "80:80"
  ingresses:
    - "http://helmx:80/helmx"
  initials:
    - image: dockercloud/hello-world
      mounts:
        - "data:/usr/share/nginx"
      command:
        - mv
      args:
        - /www
        - /usr/share/nginx/html

jobs:
  doonce:
    image: busybox
    backoffLimit: 4
  dosomecron:
    image: busybox
    cron:
      schedule: "*/1 * * * *"

envs:
  env: "test"

resources:
  cpu: 10/20
  memory: 20

tolerations:
  - env=test
  - project=test

volumes:
  data:
    emptyDir: {}

upstreams:
  - redis
  - mysql
`))

	hx.ExecuteAll(os.Stdout, &hx.Spec)
}

func TestTemplates(t *testing.T) {
	baseProject := `
project:
  name: helmx
  feature: test
  group: helmx
  version: 0.0.0
  description: helmx
`

	t.Run("service", func(t *testing.T) {
		check(t, baseProject+`
service:
  ports:
    - "80:80"
`,
			service,
			`
--- 

apiVersion: v1
kind: Service
metadata:
  name: helmx--test
  annotations: 
    helmx/project: >-
      {"name":"helmx","feature":"test","version":"0.0.0","group":"helmx","description":"helmx"}
    helmx/upstreams: ""
spec:
  selector:
    srv: helmx--test
  type: ClusterIP
  ports:
  - port: 80
    targetPort: 80
    protocol: TCP
`,
		)
	})

	t.Run("service with nodePort", func(t *testing.T) {
		check(t, baseProject+`
service:
  ports:
    - "!80:80"
`,
			service,
			`
--- 

apiVersion: v1
kind: Service
metadata:
  name: helmx--test
  annotations: 
    helmx/project: >-
      {"name":"helmx","feature":"test","version":"0.0.0","group":"helmx","description":"helmx"}
    helmx/upstreams: ""
spec:
  selector:
    srv: helmx--test
  type: NodePort
  ports:
  - nodePort: 80
    port: 80
    targetPort: 80
    protocol: TCP
`,
		)
	})

	t.Run("deployment", func(t *testing.T) {
		check(t, baseProject+`
service:
  ports:
    - "80:80"
`,
			deployment,
			`

---

apiVersion: apps/v1
kind: Deployment
metadata:
  name: helmx--test
  annotations: 
    helmx: "{\"project\":{\"name\":\"helmx\",\"feature\":\"test\",\"version\":\"0.0.0\",\"group\":\"helmx\",\"description\":\"helmx\"},\"service\":{\"ports\":[\"80\"]}}"
spec:
  selector:
    matchLabels:
      srv: helmx--test
  template:
    metadata:
      labels:
        srv: helmx--test
    spec:
      containers:
      - name: helmx--test
        image: ccr.ccs.tencentyun.com/prefix-helmx/helmx:0.0.0
        ports:
        - containerPort: 80
          protocol: TCP
      imagePullSecrets:
      - name: qcloud-registry
`,
		)
	})

	t.Run("job", func(t *testing.T) {
		check(t, baseProject+`
jobs:
  doonce:
    image: busybox
    restartPolicy: Never
    backoffLimit: 4
  docron:
    image: busybox
    restartPolicy: Never
    cron:
      schedule: "*/1 * * * *"
`,
			job,
			`
---

apiVersion: batch/v1
kind: Job
metadata:
  name: helmx--test--doonce
spec:
  backoffLimit: 4
  template:
    spec:
      containers:
      - name: helmx--test
        image: busybox
      imagePullSecrets:
      - name: qcloud-registry
      restartPolicy: Never
`,
		)
	})

	t.Run("cronJob", func(t *testing.T) {
		check(t, baseProject+`
jobs:
  doonce:
    image: busybox
    restartPolicy: Never
    backoffLimit: 4
  docron:
    image: busybox
    restartPolicy: Never
    cron:
      schedule: "*/1 * * * *"
`,
			cronJob,
			`
---

apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: helmx--test--docron
spec:
  schedule: '*/1 * * * *'
  jobTemplate:
    spec:
      template:
        spec:
          containers:
          - name: helmx--test
            image: busybox
          imagePullSecrets:
          - name: qcloud-registry
          restartPolicy: Never
`,
		)
	})
}

var (
	service = `
{{ if ( and ( exists .Service ) ( gt ( len .Service.Ports ) 0 ) ) }}
--- 

apiVersion: v1
kind: Service
metadata:
  name: {{ ( .Project.FullName ) }}
  annotations: 
    helmx/project: >-
      {{ toJson .Project }}
    helmx/upstreams: {{ join .Upstreams "," | quote }}
spec:
  selector:
    srv: {{ ( .Project.FullName ) }}
{{ spaces 2 | toYamlIndent ( toKubeServiceSpec . )  }}
{{ end }}
`
	deployment = `
{{ if ( exists .Service ) }}
---

apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ ( .Project.FullName ) }}
  annotations: 
    helmx: {{ toJson . | quote }}
spec:
  selector:
    matchLabels:
      srv: {{ ( .Project.FullName ) }}
{{ spaces 2 | toYamlIndent ( toKubeDeploymentSpec . )  }}
{{ end }}
`
	job = `
{{ $spec := .}}
{{ range $name, $job := .Jobs }}
{{ if (not (exists $job.Cron)) }}
---

apiVersion: batch/v1
kind: Job
metadata:
  name: {{ ( $spec.Project.FullName ) }}--{{ $name }}
spec:
{{ spaces 2 | toYamlIndent ( toKubeJobSpec $spec $job )  }}
{{ end }}
{{ end }}
`

	cronJob = `
{{ $spec := .}}
{{ range $name, $job := .Jobs }}
{{ if (exists $job.Cron) }}
---

apiVersion: batch/v1beta1
kind: CronJob
metadata:
  name: {{ ( $spec.Project.FullName ) }}--{{ $name }}
spec:
{{ spaces 2 | toYamlIndent ( toKubeCronJobSpec $spec $job )  }}
{{ end }}
{{ end }}
`

	ingress = `
{{ if ( len .Service.Ingresses ) }}
--- 

apiVersion: extensions/v1beta1
kind: Ingress
metadata:
  name: {{ ( .Project.FullName ) }}
  annotations:
    kubernetes.io/ingress.class: "nginx"
spec:
{{ spaces 2 | toYamlIndent ( toKubeIngressSpec . )}}
{{ end }}
`
)

func check(t *testing.T, helmx string, tmpl string, expect string) {
	hx := NewHelmX()
	err := hx.FromYAML([]byte(helmx))
	require.NoError(t, err)

	hx.AddTemplate("tmpl", tmpl)

	buf := &bytes.Buffer{}
	err = hx.ExecuteAll(buf, &hx.Spec)
	require.NoError(t, err)

	//fmt.Println(buf.String())

	require.Equal(t, strings.TrimSpace(expect), strings.TrimSpace(buf.String()))
}
