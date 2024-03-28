package utils

func MergeSlices(slices ...[]string) []string {
	var merged []string

	for _, slice := range slices {
		merged = append(merged, slice...)
	}

	return merged
}
