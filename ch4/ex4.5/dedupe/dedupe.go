package dedupe

func Dedupe(strs []string) []string {
	if len(strs) == 0 {
		return strs[:]
	}

	index := 0
	for _, s := range strs {
		if strs[index] != s {
			index++
			strs[index] = s
		}
	}
	return strs[:index+1]
}
