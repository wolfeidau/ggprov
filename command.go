package ggprov

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"syscall"

	"github.com/pkg/errors"
)

type ExitError struct {
	*os.ProcessState
	ExitCode int
	Stderr   []byte
}

func (ee *ExitError) Error() string {
	return ee.ProcessState.String()
}

// UserExist check if user exists on the system
func UserExist(username string) (bool, error) {
	err := RunCommand("id", []string{"-u", username})
	if err != nil {
		if exitError, ok := err.(*ExitError); ok {
			if exitError.ExitCode == 1 {
				return false, nil
			}
		}
		return false, err
	}

	return true, nil
}

// NewSystemUserAndGroup create a new system user and group
func NewSystemUserAndGroup(username string) error {
	cmdName := "useradd"
	cmdArgs := []string{"--system", "--no-create-home", "--user-group", username}

	return RunCommand(cmdName, cmdArgs)
}

// RunCommand run a command with the supplied arguments
func RunCommand(cmdName string, cmdArgs []string) error {
	cmd := exec.Command(cmdName, cmdArgs...)

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return errors.Wrap(err, "Failed to build reader")
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Printf("%s out | %s\n", cmdName, scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		return errors.Wrap(err, "Failed to start command")
	}

	err = cmd.Wait()
	if err != nil {
		if exitError, ok := err.(*exec.ExitError); ok {
			if waitStatus, ok := exitError.Sys().(syscall.WaitStatus); ok {
				return &ExitError{ProcessState: exitError.ProcessState, ExitCode: waitStatus.ExitStatus()}
			}
		}
		return errors.Wrap(err, "Failed while waiting for command to complete")
	}

	return nil
}
