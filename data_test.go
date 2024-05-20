package main

import (
	"testing"
	"os/exec"
)

func mockRunner( script string, args ...string ) (output string, err error ) { // TODO -- refactor
	shell := "/usr/bin/bash"
	path := "./tests/scripts-mock.sh"
	pathAndArgs := append([]string{path}, args...)
	cmd := exec.Command(shell, pathAndArgs...)
	stdout, err := cmd.Output()
	return string(stdout), err
}

func TestCanBuildData(t *testing.T )  {
	_, err := NewData(mockRunner)

	if err != nil {
		t.Error()
	}
}

// TODO -- test []Package
