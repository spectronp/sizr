package main

import (
	"fmt"
	"os"
	"sort"

	"github.com/spf13/pflag"
)

var (
	VERSION string
	HELP_MESSAGE string
	ENV string
)

// TODO -- use concurrency | paralelism
func calcSize(pack string, data *Data) uint {
	target := data.GetPackage(pack)
	explictPackages := data.GetExplicit()
	delete(explictPackages, pack)

	depsToIgnore := make(map[string]Package) // NOTE -- could I use map[string]bool ?
	for _, explicitPack := range explictPackages {
		listTree(explicitPack, depsToIgnore, data)
	} 
	return sumSize(target, depsToIgnore, data)
}

func listTree(target Package, depsToIgnore map[string]Package, data *Data) {
	for _, packName := range target.Deps {
		pack := data.GetPackage(packName)
		if _, alreadyInList := depsToIgnore[packName]; alreadyInList {
			continue
		} 
		depsToIgnore[pack.Name] = pack // this can add an already added package and run another branch of listTree

		listTree(pack, depsToIgnore, data)
	}
}

/*
	NOTE -- maybe use listTree first with an ignorePackages parameter and then iterate the list doing the sum ?
*/
func sumSize(start Package, ignorePackages map[string]Package, data *Data) uint {
	totalSize := start.Size

	for _, packName := range start.Deps {
		if _, shouldBeIgnored := ignorePackages[packName]; shouldBeIgnored {
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

func orderBySumSize(data *Data) []PackageNameWithSum {
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
	// TODO -- display human readable size
	var limitedPackages []PackageNameWithSum
	if uint8(len(packages)) > limit {
		limitedPackages = packages[:limit]
	} else {
		limitedPackages = packages
	}
	
	longestPackageName := 0
	for _, pack := range limitedPackages {
		runeCount := len([]rune(pack.Name))
		if runeCount > longestPackageName {
			longestPackageName = runeCount
		}
	}
	
	for _, pack := range limitedPackages {
		gapNum := 5 + ( longestPackageName - len([]rune(pack.Name)) )
		gapString := ""
		for i := gapNum; i > 0; i-- {
			gapString += " "	
		}

		fmt.Printf("%s%s%d\n", pack.Name, gapString, pack.Size)
	}
}

func Run(args []string) int {
	ENV = os.Getenv("SIZR_ENV")

	flag := pflag.NewFlagSet("", 1)

	helpWanted := flag.BoolP("help", "h", false, "Show this help message")
	wantVersion := flag.BoolP("version", "v", false, "Show sizr version")
	reportLimit := flag.Uint8P("limit", "n", 30, "Set the limit of packages to show (Default: 30)")
	flag.Parse(args)

	if *helpWanted {
		fmt.Print(HELP_MESSAGE)
		return 0
	}
	if *wantVersion {
		fmt.Printf("sizr %s\n", VERSION)
		return 0
	}

	data, _ := NewData(RunScript)
	// CLI or TUI

	report(orderBySumSize(&data), *reportLimit)

	return 0
}

func main() {
	Run(os.Args)
	// NOTE -- use os.Exit() with return code here ?
}
