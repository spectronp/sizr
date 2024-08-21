package main

import (
	"github.com/spectronp/sizr/data"
	"github.com/spectronp/sizr/tests"
	"github.com/spectronp/sizr/types"
	"github.com/spectronp/sizr/utils"
	"github.com/spectronp/sizr/vars"

	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var mockData data.Data

func TestMain(m *testing.M) {
	VERSION = "v0.1"
	vars.BASEDIR = os.Getenv("BASEDIR")
	helpBytes, err := os.ReadFile("help.txt")
	if err != nil {
		log.Println("Error on reading help.txt")
		panic(err)
	}
	HELP_MESSAGE = string(helpBytes)

	vars.DB_FILE = "/tmp/sizr_db"
	defer os.Remove("/tmp/sizr_db")
	if err := utils.SaveJson(map[string]any{}, vars.DB_FILE); err != nil {
		panic(err)
	}

	mockData, _ = data.NewData(helpers.MockRunner) // NOTE: if Data is broken, this data will break this tests too, maybe use mockgen ? ( data is used more along the file )
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
	depsToIgnore := map[string]bool{
		"dep13": true,
	}

	expectedPackages := map[string]bool{"dep13": true, "dep11": true, "dep12": true, "dep15": true, "rev2": true, "rev3": true, "dep5": true, "dep6": true, "dep9": true, "dep10": true}

	listTree(firstPack, depsToIgnore, &mockData)
	if !cmp.Equal(expectedPackages, depsToIgnore) {
		fmt.Println(cmp.Diff(expectedPackages, depsToIgnore))
		t.Error("listTree function returned something different from expected")
	}
}

func TestSumSize(t *testing.T) {
	start := mockData.GetPackage("exp1")
	ignorePackagesNames := []string{"rev1", "rev2", "rev3", "rev4"}
	ignoredPackages := make(map[string]bool)
	for _, ignoredName := range ignorePackagesNames {
		ignoredPackages[ignoredName] = true
	}
	alreadyCounted := map[string]bool{}

	expectedSize := uint(51200)
	actualSize := sumSize(start, ignoredPackages, alreadyCounted, &mockData)

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

	if returnCode != 0 {
		t.Errorf("Expected return code 0, got %d", returnCode)
	}

	expectedOutput := "exp1     50 KiB\nexp3     40 KiB\nexp2     20 KiB\n"
	if output != expectedOutput {
		fmt.Print(cmp.Diff(expectedOutput, output))
		t.Error("Output is different from the expected")
	}
}

func TestLimitReport(t *testing.T) {
	args := []string{"--limit", "2"}
	expectedOutput := "exp1     50 KiB\nexp3     40 KiB\n"

	returnCode, output := runApp(args)

	if returnCode != 0 {
		t.Errorf("Expected return code 0, got %d", returnCode)
	}

	if output != expectedOutput {
		fmt.Print(cmp.Diff(expectedOutput, output))
		t.Error("Output is different from expected")
	}
}
