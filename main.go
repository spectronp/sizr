package main

import (
	"github.com/dustin/go-humanize"
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

	depsToIgnore := make(map[string]bool)
	for _, explicitPack := range explictPackages {
		listTree(explicitPack, depsToIgnore, dataObj)
	}
	alreadyCounted := map[string]bool{}
	return sumSize(target, depsToIgnore, alreadyCounted, dataObj)
}

func listTree(target types.Package, depsToIgnore map[string]bool, dataObj *data.Data) {
	for _, packName := range target.Deps {
		pack := dataObj.GetPackage(packName)
		if depsToIgnore[packName] {
			continue
		}
		depsToIgnore[pack.Name] = true

		listTree(pack, depsToIgnore, dataObj)
	}
}

func sumSize(start types.Package, ignorePackages map[string]bool, alreadyCounted map[string]bool, dataObj *data.Data) uint {
	totalSize := start.Size

	for _, packName := range start.Deps {
		if ignorePackages[packName] {
			continue
		}
		if alreadyCounted[packName] {
			continue
		}

		pack := dataObj.GetPackage(packName)
		alreadyCounted[packName] = true
		totalSize += sumSize(pack, ignorePackages, alreadyCounted, dataObj)
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
		orderedPacks = append(orderedPacks[:insertIndex+1], orderedPacks[insertIndex:]...)
		orderedPacks[insertIndex] = PackageNameWithSum{Name: pack.Name, Size: packSize}
	}
	return orderedPacks
}

func report(packages []PackageNameWithSum, limit uint8) {
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

		fmt.Printf("%s%s%s\n", pack.Name, gapString, humanize.IBytes(uint64(pack.Size)))
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
