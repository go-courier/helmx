package helmx

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/go-courier/helmx/spec"
	"github.com/onsi/gomega"
)

func init() {
	os.Setenv(spec.EnvKeyImagePullSecret, "qcloud-registry://ccr.ccs.tencentyun.com/prefix-")
}

func Test(t *testing.T) {
	hx := NewHelmX()

	hx.AddTemplate("serviceAccount", serviceAccount)
	hx.AddTemplate("ingress", ingress)
	hx.AddTemplate("service", service)
	hx.AddTemplate("deployment", deployment)
	hx.AddTemplate("job", job)
	hx.AddTemplate("cronJob", cronJob)

	_ = hx.FromYAML([]byte(
		`
project:
  name: helmx
  feature: test
  group: helmx
  version: 0.0.0
  description: helmx

service:
  hostAliases:
    - "127.0.0.1:test1.com,test2.com"
    - "127.0.0.2:test2.com,test3.com"
  mounts:
    - "data:/usr/share/nginx:ro"
  ports:
    - "80:80"
    - "!20000:80"
  livenessProbe:
    action: "http://:80"
  lifecycle:
    preStop: "nginx -s quit"
  ingresses:
    - "http://helmx:80/helmx"
  serviceAccountName: test
  serviceAccountRoleRules:
    - secrets#get,update

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
  cpu: 10/20m
  memory: 0/20Mi
  nvidia.com/gpu: 0/20

tolerations:
  - env=test
  - project

volumes:
  data:
    emptyDir:
      medium: Memory
      sizeLimit: "1Gi"

upstreams:
  - redis
  - mysql

labels:
   testKey1: testValue1 
   testKey2: testValue2
`))

	if err := hx.ExecuteAll(os.Stdout, &hx.Spec); err != nil {
		panic(err)
	}
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
  - name: http-80
    port: 80
    targetPort: 80
    protocol: TCP
`,
		)
	})

	t.Run("service with nodePort", func(t *testing.T) {
		check(t, baseProject+`
service:
  ports:
    - "!20000:80"
    - "!80"
    - "!25000:25000"
    - "!40000:80"
    - "80:8080"
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
  - name: np-http-20000
    nodePort: 20000
    port: 20000
    targetPort: 80
    protocol: TCP
  - name: np-http-80
    port: 80
    targetPort: 80
    protocol: TCP
  - name: np-http-25000
    nodePort: 25000
    port: 25000
    targetPort: 25000
    protocol: TCP
  - name: np-http-40000
    nodePort: 40000
    port: 40000
    targetPort: 80
    protocol: TCP
  - name: http-80
    port: 80
    targetPort: 8080
    protocol: TCP
`,
		)
	})

	t.Run("deployment", func(t *testing.T) {
		check(t, baseProject+`
service:
  hosts:
    - "127.0.0.1:test1.com,test2.com"
    - "127.0.0.2:test3.com,test4.com"
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
  labels:
    app: helmx--test
    version: 0.0.0
  annotations:
    helmx: "{\"project\":{\"name\":\"helmx\",\"feature\":\"test\",\"version\":\"0.0.0\",\"group\":\"helmx\",\"description\":\"helmx\"},\"service\":{\"hosts\":[\"127.0.0.1:test1.com,test2.com\",\"127.0.0.2:test3.com,test4.com\"],\"ports\":[\"80\"]}}"
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
      hostAliases:
      - ip: 127.0.0.1
        hostnames:
        - test1.com
        - test2.com
      - ip: 127.0.0.2
        hostnames:
        - test3.com
        - test4.com
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
  labels:
    app: {{ ( .Project.FullName ) }}
    version: {{ ( .Project.Version ) }}
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
	serviceAccount = `
{{ if ( len .Service.ServiceAccountRoleRules ) }}

--- 
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: {{ ( .Service.ServiceAccountName ) }}
rules:
{{ spaces 2 | toYamlIndent ( toKubeRoleRules . )}}

---

apiVersion: v1
kind: ServiceAccount
metadata:
  name: {{ ( .Service.ServiceAccountName ) }}

---
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: {{ ( .Service.ServiceAccountName ) }}
subjects:
  - kind: ServiceAccount
    name: {{ ( .Service.ServiceAccountName ) }}
    apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: Role
  name: {{ ( .Service.ServiceAccountName ) }}
  apiGroup: rbac.authorization.k8s.io

{{ end }}
`
)

func check(t *testing.T, helmx string, tmpl string, expect string) {
	hx := NewHelmX()
	err := hx.FromYAML([]byte(helmx))
	gomega.NewWithT(t).Expect(err).To(gomega.BeNil())

	hx.AddTemplate("tmpl", tmpl)

	buf := &bytes.Buffer{}
	err = hx.ExecuteAll(buf, &hx.Spec)
	gomega.NewWithT(t).Expect(err).To(gomega.BeNil())

	gomega.NewWithT(t).Expect(strings.TrimSpace(buf.String())).To(gomega.Equal(strings.TrimSpace(expect)))
}
