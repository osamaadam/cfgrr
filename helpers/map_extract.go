package helpers

// Returns the keys of a map as a slice.
func GetMapKeys[K comparable, V interface{}](m map[K]V) []K {
	arr := make([]K, len(m))
	i := 0
	for k := range m {
		arr[i] = k
		i++
	}

	return arr
}

// Returns the values of a map as a slice.
func GetMapValues[K comparable, V interface{}](m map[K]V) []V {
	arr := make([]V, len(m))
	i := 0
	for _, v := range m {
		arr[i] = v
		i++
	}

	return arr
}

// Returns the keys and values of a map as slices.
func GetMapKeysValues[K comparable, V interface{}](m map[K]V) ([]K, []V) {
	keys := make([]K, len(m))
	values := make([]V, len(m))
	i := 0
	for k, v := range m {
		keys[i] = k
		values[i] = v
		i++
	}

	return keys, values
}
