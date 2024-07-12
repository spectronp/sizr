package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/schollz/progressbar/v3"
)

type Package struct { // TODO -- make fields exported and immutable
	Name string
	IsExplicit bool
	Version string
	Size uint
	Deps []string // NOTE -- change this to []*Package or map[string]*Packge ?
}


type Data struct {
	Manager string
	PackageList map[string]Package // NOTE -- change name to Packages or PackageMap ? | Should this be private ?
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

func getPackageInfoWorker(packageNames <-chan string, returnedPacks chan<- Package, runner ScriptRunner, manager string) {
	for packName := range packageNames {
		var newPack Package	
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

func getPackages(manager string, runner ScriptRunner) map[string]Package {
	// TODO -- cache system

	packages := make(map[string]Package)
	raw_result, err := runner(manager + "/get-all")
	if err != nil {
		fmt.Printf("Error on getPackages: %s \n", err) // TODO -- use log for errors
	}
	packagesInfo := strings.Split(raw_result, "\n")

	DB := DB.Load()

	// check for names in the DB and get packages that need to be updated
	upTodate, outOfDate, err := DB.Check(packagesInfo)
	for pack := range upTodate {
		packages[pack.Name] = pack
	}

	packagesCount := len(outOfDate)
	
	namesChannel := make(chan string, packagesCount)
	packagesChannel := make(chan Package, packagesCount)

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

	for w := 1; w <= packagesCount; w++ { // NOTE -- should the progress bar output be here or in the main.go ?
		newPack := <-packagesChannel
		packages[newPack.Name] = newPack
		bar.Add(1)
	} 

	if ENV == "testing" {
		fmt.Println("@END_PROGRESSBAR@")	
	}

	DB.Update(packages, outOfDate)
	
	return packages
}

func (d Data) GetPackage(name string) Package  {
	pack := d.PackageList[name]
	return pack
}

func (d Data) GetExplicit() map[string]Package { // TODO -- try to use map[string]*Package
	explicit := make(map[string]Package)

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
