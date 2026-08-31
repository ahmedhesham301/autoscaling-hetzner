package services

import (
	"fmt"
	"strings"
)

func ConvertToHetznerLabels(labels map[string]any) string {
	stringLabels := []string{}
	for key, value := range labels {
		stringLabels = append(stringLabels, fmt.Sprintf("%s=%v", key, value))
	}
	return strings.Join(stringLabels, ",")
}
