package spec

import (
	"github.com/go-courier/helmx/kubetypes"
)

type Job struct {
	Pod               `yaml:",inline"`
	kubetypes.JobOpts `yaml:",inline"`
	Cron              *kubetypes.CronJobOpts `json:"cron,omitempty" yaml:"cron,omitempty"`
}
