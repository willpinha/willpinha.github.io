package main

import "time"

var mockRepos = []repo{
	{Name: "willpinha/daisy-components", URL: "https://github.com/willpinha/daisy-components", Stars: 437},
	{Name: "willpinha/mantine-themes", URL: "https://github.com/willpinha/mantine-themes", Stars: 27},
}

// Sorted by creation date descending, as returned by the API.
// Six entries so the home section limit of five is exercised.
var mockPullRequests = []item{
	{
		Number:    1234,
		Title:     "Add [config] parsing",
		URL:       "https://github.com/linux/hello/pull/1234",
		State:     "merged",
		CreatedAt: time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
		RepoName:  "linux/hello",
		RepoURL:   "https://github.com/linux/hello",
	},
	{
		Number:    124,
		Title:     "Fix typo in docs",
		URL:       "https://github.com/linux/world/pull/124",
		State:     "closed",
		CreatedAt: time.Date(2026, 1, 5, 8, 30, 0, 0, time.UTC),
		RepoName:  "linux/world",
		RepoURL:   "https://github.com/linux/world",
	},
	{
		Number:    14,
		Title:     "Support dark mode",
		URL:       "https://github.com/go/wails/pull/14",
		State:     "open",
		CreatedAt: time.Date(2025, 11, 20, 22, 15, 0, 0, time.UTC),
		RepoName:  "go/wails",
		RepoURL:   "https://github.com/go/wails",
	},
	{
		Number:    99,
		Title:     "Improve *error* messages",
		URL:       "https://github.com/go/wails/pull/99",
		State:     "merged",
		CreatedAt: time.Date(2025, 9, 1, 10, 0, 0, 0, time.UTC),
		RepoName:  "go/wails",
		RepoURL:   "https://github.com/go/wails",
	},
	{
		Number:    50,
		Title:     "Remove deprecated flag",
		URL:       "https://github.com/linux/hello/pull/50",
		State:     "closed",
		CreatedAt: time.Date(2025, 5, 12, 14, 45, 0, 0, time.UTC),
		RepoName:  "linux/hello",
		RepoURL:   "https://github.com/linux/hello",
	},
	{
		Number:    7,
		Title:     "Initial CI setup",
		URL:       "https://github.com/linux/world/pull/7",
		State:     "merged",
		CreatedAt: time.Date(2025, 2, 3, 9, 0, 0, 0, time.UTC),
		RepoName:  "linux/world",
		RepoURL:   "https://github.com/linux/world",
	},
}

// Sorted by update date descending, so creation dates arrive out of order.
var mockIssues = []item{
	{
		Number:    12,
		Title:     "Some issue",
		URL:       "https://github.com/willpinha/foo/issues/12",
		CreatedAt: time.Date(2024, 6, 1, 11, 0, 0, 0, time.UTC),
		RepoName:  "willpinha/foo",
		RepoURL:   "https://github.com/willpinha/foo",
	},
	{
		Number:    80,
		Title:     "Crash on <startup>",
		URL:       "https://github.com/linux/hello/issues/80",
		CreatedAt: time.Date(2026, 2, 14, 16, 20, 0, 0, time.UTC),
		RepoName:  "linux/hello",
		RepoURL:   "https://github.com/linux/hello",
	},
	{
		Number:    33,
		Title:     "Docs are outdated",
		URL:       "https://github.com/willpinha/foo/issues/33",
		CreatedAt: time.Date(2025, 8, 30, 7, 10, 0, 0, time.UTC),
		RepoName:  "willpinha/foo",
		RepoURL:   "https://github.com/willpinha/foo",
	},
}

var mockDiscussions = []item{
	{
		Number:    4,
		Title:     "Some question",
		URL:       "https://github.com/willpinha/bar/discussions/4",
		CreatedAt: time.Date(2026, 4, 22, 19, 5, 0, 0, time.UTC),
		RepoName:  "willpinha/bar",
		RepoURL:   "https://github.com/willpinha/bar",
	},
}
