package ptr

func From[T any](v *T) T {
	if v == nil {
		return *new(T)
	}
	return *v
}
