package indices

import "regexp"

func Tokenize(text string) (tokens []string) {
	// Remove all non-alphanumeric characters and split into tokens.
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	for _, t := range re.Split(text, -1) {
		if t != "" {
			tokens = append(tokens, t)
		}
	}
	return
}

func StripPunctuation(text string) string {
	// Remove all non-alphanumeric and whitespace characters.
	re := regexp.MustCompile(`[^a-zA-Z0-9 ]+`)
	return re.ReplaceAllString(text, "")
}

func EditDistance(s1 string, s2 string) int {
	// Initialize empty 2-d array
	m := make([][]int, len(s1)+1)
	_m := make([]int, (len(s1)+1) * (len(s2)+1))
	for i := range m {
		m[i], _m = _m[:len(s2)+1], _m[len(s2)+1:]
	}
	// Calculate levenshtein distance
	for i := 1; i <= len(s1); i++ {
		m[i][0] = i
	}
	for j := 1; j <= len(s2); j++ {
		m[0][j] = j
	}
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s2); j++ {
			c := min(m[i-1][j], m[i][j-1]) + 1
			if s1[i-1] == s2[j-1] {
				m[i][j] = min(m[i-1][j-1], c)
			} else {
				m[i][j] = min(m[i-1][j-1] + 1, c)
			}
		}
	}
	return m[len(s1)][len(s2)]
}

func PrefixEditDistance(s1 string, s2 string, d int) int {
	// Computes Prefix Edit Distance where len(s1) << len(s2) and PED(s1, s2) <= d.
	// Initialize empty 2-d array
	m := make([][]int, len(s1)+1)
	_m := make([]int, (len(s1)+1) * (len(s1)+d+1))
	for i := range m {
		m[i], _m = _m[:len(s1)+d+1], _m[len(s1)+d+1:]
	}
	// Calculate levenshtein distance
	for i := 1; i <= len(s1); i++ {
		m[i][0] = i
	}
	for j := 1; j <= len(s1) + d; j++ {
		m[0][j] = j
	}
	for i := 1; i <= len(s1); i++ {
		for j := 1; j <= len(s1) + d; j++ {
			c := min(m[i-1][j], m[i][j-1]) + 1
			if s1[i-1] == s2[j-1] {
				m[i][j] = min(m[i-1][j-1], c)
			} else {
				m[i][j] = min(m[i-1][j-1] + 1, c)
			}
		}
	}
	return min(m[len(s1)]...)
}

func min(a ...int) (m int) {
	m = a[0]
	for _, i := range a {
		if i < m {
		m = i
	}
	}
	return
}

func max(a ...int) (m int) {
	m = a[0]
	for _, i := range a {
		if i > m {
			m = i
		}
	}
	return
}