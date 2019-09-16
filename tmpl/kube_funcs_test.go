package tmpl_test

import (
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
