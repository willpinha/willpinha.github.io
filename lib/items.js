export function groupByYear(items) {
	const sorted = [...items].sort((a, b) => b.createdAt - a.createdAt);

	const groups = [];
	for (const item of sorted) {
		const year = item.createdAt.getUTCFullYear();
		if (groups.length === 0 || groups.at(-1).year !== year) {
			groups.push({ year, items: [] });
		}
		groups.at(-1).items.push(item);
	}
	return groups;
}
