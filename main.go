package main

import (
	"github.com/spectronp/sizr/data"
	"github.com/spectronp/sizr/types"

	"fmt"
	"os"
	"sort"

	"github.com/spf13/pflag"
)

var (
	VERSION      string
	HELP_MESSAGE string
	ENV          string
)

func calcSize(pack string, dataObj *data.Data) uint {
	target := dataObj.GetPackage(pack)
	explictPackages := dataObj.GetExplicit()
	delete(explictPackages, pack)

	depsToIgnore := make(map[string]types.Package) // NOTE: could I use map[string]bool ?
	for _, explicitPack := range explictPackages {
		listTree(explicitPack, depsToIgnore, dataObj)
	}
	return sumSize(target, depsToIgnore, dataObj)
}

func listTree(target types.Package, depsToIgnore map[string]types.Package, dataObj *data.Data) {
	for _, packName := range target.Deps {
		pack := dataObj.GetPackage(packName)
		if _, alreadyInList := depsToIgnore[packName]; alreadyInList {
			continue
		}
		depsToIgnore[pack.Name] = pack // this can add an already added package and run another branch of listTree

		listTree(pack, depsToIgnore, dataObj)
	}
}

/*
NOTE: maybe use listTree first with an ignorePackages parameter and then iterate the list doing the sum ?
*/
func sumSize(start types.Package, ignorePackages map[string]types.Package, dataObj *data.Data) uint {
	totalSize := start.Size

	for _, packName := range start.Deps {
		if _, shouldBeIgnored := ignorePackages[packName]; shouldBeIgnored {
			continue
		}
		pack := dataObj.GetPackage(packName)
		totalSize += sumSize(pack, ignorePackages, dataObj)
	}
	return totalSize
}

type PackageNameWithSum struct {
	Name string
	Size uint
}

func orderBySumSize(dataObj *data.Data) []PackageNameWithSum {
	explicitPackages := dataObj.GetExplicit()
	orderedPacks := []PackageNameWithSum{}
	for _, pack := range explicitPackages {
		packSize := calcSize(pack.Name, dataObj)
		insertIndex := sort.Search(len(orderedPacks), func(i int) bool { return orderedPacks[i].Size <= packSize })

		if insertIndex == len(orderedPacks) {
			orderedPacks = append(orderedPacks, PackageNameWithSum{Name: pack.Name, Size: packSize})
			continue
		}
		orderedPacks = append(orderedPacks[:insertIndex+1], orderedPacks[insertIndex:]...) // NOTE: should i make more readable ?
		orderedPacks[insertIndex] = PackageNameWithSum{Name: pack.Name, Size: packSize}
	}
	return orderedPacks
}

func report(packages []PackageNameWithSum, limit uint8) {
	// TODO: display human readable size
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
		gapNum := 5 + (longestPackageName - len([]rune(pack.Name)))
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

	dataObject, _ := data.NewData(data.RunScript)
	// CLI or TUI

	report(orderBySumSize(&dataObject), *reportLimit)

	return 0
}

func main() {
	statusCode := Run(os.Args)
	os.Exit(statusCode)
}
