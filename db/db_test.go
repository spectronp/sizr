package db

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/go-cmp/cmp"
	"github.com/spectronp/sizr/types"
	"github.com/spectronp/sizr/utils"
	"github.com/spectronp/sizr/vars"

	"testing"
)

var tmpDir string // TODO -- use MapFS
var packagesJson string
var tmpPackagesJson string

func TestMain(m *testing.M) {
	var err error
	tmpDir, err = os.MkdirTemp("", "sizr_tests-")
	defer os.RemoveAll(tmpDir)
	if err != nil {
		panic(err)
	}

	packagesJson = tmpDir + "/original.json"
	tmpPackagesJson = tmpDir + "/modified.json"

	var packages map[string]types.Package
	packages, err = utils.LoadJson[map[string]types.Package](vars.BASEDIR + "/tests/packages.json")
	if err != nil {
		panic(err)
	}

	keyMapPackages := map[string]types.Package{}
	for packName, pack := range packages {
		if !pack.IsExplicit {
			continue
		}
		keyMapPackages[fmt.Sprintf("%s %s", packName, pack.Version)] = pack
	}
	utils.SaveJson(keyMapPackages, packagesJson)

	vars.DB_FILE = packagesJson

	m.Run()
}

func loadPackagesJson(t *testing.T, path string) map[string]types.Package {
	t.Helper()

	packagesFile, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	var packagesMap map[string]types.Package
	json.Unmarshal(packagesFile, &packagesMap)

	return packagesMap
}

func savePackagesJson(t *testing.T, packages map[string]types.Package, path string) {
	t.Helper()

	packagesBytes, err := json.Marshal(packages)
	if err != nil {
		panic(err)
	}

	os.WriteFile(path, packagesBytes, 0644)
}

func TestLoadAndClose(t *testing.T) {
	expectedPackages := loadPackagesJson(t, packagesJson)

	db := Load()

	vars.DB_FILE = tmpPackagesJson
	db.Close()
	vars.DB_FILE = packagesJson

	actualPackages := loadPackagesJson(t, tmpPackagesJson)

	if !cmp.Equal(expectedPackages, actualPackages) {
		t.Errorf("Saved packages are different from expected: %s", cmp.Diff(expectedPackages, actualPackages))
	}
}

func TestCheck(t *testing.T) {
	expPackages := loadPackagesJson(t, packagesJson)

	expPackages["exp2 0.0.1"] = types.Package{
		Name:    "exp2",
		Version: "0.0.1",
	}
	delete(expPackages, "exp2 0.0.0")
	delete(expPackages, "exp3 0.0.0")

	expectedUpToDate := []types.Package{expPackages["exp1 0.0.0"]}
	expectedOutOfDate := []string{"exp2"}
	expectedDeleted := []string{"exp3"}

	db := Load()

	checkParam := []string{}
	for packKey := range expPackages {
		checkParam = append(checkParam, packKey)
	}

	actualUpToDate, actualOutOfDate, actualDeleted := db.Check(checkParam)

	if !cmp.Equal(expectedUpToDate, actualUpToDate) {
		t.Errorf("Up to date packages are different from expected: %s", cmp.Diff(expectedUpToDate, actualUpToDate))
	}

	if !cmp.Equal(expectedOutOfDate, actualOutOfDate) {
		t.Errorf("Out of date packages are different from expected: %s", cmp.Diff(expectedOutOfDate, actualOutOfDate))
	}

	if !cmp.Equal(expectedDeleted, actualDeleted) {
		t.Errorf("Deleted packages are different from expected: %s", cmp.Diff(expectedDeleted, actualDeleted))
	}
}

func TestUpdate(t *testing.T) {
	db := Load()

	exp2Pack := types.Package{
		Name:       "exp2",
		Version:    "0.0.1",
		IsExplicit: true,
		Size:       10240,
	}
	updated := []types.Package{
		exp2Pack,
		{
			Name: "exp3", // Deleted package
		},
	}

	expectedPackages := loadPackagesJson(t, packagesJson)
	expectedPackages["exp2 0.0.1"] = exp2Pack
	delete(expectedPackages, "exp2 0.0.0")
	delete(expectedPackages, "exp3 0.0.0")

	db.Update(updated...)

	vars.DB_FILE = tmpPackagesJson
	db.Close()
	vars.DB_FILE = packagesJson

	actualPackages := loadPackagesJson(t, tmpPackagesJson)

	if !cmp.Equal(expectedPackages, actualPackages) {
		t.Errorf("Packaged are different from expected: %s", cmp.Diff(expectedPackages, actualPackages))
	}
}
