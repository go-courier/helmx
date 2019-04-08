package spec

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAction(t *testing.T) {
	t.Run("parse & string http", func(t *testing.T) {
		action, _ := ParseAction("http://:80/healthy")

		require.Equal(t, uint16(80), action.HTTPGet.Port)
		require.Equal(t, "HTTP", action.HTTPGet.Scheme)
		require.Equal(t, "", action.HTTPGet.Host)
		require.Equal(t, "/healthy", action.HTTPGet.Path)

		require.Equal(t, "http://:80/healthy", action.String())
	})

	t.Run("parse & string tcp", func(t *testing.T) {
		action, _ := ParseAction("tcp://:80")

		require.Equal(t, uint16(80), action.TCPSocket.Port)
		require.Equal(t, "", action.TCPSocket.Host)

		require.Equal(t, "tcp://:80", action.String())
	})
}
