package data

import (
	"os"

	"github.com/spectronp/sizr/tests"
	"github.com/spectronp/sizr/types"
	"github.com/spectronp/sizr/utils"
	"github.com/spectronp/sizr/vars"

	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestMain(m *testing.M) {
	vars.BASEDIR = os.Getenv("BASEDIR")
	vars.DB_FILE = "/tmp/sizr_db"
	defer os.Remove("/tmp/sizr_db")
	if err := utils.SaveJson(map[string]any{}, vars.DB_FILE); err != nil {
		panic(err)
	}

	m.Run()
}

func TestCanBuildData(t *testing.T) {
	_, err := NewData(helpers.MockRunner)

	if err != nil {
		t.Error("Data could not be built: ", err)
	}
}

func TestManagerIsSet(t *testing.T) {
	data, _ := NewData(helpers.MockRunner)

	if data.manager != "pm" {
		t.Errorf("data.manager should be 'pm', is '%s' instead", data.manager)
	}
}

func TestGetPackage(t *testing.T) {
	data, _ := NewData(helpers.MockRunner)

	expectedPackage := types.Package{Name: "exp1", Size: 10240, IsExplicit: true, Version: "0.0.0", Deps: []string{"rev1", "dep1", "rev4"}}
	actualPackage := data.GetPackage("exp1")

	if !cmp.Equal(actualPackage, expectedPackage) {
		diff := cmp.Diff(actualPackage, expectedPackage)
		t.Errorf("Package exp1 is not what was expected\n%s", diff)
	}
}

func TestGetExplicit(t *testing.T) {
	expectedPackages := map[string]types.Package{
		"exp1": {
			Name:       "exp1",
			Size:       10240,
			IsExplicit: true,
			Version:    "0.0.0",
			Deps:       []string{"rev1", "dep1", "rev4"},
		},
		"exp2": {
			Name:       "exp2",
			Size:       10240,
			IsExplicit: true,
			Version:    "0.0.0",
			Deps:       []string{"dep13", "dep14", "rev4"},
		},
		"exp3": {
			Name:       "exp3",
			Size:       10240,
			IsExplicit: true,
			Version:    "0.0.0",
			Deps:       []string{"dep11", "dep12"},
		},
	}

	data, _ := NewData(helpers.MockRunner)

	actualPackages := data.GetExplicit()

	if !cmp.Equal(expectedPackages, actualPackages) {
		fmt.Print(cmp.Diff(expectedPackages, actualPackages))
		t.Errorf("Returned packages are different than expected")
	}
}
