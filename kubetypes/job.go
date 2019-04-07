package kubetypes

type KubeJobSpec struct {
	JobOpts  `yaml:",inline"`
	Template struct {
		KubeMetadata `yaml:",inline"`
		Spec         KubePodSpec `yaml:"spec"`
	} `yaml:"template"`
}

type JobOpts struct {
	Parallelism             *int32 `yaml:"parallelism,omitempty" json:"parallelism,omitempty"`
	Completions             *int32 `yaml:"completions,omitempty" json:"completions,omitempty"`
	BackoffLimit            *int32 `yaml:"backoffLimit,omitempty" json:"backoffLimit,omitempty"`
	TTLSecondsAfterFinished *int32 `yaml:"ttlSecondsAfterFinished,omitempty" json:"ttlSecondsAfterFinished,omitempty"`
}

type KubeCronJobSpec struct {
	CronJobOpts `yaml:",inline"`
	Template    struct {
		KubeMetadata `yaml:",inline"`
		Spec         KubeJobSpec `yaml:"spec"`
	} `yaml:"jobTemplate"`
}

type CronJobOpts struct {
	Schedule                   string `yaml:"schedule" json:"schedule"`
	StartingDeadlineSeconds    *int64 `yaml:"startingDeadlineSeconds,omitempty" json:"startingDeadlineSeconds,omitempty"`
	ConcurrencyPolicy          string `yaml:"concurrencyPolicy,omitempty" json:"concurrencyPolicy,omitempty"`
	Suspend                    *bool  `yaml:"suspend,omitempty" json:"suspend,omitempty"`
	SuccessfulJobsHistoryLimit *int32 `yaml:"successfulJobsHistoryLimit,omitempty" json:"successfulJobsHistoryLimit,omitempty"`
	FailedJobsHistoryLimit     *int32 `yaml:"failedJobsHistoryLimit,omitempty" json:"failedJobsHistoryLimit,omitempty"`
}
