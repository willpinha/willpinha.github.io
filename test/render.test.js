import Eleventy from "@11ty/eleventy";
import { readdirSync } from "node:fs";
import { afterAll, beforeAll, expect, test, vi } from "vitest";

const snapshotUrls = ["/", "/blog/", "/contributions/", "/pull-requests/", "/issues/", "/discussions/"];

const postSlugs = readdirSync("posts")
	.filter((file) => file.endsWith(".md"))
	.map((file) => file.slice(0, -".md".length));

let pagesByUrl;

beforeAll(async () => {
	vi.useFakeTimers({ toFake: ["Date"], now: new Date("2026-07-13T13:27:00Z") });

	// Serve mode renders with mock data, as the dev server does
	const eleventy = new Eleventy(".", "dist", { quietMode: true, runMode: "serve" });
	const pages = await eleventy.toJSON();
	pagesByUrl = new Map(pages.map((page) => [page.url, page.content]));
});

afterAll(() => {
	vi.useRealTimers();
});

test("renders exactly the expected pages", () => {
	const expected = [...snapshotUrls, "/feed.xml", ...postSlugs.map((slug) => `/blog/${slug}/`)];
	expect([...pagesByUrl.keys()].sort()).toEqual(expected.sort());
});

test.each(snapshotUrls)("renders %s", (url) => {
	expect(pagesByUrl.get(url)).toMatchSnapshot();
});

test.each(postSlugs)("renders post %s with title and date", (slug) => {
	const html = pagesByUrl.get(`/blog/${slug}/`);
	expect(html).toMatch(/<h1>.+<\/h1>/);
	expect(html).toMatch(/<time datetime="\d{4}-\d{2}-\d{2}">\d{2} [A-Z][a-z]{2}, \d{4}<\/time>/);
});

test("lists every post in the feed", () => {
	const feed = pagesByUrl.get("/feed.xml");
	for (const slug of postSlugs) {
		expect(feed).toContain(`https://willpinha.github.io/blog/${slug}/`);
	}
});
