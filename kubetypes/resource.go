package kubetypes

type ResourceRequirements struct {
	Requests map[string]string `yaml:"requests,omitempty"`
	Limits   map[string]string `yaml:"limits,omitempty"`
}

func (r *ResourceRequirements) Add(tpe string, request string, limit string) {
	if request != "" {
		if r.Requests == nil {
			r.Requests = map[string]string{}
		}

		r.Requests[tpe] = request
	}

	if limit != "" {
		if r.Limits == nil {
			r.Limits = map[string]string{}
		}

		r.Limits[tpe] = limit
	}
}
