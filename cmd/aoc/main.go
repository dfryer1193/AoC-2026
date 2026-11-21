// Command aoc scaffolds, fetches, runs, tests, and submits Advent of Code days.
package main

import (
	"bufio"
	"bytes"
	"embed"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"text/template"
	"time"
)

const defaultYear = 2026

//go:embed templates/*.tmpl
var templates embed.FS

const usage = `usage: aoc <command> [day] [args]

Commands:
  new    [day]            create dayNN/ from the template and fetch its input
  fetch  [day]            download dayNN/input.txt (skipped if it already exists)
  run    [day] [flags]    run a day; flags are passed to the solution:
                            -p N   run only part N
                            -i F   input file relative to the day dir (- for stdin)
                            -q     print only the answer(s)
  test   [day|all] [args] run go test for a day (extra args go to go test)
  submit [day] -p N [-y]  run part N and submit its answer (-y skips the prompt)
  open   [day]            open the puzzle page in a browser

[day] defaults to today's puzzle during December (puzzles unlock at midnight EST).

Environment:
  AOC_SESSION     session cookie (otherwise read from .session in the repo root)
  AOC_YEAR        puzzle year (default 2026)
  AOC_USER_AGENT  User-Agent for requests; AoC asks automated tools to include contact info
`

func main() {
	if len(os.Args) < 2 {
		fmt.Fprint(os.Stderr, usage)
		os.Exit(2)
	}
	cmd, args := os.Args[1], os.Args[2:]

	var err error
	switch cmd {
	case "new":
		err = cmdNew(args)
	case "fetch":
		err = cmdFetch(args)
	case "run":
		err = cmdRun(args)
	case "test":
		err = cmdTest(args)
	case "submit":
		err = cmdSubmit(args)
	case "open":
		err = cmdOpen(args)
	case "help", "-h", "-help", "--help":
		fmt.Print(usage)
	default:
		err = fmt.Errorf("unknown command %q\n\n%s", cmd, usage)
	}

	var exitErr *exec.ExitError
	switch {
	case errors.As(err, &exitErr):
		os.Exit(exitErr.ExitCode()) // the child already reported its failure
	case err != nil:
		fmt.Fprintln(os.Stderr, "aoc:", err)
		os.Exit(1)
	}
}

// env describes the puzzle being worked on.
type env struct {
	root string // repository root
	year int
	day  int
}

func (e env) dirName() string { return fmt.Sprintf("day%02d", e.day) }
func (e env) dir() string     { return filepath.Join(e.root, e.dirName()) }
func (e env) pkg() string     { return "./" + e.dirName() }
func (e env) url() string     { return fmt.Sprintf("https://adventofcode.com/%d/day/%d", e.year, e.day) }

// setup resolves the repo root, year, and day. If args starts with a number it
// is taken as the day; the remaining args are returned.
func setup(args []string) (env, []string, error) {
	var e env
	var err error
	if e.root, err = findRoot(); err != nil {
		return e, nil, err
	}

	e.year = defaultYear
	if y := os.Getenv("AOC_YEAR"); y != "" {
		if e.year, err = strconv.Atoi(y); err != nil {
			return e, nil, fmt.Errorf("bad AOC_YEAR %q", y)
		}
	}

	if len(args) > 0 {
		if d, err := strconv.Atoi(args[0]); err == nil {
			if d < 1 || d > 25 {
				return e, nil, fmt.Errorf("day %d out of range", d)
			}
			e.day = d
			return e, args[1:], nil
		}
	}

	now := time.Now().In(time.FixedZone("EST", -5*60*60))
	if now.Year() != e.year || now.Month() != time.December || now.Day() > 25 {
		return e, nil, errors.New("no puzzle today; specify a day")
	}
	e.day = now.Day()
	return e, args, nil
}

// findRoot walks up from the working directory to the repository root.
func findRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			if _, err := os.Stat(filepath.Join(dir, "cmd", "aoc")); err == nil {
				return dir, nil
			}
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			return "", errors.New("not inside the advent-of-code repository")
		}
		dir = parent
	}
}

func cmdNew(args []string) error {
	e, _, err := setup(args)
	if err != nil {
		return err
	}

	if _, err := os.Stat(e.dir()); err == nil {
		fmt.Printf("%s already exists\n", e.dirName())
	} else {
		if err := scaffold(e); err != nil {
			return err
		}
		fmt.Printf("created %s/\n", e.dirName())
	}

	if err := fetchInput(e); err != nil {
		fmt.Fprintf(os.Stderr, "aoc: could not fetch input: %v\n", err)
	}
	fmt.Println(e.url())
	return nil
}

func scaffold(e env) error {
	module, err := modulePath(e.root)
	if err != nil {
		return err
	}
	data := map[string]any{"Module": module, "URL": e.url(), "Day": e.day, "Year": e.year}

	if err := os.MkdirAll(e.dir(), 0o755); err != nil {
		return err
	}
	for _, name := range []string{"main.go", "main_test.go"} {
		tmpl, err := template.ParseFS(templates, "templates/"+name+".tmpl")
		if err != nil {
			return err
		}
		var buf bytes.Buffer
		if err := tmpl.Execute(&buf, data); err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(e.dir(), name), buf.Bytes(), 0o644); err != nil {
			return err
		}
	}
	return nil
}

var moduleRE = regexp.MustCompile(`(?m)^module\s+(\S+)`)

func modulePath(root string) (string, error) {
	data, err := os.ReadFile(filepath.Join(root, "go.mod"))
	if err != nil {
		return "", err
	}
	m := moduleRE.FindSubmatch(data)
	if m == nil {
		return "", errors.New("no module line in go.mod")
	}
	return string(m[1]), nil
}

func cmdFetch(args []string) error {
	e, _, err := setup(args)
	if err != nil {
		return err
	}
	if _, err := os.Stat(e.dir()); err != nil {
		return fmt.Errorf("%s does not exist; run `aoc new %d`", e.dirName(), e.day)
	}
	return fetchInput(e)
}

func fetchInput(e env) error {
	path := filepath.Join(e.dir(), "input.txt")
	if _, err := os.Stat(path); err == nil {
		fmt.Printf("%s/input.txt already exists\n", e.dirName())
		return nil
	}
	body, err := request(e.root, "GET", e.url()+"/input", nil)
	if err != nil {
		return err
	}
	if err := os.WriteFile(path, body, 0o644); err != nil {
		return err
	}
	fmt.Printf("saved %s/input.txt\n", e.dirName())
	return nil
}

func cmdRun(args []string) error {
	e, rest, err := setup(args)
	if err != nil {
		return err
	}
	return goCmd(e.root, os.Stdout, append([]string{"run", e.pkg()}, rest...)...)
}

func cmdTest(args []string) error {
	if len(args) > 0 && args[0] == "all" {
		root, err := findRoot()
		if err != nil {
			return err
		}
		return goCmd(root, os.Stdout, append([]string{"test", "./..."}, args[1:]...)...)
	}
	e, rest, err := setup(args)
	if err != nil {
		return err
	}
	return goCmd(e.root, os.Stdout, append([]string{"test", e.pkg()}, rest...)...)
}

func cmdSubmit(args []string) error {
	e, rest, err := setup(args)
	if err != nil {
		return err
	}

	part, yes := 0, false
	for i := 0; i < len(rest); i++ {
		switch rest[i] {
		case "-y":
			yes = true
		case "-p":
			if i+1 < len(rest) {
				part, _ = strconv.Atoi(rest[i+1])
				i++
			}
		default:
			return fmt.Errorf("unexpected argument %q", rest[i])
		}
	}
	if part != 1 && part != 2 {
		return errors.New("submit needs -p 1 or -p 2")
	}

	var out bytes.Buffer
	if err := goCmd(e.root, &out, "run", e.pkg(), "-p", strconv.Itoa(part), "-q"); err != nil {
		return err
	}
	answer := lastLine(out.String())
	if answer == "" {
		return errors.New("solution printed no answer")
	}

	fmt.Printf("Day %d part %d answer: %s\n", e.day, part, answer)
	if !yes && !confirm("Submit?") {
		return nil
	}

	form := fmt.Sprintf("level=%d&answer=%s", part, urlEncode(answer))
	body, err := request(e.root, "POST", e.url()+"/answer", strings.NewReader(form))
	if err != nil {
		return err
	}
	msg := articleText(string(body))
	fmt.Println(msg)
	if !strings.Contains(msg, "That's the right answer") {
		os.Exit(1)
	}
	return nil
}

func cmdOpen(args []string) error {
	e, _, err := setup(args)
	if err != nil {
		return err
	}
	opener := "xdg-open"
	if _, err := exec.LookPath(opener); err != nil {
		opener = "open"
	}
	return exec.Command(opener, e.url()).Start()
}

// goCmd runs the go tool with stdout sent to out and stderr/stdin passed through.
func goCmd(dir string, out io.Writer, args ...string) error {
	cmd := exec.Command("go", args...)
	cmd.Dir = dir
	cmd.Stdin = os.Stdin
	cmd.Stdout = out
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func lastLine(s string) string {
	lines := strings.Split(strings.TrimSpace(s), "\n")
	return strings.TrimSpace(lines[len(lines)-1])
}

func confirm(prompt string) bool {
	fmt.Printf("%s [y/N] ", prompt)
	line, _ := bufio.NewReader(os.Stdin).ReadString('\n')
	line = strings.ToLower(strings.TrimSpace(line))
	return line == "y" || line == "yes"
}
