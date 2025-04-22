package utils

import "encoding/json"

func ConvertToMap[T any](data T) (map[string]interface{}, error) {
	var result map[string]interface{}

	// Konversi struct ke JSON bytes
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	// Konversi JSON bytes ke map
	err = json.Unmarshal(jsonBytes, &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}
