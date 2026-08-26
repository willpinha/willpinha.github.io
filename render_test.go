package main

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"
	"time"
)

var update = flag.Bool("update", false, "rewrite golden files with the current output")

func TestRenderHome(t *testing.T) {
	now := time.Date(2026, 7, 13, 13, 27, 0, 0, time.UTC)
	var buf bytes.Buffer
	if err := renderHome(&buf, now, mockRepos, mockPullRequests, mockIssues, mockDiscussions); err != nil {
		t.Fatal(err)
	}
	assertGolden(t, "home.golden.html", buf.String())
}

func TestRenderSeeAll(t *testing.T) {
	now := time.Date(2026, 7, 13, 13, 27, 0, 0, time.UTC)
	tests := []struct {
		name      string
		pageTitle string
		items     []item
		golden    string
	}{
		{
			name:      "pull requests with state grouped by year",
			pageTitle: "Pull requests I created",
			items:     mockPullRequests,
			golden:    "pull-requests.golden.html",
		},
		{
			name:      "issues without state and with unsorted input",
			pageTitle: "Issues I participated in",
			items:     mockIssues,
			golden:    "issues.golden.html",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			if err := renderSeeAll(&buf, now, tt.pageTitle, tt.items); err != nil {
				t.Fatal(err)
			}
			assertGolden(t, tt.golden, buf.String())
		})
	}
}

func assertGolden(t *testing.T, name, got string) {
	t.Helper()
	path := filepath.Join("testdata", name)
	if *update {
		if err := os.WriteFile(path, []byte(got), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	expected, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if got != string(expected) {
		t.Errorf("output differs from %s\ngot:\n%s\nexpected:\n%s", path, got, expected)
	}
}
