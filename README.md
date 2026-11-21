# Advent of Code 2026

Solutions in Go. Each day is a self-contained `main` package in `dayNN/`.

## Setup

```sh
go install ./cmd/aoc        # puts `aoc` on your PATH (via $GOBIN / ~/go/bin)
```

### Session cookie

Fetching inputs and submitting answers need your adventofcode.com `session`
cookie. In a Chromium-based browser (Chrome, Chromium, Brave, Edge, Vivaldi):

1. Log in at <https://adventofcode.com>.
2. Open DevTools (`F12` or `Ctrl+Shift+I`) and go to the **Application** tab.
3. In the sidebar, expand **Storage → Cookies** and select `https://adventofcode.com`.
4. Click the `session` row and copy its **Cookie Value** from the pane below.

Save it to `.session` at the repo root (gitignored) without it landing in your
shell history:

```sh
read -s 'token?Session cookie: ' && print -r -- "$token" > .session && chmod 600 .session && unset token
```

Or export it as `AOC_SESSION` instead, which takes priority over `.session`.
The cookie lasts a long time, but if `aoc fetch` starts failing, grab a fresh one.

AoC asks automated tools to identify themselves, so set `AOC_USER_AGENT` to
something with your contact info, e.g. a repo URL.

## Daily workflow

```sh
aoc new            # scaffold today's dayNN/ and download input.txt
$EDITOR dayNN/main.go dayNN/main_test.go
aoc test           # check against the example pasted into main_test.go
aoc run            # run both parts on the real input
aoc submit -p 1    # run part 1 and submit the answer
```

The day defaults to today's puzzle in December; pass a number to pick one:
`aoc run 3 -p 2`, `aoc test 3 -v`, `aoc new 7`. Everything works with plain Go too:

```sh
go run ./day03 -p 1 -i example.txt   # paths are relative to the day directory
go test ./day03
```

| Command | What it does |
| --- | --- |
| `aoc new [day]` | Create `dayNN/` from the template and fetch the input |
| `aoc fetch [day]` | Download `dayNN/input.txt` if it's missing |
| `aoc run [day] [-p N] [-i file] [-q]` | Run a day |
| `aoc test [day\|all] [go test args]` | Run tests |
| `aoc submit [day] -p N [-y]` | Run part N and submit it (`-y` skips the prompt) |
| `aoc open [day]` | Open the puzzle in a browser |

## Writing a solution

```go
func part1(input string) any { ... }   // return nil until implemented
```

Input arrives with trailing newlines trimmed. In tests, paste the example into
the raw string and fill in the expected answers; a `nil` answer skips the test:

```go
aoctest.Check(t, part1, example, 142)
```

Helpers live in `internal/`:

- `parse` – `Lines`, `Blocks`, `Ints` (pulls every integer out of a string), `IntLines`, `Atoi`, `Digits`
- `grid` – `Point` with `Add`/`RotateCW`/`Manhattan`, `Dirs4`/`Dirs8`, and a byte `Grid` with `Parse`, `At`, `Find`, `Neighbors`, `All`
- `xmath` – `Abs`, `Sign`, `Mod` (always non-negative), `GCD`, `LCM`, `Sum`
