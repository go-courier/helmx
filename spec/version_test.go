package spec

import (
	"testing"

	"github.com/stretchr/testify/require"
	"gopkg.in/yaml.v2"
)

func TestVersion(t *testing.T) {
	t.Run("simple-version", func(t *testing.T) {
		v, _ := ParseVersion("1.2.3")

		require.Equal(t, 1, v.Major)
		require.Equal(t, 2, v.Minor)
		require.Equal(t, 3, v.Patch)

		require.Equal(t, v.String(), "1.2.3")
	})

	t.Run("version with suffix", func(t *testing.T) {
		v, _ := ParseVersion("1.2.3-xxxxxx")

		require.Equal(t, 1, v.Major)
		require.Equal(t, 2, v.Minor)
		require.Equal(t, 3, v.Patch)
		require.Equal(t, "xxxxxx", v.Suffix)

		require.Equal(t, "1.2.3-xxxxxx", v.String())
	})

	t.Run("version with linePrefix", func(t *testing.T) {
		v, _ := ParseVersion("feat-1.2.3")

		require.Equal(t, "feat", v.Prefix)

		require.Equal(t, 1, v.Major)
		require.Equal(t, 2, v.Minor)
		require.Equal(t, 3, v.Patch)
		require.Equal(t, "feat-1.2.3", v.String())
	})

	t.Run("version with linePrefix and suffix", func(t *testing.T) {
		v, _ := ParseVersion("feat-1.2.3-xxxxxx")

		require.Equal(t, "feat", v.Prefix)
		require.Equal(t, 1, v.Major)
		require.Equal(t, 2, v.Minor)
		require.Equal(t, 3, v.Patch)
		require.Equal(t, "xxxxxx", v.Suffix)

		require.Equal(t, "feat-1.2.3-xxxxxx", v.String())
	})
}

func TestVersion_Yaml(t *testing.T) {
	t.Run("marshal & unmarshal", func(t *testing.T) {
		data, err := yaml.Marshal(struct {
			Version Version `yaml:"version"`
		}{
			Version: Version{
				Prefix: "feat",
				Suffix: "xxx",
				Major:  1,
				Minor:  2,
				Patch:  3,
			},
		})
		require.NoError(t, err)
		require.Equal(t, "version: feat-1.2.3-xxx\n", string(data))

		v := struct {
			Version Version `yaml:"version"`
		}{
			Version: Version{},
		}

		err = yaml.Unmarshal(data, &v)
		require.NoError(t, err)
		require.Equal(t, "feat-1.2.3-xxx", v.Version.String())
	})
}
