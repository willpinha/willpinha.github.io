const graphqlEndpoint = "https://api.github.com/graphql";
const requestTimeout = 30_000;

// The GitHub search API never returns more than 1000 results per query
const searchResultCap = 1000;

const repoSearchQuery = `
query ($search: String!, $cursor: String) {
	search(type: REPOSITORY, query: $search, first: 100, after: $cursor) {
		pageInfo {
			hasNextPage
			endCursor
		}
		nodes {
			... on Repository {
				nameWithOwner
				url
				stargazerCount
				isPrivate
			}
		}
	}
}`;

const issueSearchQuery = `
query ($search: String!, $cursor: String) {
	search(type: ISSUE, first: 100, query: $search, after: $cursor) {
		pageInfo {
			hasNextPage
			endCursor
		}
		nodes {
			... on Issue {
				number
				title
				url
				createdAt
				repository {
					nameWithOwner
					url
					isPrivate
					owner {
						login
						avatarUrl(size: 64)
					}
				}
			}
			... on PullRequest {
				number
				title
				url
				createdAt
				state
				repository {
					nameWithOwner
					url
					isPrivate
					owner {
						login
						avatarUrl(size: 64)
					}
				}
			}
		}
	}
}`;

const discussionSearchQuery = `
query ($search: String!, $cursor: String) {
	search(type: DISCUSSION, first: 100, query: $search, after: $cursor) {
		pageInfo {
			hasNextPage
			endCursor
		}
		nodes {
			... on Discussion {
				number
				title
				url
				createdAt
				repository {
					nameWithOwner
					url
					isPrivate
					owner {
						login
						avatarUrl(size: 64)
					}
				}
			}
		}
	}
}`;

const contributionCalendarQuery = `
query ($login: String!) {
	user(login: $login) {
		contributionsCollection {
			contributionCalendar {
				totalContributions
				weeks {
					contributionDays {
						date
						contributionCount
						contributionLevel
					}
				}
			}
		}
	}
}`;

const contributionLevels = {
	FIRST_QUARTILE: 1,
	SECOND_QUARTILE: 2,
	THIRD_QUARTILE: 3,
	FOURTH_QUARTILE: 4,
};

export class GitHubClient {
	#token;
	#login;

	constructor(token, login) {
		this.#token = token;
		this.#login = login;
	}

	async famousRepos(minStars) {
		const search = `user:${this.#login} stars:>=${minStars} fork:false is:public sort:stars-desc`;
		const nodes = await this.#searchAll(repoSearchQuery, search);
		const prefix = `${this.#login}/`;
		return nodes
			.filter((n) => !n.isPrivate)
			.map((n) => ({
				name: n.nameWithOwner.startsWith(prefix)
					? n.nameWithOwner.slice(prefix.length)
					: n.nameWithOwner,
				url: n.url,
				stars: n.stargazerCount,
			}));
	}

	createdPullRequests() {
		const search = `is:pr is:public author:${this.#login} -user:${this.#login} sort:created-desc`;
		return this.#searchItems(issueSearchQuery, search);
	}

	participatedIssues() {
		const search = `is:issue is:public involves:${this.#login} sort:updated-desc`;
		return this.#searchItems(issueSearchQuery, search);
	}

	participatedDiscussions() {
		const search = `involves:${this.#login} sort:updated-desc`;
		return this.#searchItems(discussionSearchQuery, search);
	}

	async contributionCalendar() {
		const data = await this.#query(contributionCalendarQuery, { login: this.#login });
		const calendar = data.user.contributionsCollection.contributionCalendar;
		return {
			total: calendar.totalContributions,
			weeks: calendar.weeks.map((week) =>
				week.contributionDays.map((day) => ({
					date: parseDate(day.date),
					count: day.contributionCount,
					level: contributionLevels[day.contributionLevel] ?? 0,
				})),
			),
		};
	}

	async #searchItems(query, search) {
		const nodes = await this.#searchAll(query, search);
		return (
			nodes
				// Guards against private results when running with a broadly scoped local token
				.filter((n) => !n.repository.isPrivate)
				.map((n) => ({
					number: n.number,
					title: n.title,
					url: n.url,
					state: (n.state ?? "").toLowerCase(),
					createdAt: new Date(n.createdAt),
					repoName: n.repository.nameWithOwner,
					repoUrl: n.repository.url,
					repoOwnerLogin: n.repository.owner.login,
					repoOwnerAvatarUrl: n.repository.owner.avatarUrl,
				}))
		);
	}

	async #searchAll(query, search) {
		const all = [];
		let cursor = null;
		for (;;) {
			const { search: result } = await this.#query(query, { search, cursor });
			all.push(...result.nodes);
			if (!result.pageInfo.hasNextPage || all.length >= searchResultCap) {
				return all;
			}
			cursor = result.pageInfo.endCursor;
		}
	}

	async #query(query, variables) {
		const response = await fetch(graphqlEndpoint, {
			method: "POST",
			headers: {
				Authorization: `Bearer ${this.#token}`,
				"Content-Type": "application/json",
			},
			body: JSON.stringify({ query, variables }),
			signal: AbortSignal.timeout(requestTimeout),
		});
		if (!response.ok) {
			throw new Error(`graphql request returned status ${response.status} ${response.statusText}`);
		}
		const { data, errors } = await response.json();
		if (errors?.length) {
			throw new Error(`graphql error: ${errors[0].message}`);
		}
		return data;
	}
}

function parseDate(value) {
	const date = new Date(value);
	if (Number.isNaN(date.getTime())) {
		throw new Error(`parse contribution date "${value}"`);
	}
	return date;
}
