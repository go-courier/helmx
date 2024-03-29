package spec

import (
    "fmt"
    "regexp"
    "strings"
)

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

type EnvsWithValueFrom map[string]*EnvValue

func ParseEnvsWithValueFrom(envMap Envs) (EnvsWithValueFrom, error) {
    envs := make(EnvsWithValueFrom)
    for k, v := range envMap {
        envValue, err := ParseEnvValue(v)
        if err != nil {
            return nil, err
        }
        envs[k] = envValue
    }
    return envs, nil
}

type EnvValueFromConfigMap struct {
    ConfigMapName string `json:"configMapName" yaml:"configMapName"`
    Key           string `json:"key" yaml:"key"`
}

func (c *EnvValueFromConfigMap) String() string {
    return fmt.Sprintf("%s|%s", c.ConfigMapName, c.Key)
}

type EnvValueFromSecret struct {
    SecretName string `json:"secretName" yaml:"secretName"`
    Key        string `json:"key"`
    Optional   bool   `json:"optional,omitempty" yaml:"optional,omitempty"`
}

func (c *EnvValueFromSecret) String() string {
    return fmt.Sprintf("%s|%s|%v", c.SecretName, c.Key, c.Optional)
}

type EnvValue struct {
    Value              string                 `json:"value" yaml:"value"`
    ValueFromConfigMap *EnvValueFromConfigMap `json:"valueFromConfigMap,omitempty" yaml:"valueFromConfigMap,omitempty"`
    ValueFromSecret    *EnvValueFromSecret    `json:"envValueFromSecret,omitempty" yaml:"envValueFromSecret,omitempty"`
}

func (c *EnvValue) String() string {

    if c.ValueFromConfigMap != nil {
        return c.ValueFromConfigMap.String()
    }

    if c.ValueFromSecret != nil {
        return c.ValueFromSecret.String()
    }
    return c.Value
}

var (
    // ValueFromReg prefix: ####   suffix: ####
    ValueFromReg = `(?P<valueFromPrefix>^####)(?P<valueFrom>(.*))(?P<valueFromSuffix>####$)`
)

func IsValueFrom(value string) (bool, string) {
    result := make(map[string]string)
    reg := regexp.MustCompile(ValueFromReg)
    if !reg.MatchString(value) {
        return false, ""
    }

    valueFrom := reg.FindStringSubmatch(value)

    groupName := reg.SubexpNames()
    for index, name := range groupName {
        if index != 0 && name != "" {
            result[name] = valueFrom[index]
        }
    }
    return true, result["valueFrom"]
}

func ParseEnvValue(value string) (*EnvValue, error) {
    envValue := &EnvValue{}
    isValueFrom, v := IsValueFrom(value)
    if !isValueFrom {
        envValue.Value = value
        return envValue, nil
    }
    valueStr := strings.Split(v, ".")
    switch len(valueStr) {
    case 2:
        envValue.ValueFromConfigMap = &EnvValueFromConfigMap{}
        envValue.ValueFromConfigMap.ConfigMapName = valueStr[0]
        envValue.ValueFromConfigMap.Key = valueStr[1]
    case 3:
        envValue.ValueFromSecret = &EnvValueFromSecret{}
        envValue.ValueFromSecret.SecretName = valueStr[0]
        envValue.ValueFromSecret.Key = valueStr[1]
        if strings.ToLower(valueStr[2]) == "true" {
            envValue.ValueFromSecret.Optional = true
        } else if strings.ToLower(valueStr[2]) == "false" {
            envValue.ValueFromSecret.Optional = false
        } else {
            return nil, fmt.Errorf("secret optional str error, should be  true or false")
        }
    }
    return envValue, nil
}
