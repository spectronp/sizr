package helpers

import (
	"github.com/spectronp/sizr/vars"

	"fmt"
	"os/exec"
)

func MockRunner(script string, args ...string) (output string, err error) {
	shell := "/usr/bin/bash" // TODO -- run scripts-mock directly
	path := vars.BASEDIR + "/tests/scripts-mock.sh"
	pathAndArgs := append([]string{path, script}, args...)
	cmd := exec.Command(shell, pathAndArgs...)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error on MockRunner: %s", err)
	}
	return string(stdout), err
}
