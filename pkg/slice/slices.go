package slice

func ReverseSliceInPlace[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}

func ReverseSlice[S ~[]E, E any](s S) S {
	new := make([]E, len(s))
	copy(new, s)
	ReverseSliceInPlace(new)
	return new
}
