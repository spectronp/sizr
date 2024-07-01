package main

import (
	"fmt"
	"maps"
	"sort"
)

// TODO -- use concurrency | paralelism
func calcSize(pack string, data Data) uint {
	target := data.GetPackage(pack)
	explictPackages := data.GetExplicit() // TODO -- fix typo explict to explicit
	delete(explictPackages, pack)

	depsToIgnore := make(map[string]Package) // NOTE -- could I use map[string]bool ?
	for _, explicitPack := range explictPackages {
		maps.Copy(depsToIgnore, listTree(explicitPack, data))
	} 
	return sumSize(target, depsToIgnore, data)
}

func listTree(target Package, data Data) map[string]Package {
	finalMap := make(map[string]Package)

	for _, packName := range target.Deps {
		pack := data.GetPackage(packName)
		finalMap[pack.Name] = pack

		maps.Copy(finalMap, listTree(pack, data))
	}
	return finalMap
}

/*
	NOTE -- maybe use listTree first with an ignorePackages parameter and then iterate the list doing the sum ?
*/
func sumSize(start Package, ignorePackages map[string]Package, data Data) uint {
	totalSize := start.Size

	for _, packName := range start.Deps {
		if (ignorePackages[packName].Name == packName){ // TODO -- refactor this line
			continue
		}
		pack := data.GetPackage(packName)
		totalSize += sumSize(pack, ignorePackages, data)
	}
	return totalSize
}

type PackageNameWithSum struct {
	Name string
	Size uint
}

func orderBySumSize(data Data) []PackageNameWithSum {
	explicitPackages := data.GetExplicit()	
	orderedPacks := []PackageNameWithSum{}
	for _, pack := range explicitPackages {
		packSize := calcSize(pack.Name, data)		
		insertIndex := sort.Search(len(orderedPacks), func(i int) bool { return orderedPacks[i].Size <= packSize })
		
		if insertIndex == len(orderedPacks) {
			orderedPacks = append(orderedPacks, PackageNameWithSum{ Name: pack.Name, Size: packSize })
			continue
		}
		orderedPacks = append(orderedPacks[:insertIndex+1], orderedPacks[insertIndex:]...) // NOTE -- should i make more readable ?
		orderedPacks[insertIndex] = PackageNameWithSum{ Name: pack.Name, Size: packSize }
	}
	return orderedPacks
}

func report(packages []PackageNameWithSum, limit uint8) {
	if limit == 0 {
		limit = 30
	}

	// TODO -- display human readable size

	var i uint8
	for i = 0; i <= limit; i++ {
		if i >= uint8(len(packages)) {
			break
		}
		fmt.Printf("%s\t%d\n", packages[i].Name, packages[i].Size)
	}
}

func Run() int {
	// arg parse
	data, _ := NewData(RunScript)
	// CLI or TUI

	report(orderBySumSize(data), 0)

	return 0
}

func main() {
	Run()
}
