package main

import (
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	t.Run("env dir", func(t *testing.T) {
		expected := Environment{
			"BAR":   EnvValue{Value: "bar", NeedRemove: false},
			"EMPTY": EnvValue{Value: "", NeedRemove: true},
			"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
			"HELLO": EnvValue{Value: "\"hello\"", NeedRemove: false},
			"UNSET": EnvValue{Value: "", NeedRemove: true},
		}

		env, err := ReadDir("testdata/env")

		require.NoError(t, err)
		require.Equal(t, expected, env)
	})

	t.Run("empty env dir", func(t *testing.T) {
		d, e := os.MkdirTemp("", "empty-dir")
		if e != nil {
			t.Fatal("can't create temp dir: ", e)
		}

		defer os.RemoveAll(d)

		env, err := ReadDir(d)

		require.NoError(t, err)
		require.Len(t, env, 0)
	})

	t.Run("env dir not exists", func(t *testing.T) {
		envPath := path.Join("testdata", "not-exists")
		env, err := ReadDir(envPath)

		require.Error(t, err)
		require.Equal(t, Environment(nil), env)
	})
}
