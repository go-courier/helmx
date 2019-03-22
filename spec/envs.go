package spec

type Envs map[string]string

func (envs Envs) Merge(srcEnvs Envs) Envs {
	es := Envs{}
	for k, v := range envs {
		es[k] = v
	}
	for k, v := range srcEnvs {
		es[k] = v
	}
	return es
}
