package main

import (
	"errors"
	"fmt"
	"html"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const defaultUserAgent = "aoc-go-cli (set AOC_USER_AGENT to add contact info)"

var client = &http.Client{Timeout: 30 * time.Second}

// request performs an authenticated request against adventofcode.com.
func request(root, method, target string, body io.Reader) ([]byte, error) {
	session, err := sessionToken(root)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(method, target, body)
	if err != nil {
		return nil, err
	}
	req.AddCookie(&http.Cookie{Name: "session", Value: session})
	ua := os.Getenv("AOC_USER_AGENT")
	if ua == "" {
		ua = defaultUserAgent
	}
	req.Header.Set("User-Agent", ua)
	if body != nil {
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%s %s: %s: %s", method, target, resp.Status, strings.TrimSpace(string(data)))
	}
	return data, nil
}

func sessionToken(root string) (string, error) {
	if s := strings.TrimSpace(os.Getenv("AOC_SESSION")); s != "" {
		return s, nil
	}
	data, err := os.ReadFile(filepath.Join(root, ".session"))
	if errors.Is(err, os.ErrNotExist) {
		return "", errors.New("no session: set AOC_SESSION or put your session cookie in .session")
	}
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

var (
	articleRE = regexp.MustCompile(`(?s)<article[^>]*>(.*?)</article>`)
	tagRE     = regexp.MustCompile(`<[^>]+>`)
	spaceRE   = regexp.MustCompile(`\s+`)
)

// articleText extracts the human-readable message from an answer response.
func articleText(page string) string {
	m := articleRE.FindStringSubmatch(page)
	if m == nil {
		return strings.TrimSpace(page)
	}
	text := html.UnescapeString(tagRE.ReplaceAllString(m[1], ""))
	return strings.TrimSpace(spaceRE.ReplaceAllString(text, " "))
}

func urlEncode(s string) string { return url.QueryEscape(s) }
