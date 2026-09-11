import site from "./site.json" with { type: "json" };
import { buildContributionGraph, buildContributors } from "../lib/contributions.js";
import { formatTimestamp } from "../lib/dates.js";
import { GitHubClient } from "../lib/github.js";
import { groupByYear } from "../lib/items.js";
import {
	mockContributionCalendar,
	mockDiscussions,
	mockIssues,
	mockPullRequests,
	mockRepos,
} from "../lib/mock.js";

const minRepoStars = 10;

// Caps how many items each "Latest" section of the contributions page shows
const sectionItemLimit = 5;

export default async function () {
	const now = new Date();
	const { repos, pullRequests, issues, discussions, calendar } =
		process.env.ELEVENTY_RUN_MODE === "build" ? await fetchData() : mockData(now);

	return {
		updatedAt: formatTimestamp(now),
		repos,
		pullRequests: section(pullRequests),
		issues: section(issues),
		discussions: section(discussions),
		contributions: {
			...buildContributionGraph(calendar),
			contributors: buildContributors([...pullRequests, ...issues, ...discussions], site.login),
		},
	};
}

async function fetchData() {
	const token = process.env.GITHUB_TOKEN;
	if (!token) {
		throw new Error("GITHUB_TOKEN environment variable is not set");
	}
	const client = new GitHubClient(token, site.login);
	return {
		repos: await client.famousRepos(minRepoStars),
		pullRequests: await client.createdPullRequests(),
		issues: await client.participatedIssues(),
		discussions: await client.participatedDiscussions(),
		calendar: await client.contributionCalendar(),
	};
}

function mockData(now) {
	return {
		repos: mockRepos,
		pullRequests: mockPullRequests,
		issues: mockIssues,
		discussions: mockDiscussions,
		calendar: mockContributionCalendar(now),
	};
}

function section(items) {
	return {
		latest: items.slice(0, sectionItemLimit),
		byYear: groupByYear(items),
	};
}
