export default function (eleventyConfig) {
	eleventyConfig.addPassthroughCopy("assets");
	eleventyConfig.addWatchTarget("lib/");

	return {
		templateFormats: ["njk"],
		dir: {
			input: ".",
			output: "dist",
		},
	};
}
