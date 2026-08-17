package services

func MapStructToSchema[I any, O any](items []I, mapper func(I) O) []O {
	result := make([]O, 0, len(items))

	for _, item := range items {
		result = append(result, mapper(item))
	}

	return result
}

func AppendManagedLabel(labels map[string]string) map[string]string {
	if labels == nil {
		labels = make(map[string]string)
	}
	labels["managed"] = "true"
	return labels
}
