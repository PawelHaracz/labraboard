package helpers

func Contains[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Remove[T comparable](s []T, v T) []T {
	i := FindIndex(s, v)
	if i == -1 {
		return s
	}
	return append(s[:i], s[i+1:]...)
}

func FindIndex[T comparable](s []T, v T) int {
	for i, vs := range s {
		if vs == v {
			return i
		}
	}
	return -1
}
