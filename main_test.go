package main

import (
	"github.com/spectronp/sizr/data"
	"github.com/spectronp/sizr/tests"
	"github.com/spectronp/sizr/types"
	"github.com/spectronp/sizr/vars"

	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var mockData data.Data

func TestMain(m *testing.M) {
	VERSION = "v0.1"
	helpBytes, err := os.ReadFile(".help")
	if err != nil {
		log.Println("Error on reading .help")
		panic(err)
	}
	HELP_MESSAGE = string(helpBytes)

	vars.DB_FILE = "/tmp/sizr_db"

	mockData, _ = data.NewData(helpers.MockRunner) // NOTE -- if Data is broken, this data will break this tests too ( data is used more along the file )
	m.Run()
}

func TestListTree(t *testing.T) {
	firstPack := types.Package{
		Name:       "exp3",
		IsExplicit: true,
		Size:       10240,
		Version:    "0.0.0",
		Deps:       []string{"dep11", "dep12"},
	}
	depsToIgnore := map[string]types.Package{
		"dep13": {
			Name:       "dep13",
			Size:       10240,
			IsExplicit: false,
			Version:    "0.0.0",
			Deps:       []string{"dep16"},
		},
	}

	expectedList := []string{"dep13", "dep11", "dep12", "dep15", "rev2", "rev3", "dep5", "dep6", "dep9", "dep10"}
	expectedPackages := map[string]types.Package{}

	for _, name := range expectedList {
		expectedPackages[name] = mockData.GetPackage(name)
	}

	listTree(firstPack, depsToIgnore, &mockData) // TODO -- assert it ignores the packages from second arg
	if !cmp.Equal(expectedPackages, depsToIgnore) {
		fmt.Println(cmp.Diff(expectedPackages, depsToIgnore))
		t.Error("listTree function returned something different from expected")
	}
}

func TestSumSize(t *testing.T) {
	start := mockData.GetPackage("exp1")
	ignorePackagesNames := []string{"rev1", "rev2", "rev3", "rev4"}
	ignoredPacakges := make(map[string]types.Package)
	for _, ignoredName := range ignorePackagesNames {
		ignoredPacakges[ignoredName] = mockData.GetPackage(ignoredName)
	}

	expectedSize := uint(51200)
	actualSize := sumSize(start, ignoredPacakges, &mockData)

	if expectedSize != actualSize {
		t.Errorf("Expected size %d, got size %d", expectedSize, actualSize)
	}
}

func TestCalcSize(t *testing.T) {
	expectedSize := 51200
	actualSize := calcSize("exp1", &mockData)

	if actualSize != uint(expectedSize) {
		t.Errorf("Expected size %d, got size %d", expectedSize, actualSize)
	}
}

func TestOrderBySum(t *testing.T) {
	expectedPackages := []PackageNameWithSum{
		{Name: "exp1", Size: 51200},
		{Name: "exp3", Size: 40960},
		{Name: "exp2", Size: 20480},
	}
	actualPackages := orderBySumSize(&mockData)

	if !cmp.Equal(expectedPackages, actualPackages) {
		fmt.Println(cmp.Diff(expectedPackages, actualPackages))
		t.Errorf("The received packages are different from expected")
	}
}

// E2E Tests

func removeProgressBar(output string) string {
	_, outputWithoutBar, _ := strings.Cut(output, "@END_PROGRESSBAR@\n")
	return outputWithoutBar
}

func runApp(args []string) (int, string) {
	args = append([]string{"sizr"}, args...)

	pipeReader, pipeWriter, err := os.Pipe()
	if err != nil {
		panic(err)
	}

	stdOut := os.Stdout
	stdErr := os.Stderr
	os.Stdout = pipeWriter
	os.Stderr = pipeWriter
	var returnCode int
	go func() {
		returnCode = Run(args)
		pipeWriter.Close()
	}()

	output, err := io.ReadAll(pipeReader)
	if err != nil {
		panic("Panic at io.ReadAll()")
	}
	os.Stdout = stdOut
	os.Stderr = stdErr

	return returnCode, string(output)
}

func TestVersionOutput(t *testing.T) {
	args := []string{"--version"}

	expectedOutput := fmt.Sprintf("sizr %s\n", VERSION)
	returnCode, output := runApp(args)

	if returnCode != 0 {
		t.Errorf("Expected return code 0, got %d", returnCode)
	}

	if output != expectedOutput {
		t.Errorf("expected output: %s, got output: %s", expectedOutput, output)
	}
}

func TestHelpOutput(t *testing.T) {
	args := []string{"--help"}

	expectedOutput := HELP_MESSAGE
	returnCode, output := runApp(args)

	if returnCode != 0 {
		t.Errorf("Expected return code 0, got %d", returnCode)
	}

	if output != expectedOutput {
		fmt.Print(cmp.Diff(expectedOutput, output))
		t.Error("Output is different from expected")
	}

}

func TestListReport(t *testing.T) {
	returnCode, output := runApp([]string{})
	output = removeProgressBar(output)

	if returnCode != 0 {
		t.Errorf("Expected return code 0, got %d", returnCode)
	}

	expectedOutput := "exp1     51200\nexp3     40960\nexp2     20480\n"
	if output != expectedOutput {
		fmt.Print(cmp.Diff(expectedOutput, output))
		t.Error("Output is different from the expected")
	}
}

func TestLimitReport(t *testing.T) {
	args := []string{"--limit", "2"}
	expectedOutput := "exp1     51200\nexp3     40960\n"

	returnCode, output := runApp(args)
	output = removeProgressBar(output)

	if returnCode != 0 {
		t.Errorf("Expected return code 0, got %d", returnCode)
	}

	if output != expectedOutput {
		fmt.Print(cmp.Diff(expectedOutput, output))
		t.Error("Output is different from expected")
	}
}
