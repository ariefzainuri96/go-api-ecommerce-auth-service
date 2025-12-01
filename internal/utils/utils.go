package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"
)

func UpdateStruct(existing interface{}, updates interface{}) {
	existingVal := reflect.ValueOf(existing).Elem()
	updatesVal := reflect.ValueOf(updates)

	for i := range updatesVal.NumField() {
		field := updatesVal.Field(i)

		// Skip zero values (default values)
		if !field.IsZero() {
			existingField := existingVal.Field(i)

			// Ensure the field can be set (ignores unexported fields)
			if existingField.CanSet() {
				existingField.Set(field)
			}
		}
	}
}

func FilterSlice[T any](items []T, filterFunc func(T) bool) []T {
	result := make([]T, 0, len(items)) // Preallocate memory

	for _, item := range items {
		if filterFunc(item) {
			result = append(result, item)
		}
	}
	return result
}

func MapSlice[T any, R any](source []T, transformer func(T) R) []R {
	// 1. Fail Fast: Handle nil or empty input gracefully
	if source == nil {
		return nil
	}
	if len(source) == 0 {
		return []R{}
	}

	// 2. Initialize the result slice with the correct capacity
	result := make([]R, len(source))

	// 3. Loop and transform
	for i, item := range source {
		result[i] = transformer(item)
	}

	return result
}

// LoadJsonData loads JSON data from a file into the provided destination.
func LoadJsonData[T any](filePath string, data *T) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("error opening file: %w", err)
	}
	defer file.Close()

	jsonParser := json.NewDecoder(file)

	if err := jsonParser.Decode(data); err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}

	fmt.Println("Successfully loaded data from", filePath)
	return nil
}

// SaveToJson saves a given data structure to a JSON file.
func SaveToJson[T any](filePath string, data T) error {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return fmt.Errorf("error marshalling JSON: %w", err)
	}

	file, err := os.Create(filePath)
	if err != nil {
		return fmt.Errorf("error creating file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(jsonData)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	fmt.Println("Successfully saved data to", filePath)
	return nil
}
