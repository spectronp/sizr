package main

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func mockRunner( script string, args ...string ) (output string, err error ) { // TODO -- refactor
	shell := "/usr/bin/bash"
	path := "./tests/scripts-mock.sh"
	pathAndArgs := append([]string{path, script}, args...)
	cmd := exec.Command(shell, pathAndArgs...)
	stdout, err := cmd.Output()
	if err != nil {
		fmt.Printf("Error on mockRunner: %s", err)
	}
	return string(stdout), err
}

func TestCanBuildData(t *testing.T )  {
	_, err := NewData(mockRunner) // NOTE -- do this inside InitMain ?

	if err != nil {
		t.Error("Data could not be built: ", err)
	}
}

func TestManagerIsSet(t *testing.T) {
	data, _ := NewData(mockRunner)
	
	if data.Manager != "pm" {
		t.Errorf("data.Manager should be 'pm', is '%s' instead", data.Manager)
	}
}

func TestGetPackage(t *testing.T) {
	data, _ := NewData(mockRunner)
	
	expectedPackage := Package{Name: "exp1", Size: 10000, IsExplicit: true, Deps: []string{"rev1", "dep1", "rev4"}}

	if ! cmp.Equal(data.GetPackage("exp1"), expectedPackage) {
		diff := cmp.Diff(data.GetPackage("exp1"), expectedPackage)
		t.Errorf("Package exp1 is not what was expected\n%s", diff)
	}
}

func TestGetExplicit(t *testing.T) {
	data, _ := NewData(mockRunner)
	lenBefore := len(data.PackageList)	

	explicitPackages := data.GetExplicit()
	lenAfter := len(data.PackageList)

	if len(explicitPackages) != 3 {
		t.Errorf("Expected 3 packages, got %d", len(explicitPackages))
	}

	if lenBefore != lenAfter {
		t.Errorf("Data have been modified, len before: %d, len after: %d", lenBefore, lenAfter)
	}
	// TODO -- compare both
}
// TODO -- test []Package
