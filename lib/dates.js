const monthNames = [
	"Jan",
	"Feb",
	"Mar",
	"Apr",
	"May",
	"Jun",
	"Jul",
	"Aug",
	"Sep",
	"Oct",
	"Nov",
	"Dec",
];

export const dayMs = 24 * 60 * 60 * 1000;

function pad(n) {
	return String(n).padStart(2, "0");
}

// "2006/01/02 15:04" in UTC
export function formatTimestamp(date) {
	const day = `${date.getUTCFullYear()}/${pad(date.getUTCMonth() + 1)}/${pad(date.getUTCDate())}`;
	const time = `${pad(date.getUTCHours())}:${pad(date.getUTCMinutes())}`;
	return `${day} ${time}`;
}

// "Jan 2, 2006" in UTC
export function formatDay(date) {
	return `${monthName(date)} ${date.getUTCDate()}, ${date.getUTCFullYear()}`;
}

// "02 Jan, 2006" in UTC
export function formatPostDate(date) {
	return `${pad(date.getUTCDate())} ${monthName(date)}, ${date.getUTCFullYear()}`;
}

// "2006-01-02" in UTC
export function formatIsoDate(date) {
	return `${date.getUTCFullYear()}-${pad(date.getUTCMonth() + 1)}-${pad(date.getUTCDate())}`;
}

export function monthName(date) {
	return monthNames[date.getUTCMonth()];
}

// 1-based, matching Go's time.YearDay
export function dayOfYear(date) {
	const startOfYear = Date.UTC(date.getUTCFullYear(), 0, 0);
	return Math.round((date.getTime() - startOfYear) / dayMs);
}
