// Package runner provides the entry point shared by every day's solution.
package runner

import (
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// Solver solves one part of a puzzle. The returned value is printed with %v;
// returning nil means "not implemented yet".
type Solver func(input string) any

// Main parses flags, reads the puzzle input, and runs the requested parts.
//
// Flags:
//
//	-p N    run only part N (1 or 2); default runs both
//	-i F    input file, relative to the day's directory ("-" for stdin); default input.txt
//	-q      print only the answer(s), one per line
func Main(part1, part2 Solver) {
	dir := callerDir()

	part := flag.Int("p", 0, "part to run (1 or 2); 0 runs both")
	inputPath := flag.String("i", "input.txt", "input file, relative to the day directory ('-' for stdin)")
	quiet := flag.Bool("q", false, "print only the answer(s)")
	flag.Parse()

	if *part < 0 || *part > 2 {
		fail(fmt.Errorf("invalid part %d", *part))
	}

	input, err := readInput(dir, *inputPath)
	if err != nil {
		fail(err)
	}

	for i, solve := range []Solver{part1, part2} {
		n := i + 1
		if *part != 0 && *part != n {
			continue
		}
		start := time.Now()
		answer := solve(input)
		elapsed := time.Since(start)

		switch {
		case *quiet && answer == nil:
			fail(fmt.Errorf("part %d is not implemented", n))
		case *quiet:
			fmt.Println(answer)
		case answer == nil:
			fmt.Printf("Part %d: (not implemented)\n", n)
		default:
			fmt.Printf("Part %d: %v\t(%s)\n", n, answer, elapsed.Round(time.Microsecond))
		}
	}
}

// Normalize strips a leading newline (convenient for raw-string examples) and
// trailing newlines/carriage returns. Leading spaces are preserved because they
// are sometimes significant.
func Normalize(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	s = strings.TrimPrefix(s, "\n")
	return strings.TrimRight(s, "\n")
}

func readInput(dir, path string) (string, error) {
	var data []byte
	var err error
	if path == "-" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		if !filepath.IsAbs(path) {
			path = filepath.Join(dir, path)
		}
		data, err = os.ReadFile(path)
		if errors.Is(err, os.ErrNotExist) {
			return "", fmt.Errorf("%s not found (try `aoc fetch`)", path)
		}
	}
	if err != nil {
		return "", err
	}
	return Normalize(string(data)), nil
}

// callerDir returns the source directory of the day package that called Main,
// so inputs resolve correctly no matter where `go run` is invoked from.
func callerDir() string {
	if _, file, _, ok := runtime.Caller(2); ok {
		dir := filepath.Dir(file)
		if st, err := os.Stat(dir); err == nil && st.IsDir() {
			return dir
		}
	}
	return "."
}

func fail(err error) {
	fmt.Fprintln(os.Stderr, "error:", err)
	os.Exit(1)
}
