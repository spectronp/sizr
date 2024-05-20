package main

import (
	"fmt"
	"maps"
)

// TODO -- use concurrency | paralelism
func calcSize(pack string, data Data) uint {
	target := data.GetPackage(pack)
	explictPackages := data.GetExplicit()
	delete(explictPackages, pack)

	depsToIgnore := make(map[string]Package) // NOTE -- could I use map[string]bool ?
	for _, explicitPack := range explictPackages {
		maps.Copy(depsToIgnore, listTree(explicitPack, data))
	} 
	return sumSize(target, depsToIgnore, data)
}

func listTree(target Package, data Data) map[string]Package {
	finalMap := make(map[string]Package)

	for _, packName := range target.deps {
		pack := data.GetPackage(packName)
		finalMap[pack.name] = pack

		maps.Copy(finalMap, listTree(pack, data))
	}
	return finalMap
}

/*
	NOTE -- maybe use listTree first with an ignorePackages parameter and then iterate the list doing the sum ?
*/
func sumSize(start Package, ignorePackages map[string]Package, data Data) uint {
	totalSize := start.size
	
	for _, packName := range start.deps {
		if (ignorePackages[packName].name == packName){ // TODO -- refactor this line
			continue
		}
		pack := data.GetPackage(packName)
		totalSize += sumSize(pack, ignorePackages, data)
	}
	return totalSize
}

func main()  {
	Run()
}

func Run()  {
	data, _ := NewData(RunScript)
	fmt.Println(data.PackageList["steam"].size)

	// `
}
