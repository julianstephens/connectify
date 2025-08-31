package utils

// DefaultIfNil returns the default value if the provided pointer is nil
func DefaultIfNil[T any](value *T, def T) T {
	if value == nil {
		return def
	}
	return *value
}

// PointerTo creates a pointer to the given value
func PointerTo[T any](value T) *T {
	return &value
}
