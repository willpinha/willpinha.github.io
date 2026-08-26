package main

import "time"

func mockContributionCalendar(now time.Time) contributionCalendar {
	end := now.UTC().Truncate(24 * time.Hour)
	start := end.AddDate(-1, 0, -int(end.Weekday()))

	var weeks [][]contributionDay
	var week []contributionDay
	total := 0
	for d := start; !d.After(end); d = d.AddDate(0, 0, 1) {
		count := (d.YearDay()*7 + int(d.Weekday())*3) % 12
		total += count
		week = append(week, contributionDay{Date: d, Count: count, Level: mockContributionLevel(count)})
		if d.Weekday() == time.Saturday {
			weeks = append(weeks, week)
			week = nil
		}
	}
	if len(week) > 0 {
		weeks = append(weeks, week)
	}

	return contributionCalendar{Total: total, Weeks: weeks}
}

func mockContributionLevel(count int) int {
	switch {
	case count == 0:
		return 0
	case count <= 3:
		return 1
	case count <= 6:
		return 2
	case count <= 9:
		return 3
	default:
		return 4
	}
}

var mockRepos = []repo{
	{Name: "daisy-components", URL: "https://github.com/willpinha/daisy-components", Stars: 437},
	{Name: "mantine-themes", URL: "https://github.com/willpinha/mantine-themes", Stars: 27},
}

// Sorted by creation date descending, as returned by the API.
// Six entries so the home section limit of five is exercised.
var mockPullRequests = []item{
	{
		Number:             1234,
		Title:              "Add [config] parsing",
		URL:                "https://github.com/linux/hello/pull/1234",
		State:              "merged",
		CreatedAt:          time.Date(2026, 3, 10, 12, 0, 0, 0, time.UTC),
		RepoName:           "linux/hello",
		RepoURL:            "https://github.com/linux/hello",
		RepoOwnerLogin:     "torvalds",
		RepoOwnerAvatarURL: "https://avatars.githubusercontent.com/u/1024025?s=64&v=4",
	},
	{
		Number:             124,
		Title:              "Fix typo in docs",
		URL:                "https://github.com/linux/world/pull/124",
		State:              "closed",
		CreatedAt:          time.Date(2026, 1, 5, 8, 30, 0, 0, time.UTC),
		RepoName:           "linux/world",
		RepoURL:            "https://github.com/linux/world",
		RepoOwnerLogin:     "torvalds",
		RepoOwnerAvatarURL: "https://avatars.githubusercontent.com/u/1024025?s=64&v=4",
	},
	{
		Number:             14,
		Title:              "Support dark mode",
		URL:                "https://github.com/go/wails/pull/14",
		State:              "open",
		CreatedAt:          time.Date(2025, 11, 20, 22, 15, 0, 0, time.UTC),
		RepoName:           "go/wails",
		RepoURL:            "https://github.com/go/wails",
		RepoOwnerLogin:     "golang",
		RepoOwnerAvatarURL: "https://avatars.githubusercontent.com/u/4314092?s=64&v=4",
	},
	{
		Number:             99,
		Title:              "Improve *error* messages",
		URL:                "https://github.com/go/wails/pull/99",
		State:              "merged",
		CreatedAt:          time.Date(2025, 9, 1, 10, 0, 0, 0, time.UTC),
		RepoName:           "go/wails",
		RepoURL:            "https://github.com/go/wails",
		RepoOwnerLogin:     "golang",
		RepoOwnerAvatarURL: "https://avatars.githubusercontent.com/u/4314092?s=64&v=4",
	},
	{
		Number:             50,
		Title:              "Remove deprecated flag",
		URL:                "https://github.com/linux/hello/pull/50",
		State:              "closed",
		CreatedAt:          time.Date(2025, 5, 12, 14, 45, 0, 0, time.UTC),
		RepoName:           "linux/hello",
		RepoURL:            "https://github.com/linux/hello",
		RepoOwnerLogin:     "torvalds",
		RepoOwnerAvatarURL: "https://avatars.githubusercontent.com/u/1024025?s=64&v=4",
	},
	{
		Number:             7,
		Title:              "Initial CI setup",
		URL:                "https://github.com/linux/world/pull/7",
		State:              "merged",
		CreatedAt:          time.Date(2025, 2, 3, 9, 0, 0, 0, time.UTC),
		RepoName:           "linux/world",
		RepoURL:            "https://github.com/linux/world",
		RepoOwnerLogin:     "torvalds",
		RepoOwnerAvatarURL: "https://avatars.githubusercontent.com/u/1024025?s=64&v=4",
	},
}

// Sorted by update date descending, so creation dates arrive out of order.
var mockIssues = []item{
	{
		Number:             12,
		Title:              "Some issue",
		URL:                "https://github.com/willpinha/foo/issues/12",
		CreatedAt:          time.Date(2024, 6, 1, 11, 0, 0, 0, time.UTC),
		RepoName:           "willpinha/foo",
		RepoURL:            "https://github.com/willpinha/foo",
		RepoOwnerLogin:     "willpinha",
		RepoOwnerAvatarURL: "https://avatars.githubusercontent.com/u/86596621?s=64&v=4",
	},
	{
		Number:             80,
		Title:              "Crash on <startup>",
		URL:                "https://github.com/linux/hello/issues/80",
		CreatedAt:          time.Date(2026, 2, 14, 16, 20, 0, 0, time.UTC),
		RepoName:           "linux/hello",
		RepoURL:            "https://github.com/linux/hello",
		RepoOwnerLogin:     "torvalds",
		RepoOwnerAvatarURL: "https://avatars.githubusercontent.com/u/1024025?s=64&v=4",
	},
	{
		Number:             33,
		Title:              "Docs are outdated",
		URL:                "https://github.com/willpinha/foo/issues/33",
		CreatedAt:          time.Date(2025, 8, 30, 7, 10, 0, 0, time.UTC),
		RepoName:           "willpinha/foo",
		RepoURL:            "https://github.com/willpinha/foo",
		RepoOwnerLogin:     "willpinha",
		RepoOwnerAvatarURL: "https://avatars.githubusercontent.com/u/86596621?s=64&v=4",
	},
}

var mockDiscussions = []item{
	{
		Number:             4,
		Title:              "Some question",
		URL:                "https://github.com/willpinha/bar/discussions/4",
		CreatedAt:          time.Date(2026, 4, 22, 19, 5, 0, 0, time.UTC),
		RepoName:           "willpinha/bar",
		RepoURL:            "https://github.com/willpinha/bar",
		RepoOwnerLogin:     "willpinha",
		RepoOwnerAvatarURL: "https://avatars.githubusercontent.com/u/86596621?s=64&v=4",
	},
}
