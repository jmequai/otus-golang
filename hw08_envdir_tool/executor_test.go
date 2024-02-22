package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	t.Run("Successful without args", func(t *testing.T) {
		cmd := []string{"true"}
		exitCode := RunCmd(cmd, Environment{})

		require.Equal(t, 0, exitCode)
	})

	t.Run("Successful with args", func(t *testing.T) {
		cmd := []string{"whoami", "--help"}
		exitCode := RunCmd(cmd, Environment{})

		require.Equal(t, 0, exitCode)
	})

	t.Run("Successful with env", func(t *testing.T) {
		cmd := []string{"whoami", "--help"}
		exitCode := RunCmd(cmd, Environment{
			"OTUS_TEST": EnvValue{
				Value: "123",
			},
		})

		require.Equal(t, 0, exitCode)
		require.Contains(t, os.Environ(), "OTUS_TEST=123")
	})

	t.Run("Failure without args", func(t *testing.T) {
		cmd := []string{"false"}
		exitCode := RunCmd(cmd, Environment{})

		require.Equal(t, 1, exitCode)
	})

	t.Run("No command", func(t *testing.T) {
		cmd := []string{}
		exitCode := RunCmd(cmd, Environment{})

		require.Equal(t, DefaultErrCode, exitCode)
	})
}
