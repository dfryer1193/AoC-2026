package parse

import (
	"slices"
	"testing"
)

func TestInts(t *testing.T) {
	got := Ints("p=3,-14 v=-2,7 x42")
	want := []int{3, -14, -2, 7, 42}
	if !slices.Equal(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestBlocks(t *testing.T) {
	got := Blocks("a\nb\n\nc")
	if len(got) != 2 || got[0] != "a\nb" || got[1] != "c" {
		t.Errorf("got %q", got)
	}
}

func TestLinesEmpty(t *testing.T) {
	if got := Lines(""); got != nil {
		t.Errorf("got %q, want nil", got)
	}
}
