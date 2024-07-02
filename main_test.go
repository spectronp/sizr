package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var data Data;

func TestMain(m *testing.M) {
	VERSION = "v0.1"
	helpBytes, err := os.ReadFile(".help")
	if err != nil {
		log.Println("Error on reading .help")	
		panic(err)
	}
	HELP_MESSAGE = string(helpBytes)

	data, _ = NewData(mockRunner) // NOTE -- if Data is broken, this data will break this tests too ( data is used more along the file )
	m.Run()	
}

func TestListTree(t *testing.T)  {
	firstPack := Package{
		Name: "exp3",
		IsExplicit: true,
		Size: 10240,
		Deps: []string{"dep11", "dep12"},
	}

	expectedList := []string{"dep11", "dep12", "dep15", "rev2", "rev3", "dep5", "dep6", "dep9", "dep10"}
	expectedPackages := map[string]Package{}

	for _, name := range expectedList {
		expectedPackages[name] = data.GetPackage(name)
	}

	receivedPackages := listTree(firstPack, data)
	if ! cmp.Equal(expectedPackages, receivedPackages) {
		fmt.Println(cmp.Diff(expectedPackages, receivedPackages))
		t.Error("listTree function returned something different from expected")
	}
}



func TestSumSize(t *testing.T) {
	start := data.GetPackage("exp1")
	ignorePackagesNames := []string{"rev1", "rev2", "rev3", "rev4"}
	ignoredPacakges := make(map[string]Package)
	for _, ignoredName := range ignorePackagesNames {
		ignoredPacakges[ignoredName] = data.GetPackage(ignoredName)	
	}

	expectedSize := uint(50000)
	actualSize := sumSize(start, ignoredPacakges, data)

	if expectedSize != actualSize {
		t.Errorf("Expected size %d, got size %d", expectedSize, actualSize)	
	}
}

func TestCalcSize(t *testing.T) {
	expectedSize := 50000
	actualSize := calcSize("exp1", data)
	
	if actualSize != uint(expectedSize) {
		t.Errorf("Expected size %d, got size %d", expectedSize, actualSize)
	}
}

func TestOrderBySum(t *testing.T) {
	expectedPackages := []PackageNameWithSum{
		{Name: "exp1", Size: 50000},
		{Name: "exp3", Size: 40000},
		{Name: "exp2", Size: 30000},
	}
	actualPackages := orderBySumSize(data)	

	if ! cmp.Equal(expectedPackages, actualPackages) {
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

	expectedOutput := "sizr v0.1\n" // TODO -- get this automatticaly

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

func TestListReport( t *testing.T ) {
	returnCode, output := runApp([]string{})

	if returnCode != 0 {
		t.Errorf("Expected return code 0, got %d", returnCode)
	}
	
	expectedOutput := "exp1\t51200\nexp3\t40960\nexp2\t30720\n"
	if output != expectedOutput {
		fmt.Print(cmp.Diff(expectedOutput, output))
		t.Error("Output is different from the expected")
	}
}

func TestLimitReport(t *testing.T) {
	args := []string{"--limit", "2"}
	expectedOutput := "exp1\t51200\nexp3\t40960\n"	

	returnCode, output := runApp(args)

	if returnCode != 0 {
		t.Errorf("Expected return code 0, got %d", returnCode)
	}

	if output != expectedOutput {
		fmt.Print(cmp.Diff(expectedOutput, output))
		t.Error("Output is different from expected")
	}
}

