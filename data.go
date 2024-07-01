package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"
)

type Package struct { // TODO -- make fields exported and immutable
	Name string
	IsExplicit bool
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

func getPackages(manager string, runner ScriptRunner) map[string]Package {
	// TODO -- cache system

	packages := make(map[string]Package)
	raw_result, err := runner(manager + "/get-all")
	if err != nil {
		fmt.Printf("Error on getPackages: %s \n", err)
	}
	packagesNames := strings.Fields(raw_result)
	for _, pack := range packagesNames {
		var newPack Package	
		packageJson, err := runner(manager + "/info", pack)
		if err != nil {
			fmt.Printf("Error on running info.sh: %s\n", err)
		}
		err = json.Unmarshal([]byte(packageJson), &newPack)	
		if err != nil {
			fmt.Printf("Error on Unmarshal: %s\n", err)
		}
		packages[newPack.Name] = newPack
	}
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
