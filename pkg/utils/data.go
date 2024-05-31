package utils

// MergeSlices merges multiple slices of strings into a single slice.
//
// Parameters:
//   - slices: Variadic parameter representing multiple slices of strings to be merged.
//
// Returns:
//   - []string: A single slice containing all the elements from the input slices.
//
// Example:
//
//	slice1 := []string{"a", "b"}
//	slice2 := []string{"c", "d"}
//	merged := MergeSlices(slice1, slice2)
//	fmt.Printf("Merged slice: %v\n", merged)
func MergeSlices(slices ...[]string) []string {
	var merged []string

	for _, slice := range slices {
		merged = append(merged, slice...)
	}

	return merged
}
