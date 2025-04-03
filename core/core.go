package core

func levenshtein(s1, s2 string) int {
	if len(s1) == 0 {
		return len(s2)
	}
	if len(s2) == 0 {
		return len(s1)
	}

	if s1[0] == s2[0] {
		return levenshtein(s1[1:], s2[1:])
	}

	l1 := levenshtein(s1[1:], s2)
	l2 := levenshtein(s1, s2[1:])
	l3 := levenshtein(s1[1:], s2[1:])

	min := l1
	if l2 < min {
		min = l2
	}
	if l3 < min {
		min = l3
	}
	return min + 1
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func LevenshteinSimilarity(s1, s2 string) float64 {
	distance := levenshtein(s1, s2)
	maxLen := max(len(s1), len(s2))
	if maxLen == 0 {
		return 1.0
	}
	return 1.0 - float64(distance)/float64(maxLen)
}
