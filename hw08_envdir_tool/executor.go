package main

import (
	"errors"
	"os"
	"os/exec"
)

const (
	SuccessCode = iota
	FailCode
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return FailCode
	}

	for key, val := range env {
		if val.NeedRemove {
			err := os.Unsetenv(key)
			if err != nil {
				return FailCode
			}

			continue
		}

		err := os.Setenv(key, val.Value)
		if err != nil {
			return FailCode
		}
	}

	command := exec.Command(cmd[0], cmd[1:]...) // #nosec G204
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Stdin = os.Stdin

	if err := command.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}

		return FailCode
	}

	return SuccessCode
}
