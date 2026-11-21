// Package aoctest has helpers for checking solutions against puzzle examples.
package aoctest

import (
	"fmt"
	"strings"
	"testing"

	"github.com/dfryer1193/AoC-2026/internal/runner"
)

// Check runs solve on input and compares the result to want using their %v
// representations, so an int answer can be checked against 42 or "42".
// The test is skipped when want is nil or the input is empty.
func Check(t *testing.T, solve runner.Solver, input string, want any) {
	t.Helper()
	if want == nil {
		t.Skip("no expected answer yet")
	}
	input = runner.Normalize(input)
	if strings.TrimSpace(input) == "" {
		t.Skip("no example input yet")
	}
	got := solve(input)
	if fmt.Sprint(got) != fmt.Sprint(want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
