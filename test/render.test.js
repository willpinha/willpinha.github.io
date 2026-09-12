import Eleventy from "@11ty/eleventy";
import { existsSync, readdirSync } from "node:fs";
import { afterAll, beforeAll, expect, test, vi } from "vitest";

const staticUrls = ["/", "/blog/", "/contributions/", "/photos/", "/pull-requests/", "/issues/", "/discussions/"];

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
	const expected = [...staticUrls, "/feed.xml", ...postSlugs.map((slug) => `/blog/${slug}/`)];
	expect([...pagesByUrl.keys()].sort()).toEqual(expected.sort());
});

test.each(postSlugs)("renders post %s with title and date", (slug) => {
	const html = pagesByUrl.get(`/blog/${slug}/`);
	expect(html).toMatch(/<h1>.+<\/h1>/);
	expect(html).toMatch(/<time datetime="\d{4}-\d{2}-\d{2}">\d{2} [A-Z][a-z]{2}, \d{4}<\/time>/);
});

test("every photo linked on the photos page exists in assets", () => {
	const html = pagesByUrl.get("/photos/");
	const hrefs = [...html.matchAll(/href="(\/assets\/photos\/[^"]+)"/g)].map((match) => match[1]);
	expect(hrefs.length).toBeGreaterThan(0);
	for (const href of hrefs) {
		expect(existsSync(href.slice(1)), `${href} is missing`).toBe(true);
	}
});

test("renders each photo item with a date and a count badge only when it has several photos", () => {
	const items = pagesByUrl.get("/photos/").split('<li class="photo-item">').slice(1);
	expect(items.length).toBeGreaterThan(0);
	for (const item of items) {
		expect(item).toMatch(/<time datetime="\d{4}-\d{2}-\d{2}">\d{2} [A-Z][a-z]{2}, \d{4}<\/time>/);
		const photos = item.match(/href="\/assets\/photos\//g).length;
		if (photos > 1) {
			expect(item).toContain(`<span class="photo-count">${photos}</span>`);
		} else {
			expect(item).not.toContain('class="photo-count"');
		}
	}
});

test("lists every post in the feed", () => {
	const feed = pagesByUrl.get("/feed.xml");
	for (const slug of postSlugs) {
		expect(feed).toContain(`https://willpinha.github.io/blog/${slug}/`);
	}
});
