package main

import (
	"regexp"
	"strings"
)

// Remove all non-alphanumeric characters and split into tokens.
func tokenize(text string) (tokens []string) {
	re := regexp.MustCompile(`[^a-zA-Z0-9]+`)
	for _, t := range re.Split(text, -1) {
		if t != "" {
			tokens = append(tokens, strings.ToLower(t))
		}
	}
	return
}

// Remove all non-alphanumeric and wildcard characters and split into tokens.
func tokenizeWildcard(text string) (tokens []string) {
	re := regexp.MustCompile(`[^a-zA-Z0-9*?]+`)
	for _, t := range re.Split(text, -1) {
		if t != "" {
			tokens = append(tokens, strings.ToLower(t))
		}
	}
	return
}

// Calculates edit (levenshtein) distance between two strings.
func editDistance(s1 string, s2 string) int {
	// Initialize empty 2-d array
	m := make([][]int, len(s1)+1)
	_m := make([]int, (len(s1)+1) * (len(s2)+1))
	for i := range m {
		m[i], _m = _m[:len(s2)+1], _m[len(s2)+1:]
	}
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

// Test if the input string matches the wildcard pattern.
func wildcardMatch(pattern string, str string) bool {
	// Initialize empty 2-d array
	m := make([][]bool, len(pattern)+1)
	_m := make([]bool, (len(pattern)+1) * (len(str)+1))
	for i := range m {
		m[i], _m = _m[:len(str)+1], _m[len(str)+1:]
	}
	m[0][0] = true
	for i := 1; i <= len(pattern); i++ {
		m[i][0] = pattern[i-1] == '*' && m[i-1][0]
	}
	for i := 1; i <= len(pattern); i++ {
		for j := 1; j <= len(str); j++ {
			if pattern[i-1] == str[j-1] || pattern[i-1] == '?' {
				m[i][j] = m[i-1][j-1]
			} else if pattern[i-1] == '*' {
				m[i][j] = m[i][j-1] || m[i-1][j]
			} else {
				// pattern[i-1] != str[j-1]
				m[i][j] = false
			}
		}
	}
	return m[len(pattern)][len(str)]
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