package data

import (
	"github.com/spectronp/sizr/db"
	"github.com/spectronp/sizr/types"
	"github.com/spectronp/sizr/vars"

	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
)



type Data struct {
	Manager string
	PackageList map[string]types.Package // NOTE -- change name to Packages or PackageMap ? | Should this be private ?
}

type ScriptRunner func(script string, args ...string) (output string, err error)

func NewData(runner ScriptRunner) (data Data, err error) {
	manager, err := runner("get-package-manager")
	if err != nil {
		err = fmt.Errorf("Error on getting package manager: %w\n", err)
	}
	packageList := getPackages(manager, runner)
	return Data{Manager: manager, PackageList: packageList}, err
} 

func getPackageInfoWorker(packageNames <-chan string, returnedPacks chan<- types.Package, runner ScriptRunner, manager string) {
	for packName := range packageNames {
		var newPack types.Package	
		packageJson, err := runner(manager + "/info", packName)
		if err != nil {
			fmt.Printf("Error on running info.sh: %s\n", err)
		}
		err = json.Unmarshal([]byte(packageJson), &newPack)	
		if err != nil {
			fmt.Printf("Error on Unmarshal: %s\n", err)
		}
		returnedPacks <- newPack	
	}
}

func getPackages(manager string, runner ScriptRunner) map[string]types.Package {
	// TODO -- cache system

	packages := make(map[string]types.Package)
	raw_result, err := runner(manager + "/get-all")
	if err != nil {
		fmt.Printf("Error on getPackages: %s \n", err) // TODO -- use log for errors
	}
	packagesInfo := strings.Split(raw_result, "\n")

	DB := db.Load()
	defer DB.Close()

	// check for names in the DB and get packages that need to be updated
	upToDate, outOfDate, deleted := DB.Check(packagesInfo)
	for _, pack := range upToDate {
		packages[pack.Name] = pack
	}

	packagesCount := len(outOfDate)
	
	namesChannel := make(chan string, packagesCount)
	packagesChannel := make(chan types.Package, packagesCount)

	workerCount := 6
	for w := 1; w <= workerCount; w++ {
		go getPackageInfoWorker(namesChannel, packagesChannel, runner, manager)	
	}

	for _, packageName := range outOfDate {
		namesChannel <- packageName
	}
	close(namesChannel)
	
	bar := progressbar.NewOptions(packagesCount,
		progressbar.OptionClearOnFinish(),
		progressbar.OptionShowCount(),
	)
	
	updateOnDB := []types.Package{}
	for w := 1; w <= packagesCount; w++ { // NOTE -- should the progress bar output be here or in the main.go ?
		newPack := <-packagesChannel
		packages[newPack.Name] = newPack
		updateOnDB = append(updateOnDB, newPack)
		bar.Add(1)
	} 

	if vars.ENV == "testing" {
		fmt.Println("@END_PROGRESSBAR@")	
	}

	for _, deletedName := range deleted {
		updateOnDB = append(updateOnDB, types.Package{Name: deletedName})
	}

	DB.Update(updateOnDB...)
	
	return packages
}

func (d Data) GetPackage(name string) types.Package  {
	pack := d.PackageList[name]
	return pack
}

func (d Data) GetExplicit() map[string]types.Package { // TODO -- try to use map[string]*Package
	explicit := make(map[string]types.Package)

	for _, pack := range d.PackageList {
		if pack.IsExplicit {
			explicit[pack.Name] = pack
		}		
	}	

	return explicit
} 

func RunScript(script string, args ...string) (output string, err error) {
	shell := "/bin/sh"
	path := "scripts/" + script + ".sh"	
	path, _ = filepath.Abs(path)
	pathAndArgs := append([]string{path}, args...)
	cmd := exec.Command(shell, pathAndArgs...)
	stdout, err := cmd.Output()
	if (err != nil){
		fmt.Println(err)
	}
	return string(stdout), err
}
