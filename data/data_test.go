package data

import (
	"fmt"
	"os/exec"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func mockRunner( script string, args ...string ) (output string, err error ) {
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
	
	expectedPackage := Package{Name: "exp1", Size: 10240, IsExplicit: true, Version: "0.0.0", Deps: []string{"rev1", "dep1", "rev4"}}

	if ! cmp.Equal(data.GetPackage("exp1"), expectedPackage) {
		diff := cmp.Diff(data.GetPackage("exp1"), expectedPackage)
		t.Errorf("Package exp1 is not what was expected\n%s", diff)
	}
}

func TestGetExplicit(t *testing.T) {
	expectedPackages := map[string]Package{
		"exp1": {
			Name: "exp1",
			Size: 10240,
			IsExplicit: true,
			Version: "0.0.0",
			Deps: []string{"rev1", "dep1", "rev4"},
		},
		"exp2": {
			Name: "exp2",
			Size: 10240,
			IsExplicit: true,
			Version: "0.0.0",
			Deps: []string{"dep13", "dep14", "rev4"},
		},
		"exp3": {
			Name: "exp3",
			Size: 10240,
			IsExplicit: true,
			Version: "0.0.0",
			Deps: []string{"dep11", "dep12"},
		},
	}

	data, _ := NewData(mockRunner)

	actualPackages := data.GetExplicit()

	if ! cmp.Equal(expectedPackages, actualPackages) {
		fmt.Print(cmp.Diff(expectedPackages, actualPackages))
		t.Errorf("Returned packages are different than expected")
	}
}

