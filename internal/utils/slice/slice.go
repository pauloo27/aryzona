package slice

func Map[S any, T any](arr []S, fn func(S) T) []T {
	targetArr := make([]T, len(arr))
	for i, s := range arr {
		targetArr[i] = fn(s)
	}
	return targetArr
}

func Find[S any](arr []S, fn func(S) bool) *S {
	for _, s := range arr {
		if fn(s) {
			return &s
		}
	}
	return nil
}
