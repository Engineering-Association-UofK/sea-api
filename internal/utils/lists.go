package utils

// A function that takes in a slice of structs and returns a slice of one of it's fields using a function argument
func ExtractField[T any, R any](slice []T, fn func(T) R) []R {
	result := make([]R, len(slice))
	for i, item := range slice {
		result[i] = fn(item)
	}
	return result
}
