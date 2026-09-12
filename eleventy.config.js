import { RenderPlugin } from "@11ty/eleventy";
import { feedPlugin } from "@11ty/eleventy-plugin-rss";
import syntaxHighlight from "@11ty/eleventy-plugin-syntaxhighlight";
import site from "./_data/site.json" with { type: "json" };
import { formatIsoDate, formatPostDate } from "./lib/dates.js";

export default function (eleventyConfig) {
	eleventyConfig.addPassthroughCopy("assets");
	eleventyConfig.addPassthroughCopy({ "node_modules/prism-themes/themes/prism-coldark-cold.css": "assets/prism.css" });
	eleventyConfig.addPassthroughCopy({
		"resume/output/willian-pinheiro-resume.pdf": "willian-pinheiro-resume.pdf",
	});
	eleventyConfig.addWatchTarget("lib/");

	eleventyConfig.addCollection("posts", (api) => api.getFilteredByGlob("posts/*.md"));
	eleventyConfig.addFilter("postDate", formatPostDate);
	eleventyConfig.addFilter("isoDate", formatIsoDate);

	eleventyConfig.addPlugin(RenderPlugin);
	eleventyConfig.addPlugin(syntaxHighlight);
	eleventyConfig.addPlugin(feedPlugin, {
		type: "atom",
		outputPath: "/feed.xml",
		collection: { name: "posts", limit: 0 },
		metadata: {
			language: "en",
			title: site.title,
			subtitle: "Blog",
			base: site.url,
			author: { name: site.title },
		},
	});

	return {
		templateFormats: ["njk", "md"],
		dir: {
			input: ".",
			output: "dist",
		},
	};
}
