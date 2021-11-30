package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	dir := "./testdata/env"

	t.Run("test with correct direction", func(t *testing.T) {
		os.Setenv("BAR", "123")
		os.Setenv("HELLO", `"goodbye"`)
		os.Setenv("EMPTY", "smth")
		os.Setenv("TAB", "value")
		os.Setenv("FOO", "ololo")
		os.Setenv("INVALID", "invalid")
		os.Setenv("UNSET", "set")
		defer func() {
			os.Unsetenv("BAR")
			os.Unsetenv("HELLO")
			os.Unsetenv("EMPTY")
			os.Unsetenv("TAB")
			os.Unsetenv("FOO")
			os.Unsetenv("INVALID")
			os.Unsetenv("UNSET")
		}()

		expectedEnv := Environment{
			"BAR":   {"bar", false},
			"HELLO": {`"hello"`, false},
			"EMPTY": {"", false},
			"TAB":   {"123", false},
			"FOO":   {"   foo\nwith new line", false},
			"UNSET": {"", true},
		}

		env, err := ReadDir(dir)
		require.NoError(t, err)
		require.Equal(t, env, expectedEnv)
	})

	t.Run("test with invalid direction", func(t *testing.T) {
		_, err := ReadDir("")
		require.Error(t, err)
	})
}
