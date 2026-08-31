package utils

import "encoding/json"

func ConvertMapToJsonString(m map[string]any) string {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}
