package utils

import (
	"encoding/json"
	"os"
)

func LoadJson[T any](file string) (T, error) {
	fileBytes, err := os.ReadFile(file)
	if err != nil {
		return *new(T), err
	}

	var finalData T
	if err := json.Unmarshal(fileBytes, &finalData); err != nil {
		return *new(T), err
	}

	return finalData, nil
}

func SaveJson(data any, file string) error {
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}

	if err := os.WriteFile(file, dataBytes, 0644); err != nil {
		return err
	}

	return nil
}
