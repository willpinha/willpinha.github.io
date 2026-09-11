import { dayMs, dayOfYear } from "./dates.js";

export function mockContributionCalendar(now) {
	const end = Date.UTC(now.getUTCFullYear(), now.getUTCMonth(), now.getUTCDate());
	const endDate = new Date(end);
	const start = Date.UTC(
		endDate.getUTCFullYear() - 1,
		endDate.getUTCMonth(),
		endDate.getUTCDate() - endDate.getUTCDay(),
	);

	const weeks = [];
	let week = [];
	let total = 0;
	for (let t = start; t <= end; t += dayMs) {
		const date = new Date(t);
		const count = (dayOfYear(date) * 7 + date.getUTCDay() * 3) % 12;
		total += count;
		week.push({ date, count, level: mockContributionLevel(count) });
		if (date.getUTCDay() === 6) {
			weeks.push(week);
			week = [];
		}
	}
	if (week.length > 0) {
		weeks.push(week);
	}

	return { total, weeks };
}

function mockContributionLevel(count) {
	if (count === 0) return 0;
	if (count <= 3) return 1;
	if (count <= 6) return 2;
	if (count <= 9) return 3;
	return 4;
}

export const mockRepos = [
	{ name: "daisy-components", url: "https://github.com/willpinha/daisy-components", stars: 437 },
	{ name: "mantine-themes", url: "https://github.com/willpinha/mantine-themes", stars: 27 },
];

// Sorted by creation date descending, as returned by the API.
// Six entries so the section limit of five is exercised.
export const mockPullRequests = [
	{
		number: 1234,
		title: "Add [config] parsing",
		url: "https://github.com/linux/hello/pull/1234",
		state: "merged",
		createdAt: new Date("2026-03-10T12:00:00Z"),
		repoName: "linux/hello",
		repoUrl: "https://github.com/linux/hello",
		repoOwnerLogin: "torvalds",
		repoOwnerAvatarUrl: "https://avatars.githubusercontent.com/u/1024025?s=64&v=4",
	},
	{
		number: 124,
		title: "Fix typo in docs",
		url: "https://github.com/linux/world/pull/124",
		state: "closed",
		createdAt: new Date("2026-01-05T08:30:00Z"),
		repoName: "linux/world",
		repoUrl: "https://github.com/linux/world",
		repoOwnerLogin: "torvalds",
		repoOwnerAvatarUrl: "https://avatars.githubusercontent.com/u/1024025?s=64&v=4",
	},
	{
		number: 14,
		title: "Support dark mode",
		url: "https://github.com/go/wails/pull/14",
		state: "open",
		createdAt: new Date("2025-11-20T22:15:00Z"),
		repoName: "go/wails",
		repoUrl: "https://github.com/go/wails",
		repoOwnerLogin: "golang",
		repoOwnerAvatarUrl: "https://avatars.githubusercontent.com/u/4314092?s=64&v=4",
	},
	{
		number: 99,
		title: "Improve *error* messages",
		url: "https://github.com/go/wails/pull/99",
		state: "merged",
		createdAt: new Date("2025-09-01T10:00:00Z"),
		repoName: "go/wails",
		repoUrl: "https://github.com/go/wails",
		repoOwnerLogin: "golang",
		repoOwnerAvatarUrl: "https://avatars.githubusercontent.com/u/4314092?s=64&v=4",
	},
	{
		number: 50,
		title: "Remove deprecated flag",
		url: "https://github.com/linux/hello/pull/50",
		state: "closed",
		createdAt: new Date("2025-05-12T14:45:00Z"),
		repoName: "linux/hello",
		repoUrl: "https://github.com/linux/hello",
		repoOwnerLogin: "torvalds",
		repoOwnerAvatarUrl: "https://avatars.githubusercontent.com/u/1024025?s=64&v=4",
	},
	{
		number: 7,
		title: "Initial CI setup",
		url: "https://github.com/linux/world/pull/7",
		state: "merged",
		createdAt: new Date("2025-02-03T09:00:00Z"),
		repoName: "linux/world",
		repoUrl: "https://github.com/linux/world",
		repoOwnerLogin: "torvalds",
		repoOwnerAvatarUrl: "https://avatars.githubusercontent.com/u/1024025?s=64&v=4",
	},
];

// Sorted by update date descending, so creation dates arrive out of order
export const mockIssues = [
	{
		number: 12,
		title: "Some issue",
		url: "https://github.com/willpinha/foo/issues/12",
		state: "",
		createdAt: new Date("2024-06-01T11:00:00Z"),
		repoName: "willpinha/foo",
		repoUrl: "https://github.com/willpinha/foo",
		repoOwnerLogin: "willpinha",
		repoOwnerAvatarUrl: "https://avatars.githubusercontent.com/u/86596621?s=64&v=4",
	},
	{
		number: 80,
		title: "Crash on <startup>",
		url: "https://github.com/linux/hello/issues/80",
		state: "",
		createdAt: new Date("2026-02-14T16:20:00Z"),
		repoName: "linux/hello",
		repoUrl: "https://github.com/linux/hello",
		repoOwnerLogin: "torvalds",
		repoOwnerAvatarUrl: "https://avatars.githubusercontent.com/u/1024025?s=64&v=4",
	},
	{
		number: 33,
		title: "Docs are outdated",
		url: "https://github.com/willpinha/foo/issues/33",
		state: "",
		createdAt: new Date("2025-08-30T07:10:00Z"),
		repoName: "willpinha/foo",
		repoUrl: "https://github.com/willpinha/foo",
		repoOwnerLogin: "willpinha",
		repoOwnerAvatarUrl: "https://avatars.githubusercontent.com/u/86596621?s=64&v=4",
	},
];

export const mockDiscussions = [
	{
		number: 4,
		title: "Some question",
		url: "https://github.com/willpinha/bar/discussions/4",
		state: "",
		createdAt: new Date("2026-04-22T19:05:00Z"),
		repoName: "willpinha/bar",
		repoUrl: "https://github.com/willpinha/bar",
		repoOwnerLogin: "willpinha",
		repoOwnerAvatarUrl: "https://avatars.githubusercontent.com/u/86596621?s=64&v=4",
	},
];
