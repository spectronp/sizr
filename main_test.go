package main

import (
	"fmt"
	"io"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var data Data;

func TestMain(m *testing.M) {
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

func runApp() (int, string) {
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
		returnCode = Run()
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

// TestVersionOutput

// TestHelpOutput

func TestListReport( t *testing.T ) {
	returnCode, output := runApp()

	if returnCode != 0 {
		t.Errorf("Expected return code 0, got %d", returnCode)
	}
	
	expectedOutput := "exp1\t51200\nexp3\t40960\nexp2\t30720\n"
	if output != expectedOutput {
		fmt.Print(cmp.Diff(expectedOutput, output))
		t.Error("Output is different from the expected")
	}
}

// TestLimitReport

// 
