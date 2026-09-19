package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strconv"
)

func ConvertMapToJsonString(m map[string]any) string {
	jsonBytes, err := json.Marshal(m)
	if err != nil {
		panic(err)
	}
	return string(jsonBytes)
}

func GetSnapshotID(output []byte) (int64, error) {
	marker := []byte(",artifact,0,id,")

	index := bytes.Index(output, marker)
	if index == -1 {
		return 0, fmt.Errorf("snapshot ID not found")
	}

	start := index + len(marker)

	end := bytes.IndexByte(output[start:], '\n')
	if end == -1 {
		end = len(output[start:])
	}

	id, err := strconv.ParseInt(string(output[start:start+end]), 10, 64)
	return id, err
}
