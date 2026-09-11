import Eleventy from "@11ty/eleventy";
import { afterAll, beforeAll, expect, test, vi } from "vitest";

const pageUrls = ["/", "/contributions/", "/pull-requests/", "/issues/", "/discussions/"];

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
	expect([...pagesByUrl.keys()].sort()).toEqual([...pageUrls].sort());
});

test.each(pageUrls)("renders %s", (url) => {
	expect(pagesByUrl.get(url)).toMatchSnapshot();
});
