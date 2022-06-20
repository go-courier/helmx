package spec

import (
    "github.com/stretchr/testify/require"
    "testing"
)

func TestIsValueFrom(t *testing.T) {
    t.Run("not value from", func(t *testing.T) {
        is, v := IsValueFrom("a.b.c")
        require.Equal(t, false, is)
        require.Equal(t, "", v)
    })

    t.Run("not value from with prefix #", func(t *testing.T) {
        is, v := IsValueFrom("###a.b.c####")
        require.Equal(t, false, is)
        require.Equal(t, "", v)
    })

    t.Run("not value from without prefix #", func(t *testing.T) {
        is, v := IsValueFrom("a.b.c####")
        require.Equal(t, false, is)
        require.Equal(t, "", v)
    })

    t.Run("not value from with suffix #", func(t *testing.T) {
        is, v := IsValueFrom("###a.b.c")
        require.Equal(t, false, is)
        require.Equal(t, "", v)
    })

    t.Run("not value from with suffix #", func(t *testing.T) {
        is, v := IsValueFrom("###a.b.c###")
        require.Equal(t, false, is)
        require.Equal(t, "", v)
    })

    t.Run("value from", func(t *testing.T) {
        is, v := IsValueFrom("####a.b.c####")
        require.Equal(t, true, is)
        require.Equal(t, "a.b.c", v)
    })

    t.Run("value from with #  prefix", func(t *testing.T) {
        is, v := IsValueFrom("######a.b.c####")
        require.Equal(t, true, is)
        require.Equal(t, "##a.b.c", v)
    })

    t.Run("value from with # suffix  ", func(t *testing.T) {
        is, v := IsValueFrom("####a.b.c#######")
        require.Equal(t, true, is)
        require.Equal(t, "a.b.c###", v)
    })

    t.Run("value from with # prefix and suffix  ", func(t *testing.T) {
        is, v := IsValueFrom("######a.b.c#######")
        require.Equal(t, true, is)
        require.Equal(t, "##a.b.c###", v)
    })
}

func TestParseEnvValue(t *testing.T) {
    t.Run("default", func(t *testing.T) {
        envValue, _ := ParseEnvValue("default_value")
        require.Equal(t, "default_value", envValue.Value)
    })

    t.Run("configmap", func(t *testing.T) {
        envValue, _ := ParseEnvValue("configMapName.configMapKey")
        require.Equal(t, "configMapName", envValue.ValueFromConfigMap.ConfigMapName)
        require.Equal(t, "configMapKey", envValue.ValueFromConfigMap.Key)
    })

    t.Run("secret_true", func(t *testing.T) {
        envValue, _ := ParseEnvValue("secretName.secretKey.true")
        require.Equal(t, "secretName", envValue.ValueFromSecret.SecretName)
        require.Equal(t, "secretKey", envValue.ValueFromSecret.Key)
        require.Equal(t, true, envValue.ValueFromSecret.Optional)
    })

    t.Run("secret_false", func(t *testing.T) {
        envValue, _ := ParseEnvValue("secretName.secretKey.false")
        require.Equal(t, "secretName", envValue.ValueFromSecret.SecretName)
        require.Equal(t, "secretKey", envValue.ValueFromSecret.Key)
        require.Equal(t, false, envValue.ValueFromSecret.Optional)
    })

}

func TestParseEnvs(t *testing.T) {
    t.Run("default", func(t *testing.T) {
        envs, _ := ParseEnvsWithValueFrom(map[string]string{"key": "value"})
        expectEnvs := make(EnvsWithValueFrom)
        expectEnvs["key"] = &EnvValue{
            Value: "value",
        }
        require.Equal(t, expectEnvs, envs)
    })

    t.Run("configmap", func(t *testing.T) {
        envs, _ := ParseEnvsWithValueFrom(map[string]string{"key": "configmapName.configMapKey"})
        expectEnvs := make(EnvsWithValueFrom)
        expectEnvs["key"] = &EnvValue{
            ValueFromConfigMap: &EnvValueFromConfigMap{
                ConfigMapName: "configmapName",
                Key:           "configMapKey",
            },
        }
        require.Equal(t, expectEnvs, envs)
    })

    t.Run("secret", func(t *testing.T) {
        envs, _ := ParseEnvsWithValueFrom(map[string]string{"key": "secretName.secretKey.true"})
        expectEnvs := make(EnvsWithValueFrom)
        expectEnvs["key"] = &EnvValue{
            ValueFromSecret: &EnvValueFromSecret{
                SecretName: "secretName",
                Key:        "secretKey",
                Optional:   true,
            },
        }
        require.Equal(t, expectEnvs, envs)
    })

    t.Run("mixed", func(t *testing.T) {
        envs, _ := ParseEnvsWithValueFrom(map[string]string{
            "key1": "secretName.secretKey.true",
            "key2": "configmapName.configMapKey",
            "key3": "value",
        })
        expectEnvs := make(EnvsWithValueFrom)
        expectEnvs["key1"] = &EnvValue{
            ValueFromSecret: &EnvValueFromSecret{
                SecretName: "secretName",
                Key:        "secretKey",
                Optional:   true,
            },
        }
        expectEnvs["key2"] = &EnvValue{
            ValueFromConfigMap: &EnvValueFromConfigMap{
                ConfigMapName: "configmapName",
                Key:           "configMapKey",
            },
        }
        expectEnvs["key3"] = &EnvValue{
            Value: "value",
        }
        require.Equal(t, expectEnvs, envs)
    })

}
