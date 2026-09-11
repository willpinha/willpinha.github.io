import { formatDay, monthName } from "./dates.js";

const cellSize = 9;
const cellGap = 3;
const leftPad = 24;
const topPad = 16;

const labeledWeekdays = [
	{ weekday: 1, label: "Mon" },
	{ weekday: 3, label: "Wed" },
	{ weekday: 5, label: "Fri" },
];

export function buildContributionGraph(calendar) {
	const step = cellSize + cellGap;

	const cells = [];
	const monthLabels = [];
	let lastMonth = -1;

	calendar.weeks.forEach((week, weekIndex) => {
		const x = leftPad + weekIndex * step;
		if (week.length > 0) {
			const month = week[0].date.getUTCMonth();
			if (month !== lastMonth) {
				monthLabels.push({ x, label: monthName(week[0].date) });
				lastMonth = month;
			}
		}
		for (const day of week) {
			cells.push({
				x,
				y: topPad + day.date.getUTCDay() * step,
				level: day.level,
				tooltip: contributionTooltip(day),
			});
		}
	});

	const weekdayLabels = labeledWeekdays.map(({ weekday, label }) => ({
		y: topPad + weekday * step + cellSize,
		label,
	}));

	return {
		cells,
		monthLabels,
		weekdayLabels,
		total: calendar.total,
		width: leftPad + calendar.weeks.length * step,
		height: topPad + 7 * step,
	};
}

export function buildContributors(items, ownLogin) {
	const byLogin = new Map();
	for (const item of items) {
		if (item.repoOwnerLogin === ownLogin) continue;
		let contributor = byLogin.get(item.repoOwnerLogin);
		if (!contributor) {
			contributor = { login: item.repoOwnerLogin, avatarUrl: item.repoOwnerAvatarUrl, count: 0 };
			byLogin.set(item.repoOwnerLogin, contributor);
		}
		contributor.count++;
	}

	return [...byLogin.values()]
		.map((c) => ({ ...c, tooltip: contributorTooltip(c) }))
		.sort((a, b) => b.count - a.count || compareStrings(a.login, b.login));
}

function compareStrings(a, b) {
	if (a < b) return -1;
	if (a > b) return 1;
	return 0;
}

function contributorTooltip(contributor) {
	const count = contributor.count === 1 ? "1 contribution" : `${contributor.count} contributions`;
	return `${contributor.login} (${count})`;
}

function contributionTooltip(day) {
	let count = "No contributions";
	if (day.count === 1) {
		count = "1 contribution";
	} else if (day.count > 1) {
		count = `${day.count} contributions`;
	}
	return `${count} on ${formatDay(day.date)}`;
}
