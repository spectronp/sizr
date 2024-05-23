package main

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

var data Data;

func TestMain(m *testing.M) {
	data, _ = NewData(mockRunner)
	m.Run()	
}

func TestListTree(t *testing.T)  {
	firstPack := Package{
		Name: "zjvut",
		IsExplicit: true,
		Size: 2704100,
		Deps: []string{"kcptm", "tzxvn", "lebsv"},
	}

	expectedList := []string{"kcptm", "yavjp", "ujgpd", "tqkwg", "ioleb", "fmesb", "tzxvn", "klkrz", "lebsv"}
	expectedPackages := map[string]Package{}

	for _, name := range expectedList {
		expectedPackages[name] = data.GetPackage(name)
	}

	if cmp.Equal(expectedList, listTree(firstPack, data)) {
		t.Fail()
	}
}

func TestOrderBySum(t *testing.T) {
	orderedPackages := orderBySumSize(data)	

	if len(orderedPackages) != 3 {
		t.Errorf("Expected 3 packages, got %d", len(orderedPackages))
	}
}

// TestSumSize

// TestCalcSize
