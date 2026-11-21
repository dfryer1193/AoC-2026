package grid

import "testing"

func TestGrid(t *testing.T) {
	g := Parse("#.S\n..#")
	if g.W != 3 || g.H != 2 {
		t.Fatalf("size %dx%d", g.W, g.H)
	}
	s, ok := g.Find('S')
	if !ok || s != P(2, 0) {
		t.Fatalf("Find = %v, %v", s, ok)
	}
	if n := g.Neighbors(s, Dirs4); len(n) != 2 {
		t.Errorf("neighbors of corner: %v", n)
	}
	if g.At(P(5, 5)) != 0 {
		t.Error("out of bounds At should be 0")
	}
	if g.String() != "#.S\n..#" {
		t.Errorf("String = %q", g.String())
	}
}

func TestRotate(t *testing.T) {
	if Up.RotateCW() != Right || Up.RotateCCW() != Left {
		t.Error("rotation is wrong")
	}
}
