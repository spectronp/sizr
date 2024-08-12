package helpers

import (
	"github.com/spectronp/sizr/vars"

	"fmt"
	"os/exec"
)

func MockRunner(script string, args ...string) (output string, err error) {
	scriptsMock := vars.BASEDIR + "/tests/scripts-mock.sh"
	scriptAndArgs := append([]string{script}, args...)
	cmd := exec.Command(scriptsMock, scriptAndArgs...)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error on MockRunner: %s", err)
	}
	return string(stdout), err
}
