package rotate

func Rotate(s []int, n int) {
	right := s[:n]
	left := s[n:]
	copy(s, append(left, right...))
}
