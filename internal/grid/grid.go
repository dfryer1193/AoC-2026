// Package grid has 2D points and a byte grid for map-style puzzles.
// Coordinates are (X, Y) with Y growing downward, matching how input is read.
package grid

import (
	"strings"
)

type Point struct{ X, Y int }

func P(x, y int) Point { return Point{x, y} }

func (p Point) Add(q Point) Point     { return Point{p.X + q.X, p.Y + q.Y} }
func (p Point) Sub(q Point) Point     { return Point{p.X - q.X, p.Y - q.Y} }
func (p Point) Scale(k int) Point     { return Point{p.X * k, p.Y * k} }
func (p Point) RotateCW() Point       { return Point{-p.Y, p.X} }
func (p Point) RotateCCW() Point      { return Point{p.Y, -p.X} }
func (p Point) Manhattan(q Point) int { return abs(p.X-q.X) + abs(p.Y-q.Y) }

var (
	Up    = Point{0, -1}
	Down  = Point{0, 1}
	Left  = Point{-1, 0}
	Right = Point{1, 0}

	// Dirs4 is the orthogonal directions, clockwise from Up.
	Dirs4 = []Point{Up, Right, Down, Left}
	// Dirs8 is all eight neighbors, clockwise from Up.
	Dirs8 = []Point{Up, {1, -1}, Right, {1, 1}, Down, {-1, 1}, Left, {-1, -1}}
)

// Grid is a rectangular grid of bytes, indexed Cells[y][x].
type Grid struct {
	W, H  int
	Cells [][]byte
}

// Parse builds a Grid from newline-separated rows.
func Parse(s string) *Grid {
	rows := strings.Split(s, "\n")
	g := &Grid{H: len(rows), Cells: make([][]byte, len(rows))}
	for y, row := range rows {
		g.Cells[y] = []byte(row)
		g.W = max(g.W, len(row))
	}
	return g
}

// New makes a w×h grid filled with fill.
func New(w, h int, fill byte) *Grid {
	g := &Grid{W: w, H: h, Cells: make([][]byte, h)}
	for y := range g.Cells {
		g.Cells[y] = []byte(strings.Repeat(string(fill), w))
	}
	return g
}

func (g *Grid) In(p Point) bool {
	return p.Y >= 0 && p.Y < g.H && p.X >= 0 && p.X < len(g.Cells[p.Y])
}

// At returns the byte at p, or 0 if p is out of bounds.
func (g *Grid) At(p Point) byte {
	if !g.In(p) {
		return 0
	}
	return g.Cells[p.Y][p.X]
}

func (g *Grid) Set(p Point, b byte) { g.Cells[p.Y][p.X] = b }

// Find returns the first point holding b, scanning row by row.
func (g *Grid) Find(b byte) (Point, bool) {
	for y, row := range g.Cells {
		for x, c := range row {
			if c == b {
				return Point{x, y}, true
			}
		}
	}
	return Point{}, false
}

// FindAll returns every point holding b.
func (g *Grid) FindAll(b byte) []Point {
	var out []Point
	for y, row := range g.Cells {
		for x, c := range row {
			if c == b {
				out = append(out, Point{x, y})
			}
		}
	}
	return out
}

// Neighbors returns the in-bounds points at p+d for each d in dirs.
func (g *Grid) Neighbors(p Point, dirs []Point) []Point {
	out := make([]Point, 0, len(dirs))
	for _, d := range dirs {
		if q := p.Add(d); g.In(q) {
			out = append(out, q)
		}
	}
	return out
}

// All iterates over every point and its value.
func (g *Grid) All(yield func(Point, byte) bool) {
	for y, row := range g.Cells {
		for x, c := range row {
			if !yield(Point{x, y}, c) {
				return
			}
		}
	}
}

func (g *Grid) Clone() *Grid {
	c := &Grid{W: g.W, H: g.H, Cells: make([][]byte, g.H)}
	for y, row := range g.Cells {
		c.Cells[y] = append([]byte(nil), row...)
	}
	return c
}

func (g *Grid) String() string {
	var b strings.Builder
	for y, row := range g.Cells {
		if y > 0 {
			b.WriteByte('\n')
		}
		b.Write(row)
	}
	return b.String()
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}
