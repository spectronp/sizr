package main

import (
	"encoding/json"
	"os"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var data Data;

func Init() {
	jsonData, _ := os.ReadFile("packages.json")  
	json.Unmarshal(jsonData, &data)
}

func TestListTree(t *testing.T)  {
	firstPack := Package{
		name: "zjvut",
		isExplict: true,
		size: 2704100,
		deps: []string{"kcptm", "tzxvn", "lebsv"},
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

// TestSumSize

// TestCalcSize
