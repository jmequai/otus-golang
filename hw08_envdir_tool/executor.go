package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

const DefaultErrCode = 128

var ErrCmdRequired = errors.New("command required")

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		fmt.Println(ErrCmdRequired)
		return DefaultErrCode
	}

	if err := initEnvironment(env); err != nil {
		fmt.Println(err)
		return DefaultErrCode
	}

	command := exec.Command(cmd[0], cmd[1:]...) //nolint:gosec

	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		var errExit *exec.ExitError

		if errors.As(err, &errExit) {
			return errExit.ExitCode()
		}

		fmt.Println(err)
		return DefaultErrCode
	}

	return 0
}

func initEnvironment(env Environment) error {
	var err error

	for name, val := range env {
		if val.NeedRemove {
			err = os.Unsetenv(name)
		} else {
			err = os.Setenv(name, val.Value)
		}

		if err != nil {
			return err
		}
	}

	return nil
}
