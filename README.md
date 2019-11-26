## HelmX

[![GoDoc Widget](https://godoc.org/github.com/go-courier/helmx?status.svg)](https://godoc.org/github.com/go-courier/helmx)
[![Build Status](https://travis-ci.org/go-courier/helmx.svg?branch=master)](https://travis-ci.org/go-courier/helmx)
[![codecov](https://codecov.io/gh/go-courier/helmx/branch/master/graph/badge.svg)](https://codecov.io/gh/go-courier/helmx)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-courier/helmx)](https://goreportcard.com/report/github.com/go-courier/helmx)


Deploy k8s on flying.


```yaml
# Example

project:
  name: helmx
  feature: test
  group: helmx
  version: 0.0.0
  description: helmx

service:
  hosts:
    - "127.0.0.1:test1.com,test2.com"
    - "127.0.0.2:test3.com,test4.com"
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
  - project

volumes:
  data:
    emptyDir: {}

upstreams:
  - redis
  - mysql

labels:
   testKey1: testValue1 
   testKey2: testValue2
```
