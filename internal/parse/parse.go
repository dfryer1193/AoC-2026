// Package parse has small helpers for turning puzzle input into values.
package parse

import (
	"regexp"
	"strconv"
	"strings"
)

// Lines splits s into lines.
func Lines(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n")
}

// Blocks splits s into blank-line-separated sections.
func Blocks(s string) []string {
	if s == "" {
		return nil
	}
	return strings.Split(s, "\n\n")
}

var intRE = regexp.MustCompile(`-?\d+`)

// Ints extracts every (optionally negative) integer in s, ignoring everything else.
func Ints(s string) []int {
	matches := intRE.FindAllString(s, -1)
	out := make([]int, len(matches))
	for i, m := range matches {
		out[i] = Atoi(m)
	}
	return out
}

// IntLines returns Ints for each line of s.
func IntLines(s string) [][]int {
	lines := Lines(s)
	out := make([][]int, len(lines))
	for i, l := range lines {
		out[i] = Ints(l)
	}
	return out
}

// Atoi is strconv.Atoi that panics on bad input. Puzzle input is trusted.
func Atoi(s string) int {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		panic(err)
	}
	return n
}

// Digits returns the decimal digits of s as ints, skipping non-digits.
func Digits(s string) []int {
	var out []int
	for _, r := range s {
		if r >= '0' && r <= '9' {
			out = append(out, int(r-'0'))
		}
	}
	return out
}
