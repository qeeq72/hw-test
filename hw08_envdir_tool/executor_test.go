package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("invalid command", func(t *testing.T) {
		env := Environment{
			"BAR":   {"bar", false},
			"HELLO": {`"hello"`, false},
			"EMPTY": {"", false},
			"TAB":   {"123", false},
			"FOO":   {"   foo\nwith new line", false},
			"UNSET": {"", true},
		}
		result := RunCmd(nil, env)
		require.Equal(t, 1, result)
	})

	t.Run("command without arguments", func(t *testing.T) {
		env := Environment{
			"BAR":   {"bar", false},
			"HELLO": {`"hello"`, false},
			"EMPTY": {"", false},
			"TAB":   {"123", false},
			"FOO":   {"   foo\nwith new line", false},
			"UNSET": {"", true},
		}
		result := RunCmd([]string{"pwd"}, env)
		require.Equal(t, 0, result)
	})

	t.Run("command with arguments", func(t *testing.T) {
		env := Environment{
			"BAR":   {"bar", false},
			"HELLO": {`"hello"`, false},
			"EMPTY": {"", false},
			"TAB":   {"123", false},
			"FOO":   {"   foo\nwith new line", false},
			"UNSET": {"", true},
			"PORT":  {"8080", false},
		}
		result := RunCmd([]string{"./testdata/exit.sh", "1"}, env)
		require.Equal(t, 123, result)
	})
}
