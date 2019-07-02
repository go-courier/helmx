package spec

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestRequestAndLimit(t *testing.T) {
	t.Run("parse & string", func(t *testing.T) {
		r, _ := ParseRequestAndLimit("10/500")

		require.Equal(t, 10, r.Request)
		require.Equal(t, 500, r.Limit)

		require.Equal(t, "10/500", r.String())
	})

	t.Run("parse & string simple", func(t *testing.T) {
		r, _ := ParseRequestAndLimit("10")

		require.Equal(t, 10, r.Request)
		require.Equal(t, 0, r.Limit)

		require.Equal(t, "10", r.String())
	})
}
