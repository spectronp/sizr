package db

import (
	"fmt"
	"os"
	"strings"

	"github.com/google/go-cmp/cmp"
	"github.com/spectronp/sizr/types"
	"github.com/spectronp/sizr/vars"

	"encoding/json"
)

type DB struct {
	keyMap     map[string]types.Package // "packName Version": Package
	pointerMap map[string]string        // "packName": "packName Version"
}

func Load() DB {
	jsonFile, err := os.ReadFile(vars.DB_FILE)

	if os.IsNotExist(err) {
		jsonFile = []byte("{}")
		os.WriteFile(vars.DB_FILE, jsonFile, 0644)
		err = nil
	}

	if err != nil {
		panic(err)
	}

	var localKeyMap map[string]types.Package
	json.Unmarshal(jsonFile, &localKeyMap)

	localPointerMap := map[string]string{}
	for pack := range localKeyMap {
		packName, _, _ := strings.Cut(pack, " ")
		localPointerMap[packName] = pack
	}

	return DB{
		keyMap:     localKeyMap,
		pointerMap: localPointerMap,
	}
}

func (db DB) Close() {
	jsonData, err := json.Marshal(db.keyMap)
	if err != nil {
		panic("error on json.Marshal() for db")
	}
	os.WriteFile(vars.DB_FILE, jsonData, 0644)
}

func (db DB) Check(packagesKey []string) (upToDate []types.Package, outOfDate []string, deleted []string) {
	deletedMap := map[string]bool{}
	for _, value := range db.keyMap {
		deletedMap[value.Name] = true
	}

	for _, key := range packagesKey {
		pack, presentInMap := db.keyMap[key]
		if presentInMap {
			upToDate = append(upToDate, pack)
			delete(deletedMap, pack.Name)
		} else {
			packName := strings.Split(key, " ")[0]
			outOfDate = append(outOfDate, packName)
			delete(deletedMap, packName)
		}
	}

	for deletedName := range deletedMap {
		deleted = append(deleted, deletedName)
	}

	return
}

func (db DB) Update(updated ...types.Package) {
	for _, pack := range updated {
		if cmp.Equal(pack, types.Package{Name: pack.Name}) {
			packKey := db.pointerMap[pack.Name]
			delete(db.pointerMap, pack.Name)
			delete(db.keyMap, packKey)

			continue
		}

		oldPackKey := db.pointerMap[pack.Name]
		newPackKey := fmt.Sprintf("%s %s", pack.Name, pack.Version)
		delete(db.keyMap, oldPackKey)
		db.keyMap[newPackKey] = pack
		db.pointerMap[pack.Name] = newPackKey
	}
}
