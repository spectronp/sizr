package main

import (
	"fmt"
	"maps"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
)

type Package struct {
	name string
	isExplict bool
	size uint
	deps []string // NOTE -- change this to []*Package or map[string]*Packge ?
}


type Data struct {
	Manager string
	PackageList map[string]Package // NOTE -- change name to Packages or PackageMap ?
}

type ScriptRunner func(script string, args ...string) (output string, err error)

func NewData(runner ScriptRunner) (data Data, err error) {
	manager, err := runner("get-package-manager")
	packageList := getPackages(manager, runner)
	return Data{Manager: manager, PackageList: packageList}, err
} 

func getPackages(manager string, runner ScriptRunner) map[string]Package {
	// TODO -- cache system

	packages := make(map[string]Package)
	raw_result, _ := runner(manager + "/get-all")
	packagesName := strings.Fields(raw_result)

	for _, packName := range packagesName { // TODO -- too slow
		raw_info, _ := runner(manager + "/info", packName)
		info := strings.Fields(raw_info)
		
		isExplict, err := strconv.ParseBool(info[1])
		if err != nil {
			fmt.Println("erro no explictic")		
			fmt.Println(err)
		}	

		size, err := strconv.ParseUint(info[2], 10, 0)
		if err != nil {
			fmt.Println("erro no size")
			fmt.Println(err)
		}	

		deps := strings.Fields(info[3])

		pack := Package{
			info[0],
			isExplict,
			uint(size),
			deps,	
		}
		packages[packName] = pack	
		
		fmt.Printf("%v %v", packName, size)
		fmt.Println("")
	} 
	return packages
}

func (d Data) GetPackage(name string) Package  {
	pack := d.PackageList[name]
	return pack
}

func (d Data) GetExplicit() map[string]Package { // TODO -- try to use map[string]*Package
	packs := d.PackageList // TODO -- check if this is not changing the data field itself ( maybe use maps.Clone instead )

	maps.DeleteFunc(packs, func (_ string, pack Package) bool {
		return  !pack.isExplict
	})
	return packs
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
