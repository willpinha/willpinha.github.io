package main

import (
	"embed"
	"fmt"
	"html/template"
	"io"
	"slices"
	"time"
)

//go:embed templates/*.html
var templatesFS embed.FS

// sectionItemLimit caps how many items each "Latest" home section shows.
const sectionItemLimit = 5

const timeFormat = "2006/01/02 15:04"

type repo struct {
	Name  string
	URL   string
	Stars int
}

type item struct {
	Number             int
	Title              string
	URL                string
	State              string
	CreatedAt          time.Time
	RepoName           string
	RepoURL            string
	RepoOwnerLogin     string
	RepoOwnerAvatarURL string
}

type yearGroup struct {
	Year  int
	Items []item
}

type contributionDay struct {
	Date  time.Time
	Count int
	Level int
}

type contributionCalendar struct {
	Total int
	Weeks [][]contributionDay
}

const (
	contribCellSize = 9
	contribCellGap  = 3
	contribLeftPad  = 24
	contribTopPad   = 16
)

type contributionCell struct {
	X       int
	Y       int
	Level   int
	Tooltip string
}

type contributionMonthLabel struct {
	X     int
	Label string
}

type contributionWeekdayLabel struct {
	Y     int
	Label string
}

type contributionGraph struct {
	Cells         []contributionCell
	MonthLabels   []contributionMonthLabel
	WeekdayLabels []contributionWeekdayLabel
	Total         int
	Width         int
	Height        int
}

func buildContributionGraph(cal contributionCalendar) contributionGraph {
	step := contribCellSize + contribCellGap

	var cells []contributionCell
	var monthLabels []contributionMonthLabel
	lastMonth := time.Month(0)

	for weekIndex, week := range cal.Weeks {
		x := contribLeftPad + weekIndex*step
		if len(week) > 0 {
			if month := week[0].Date.Month(); month != lastMonth {
				monthLabels = append(monthLabels, contributionMonthLabel{X: x, Label: week[0].Date.Format("Jan")})
				lastMonth = month
			}
		}
		for _, day := range week {
			y := contribTopPad + int(day.Date.Weekday())*step
			cells = append(cells, contributionCell{
				X:       x,
				Y:       y,
				Level:   day.Level,
				Tooltip: contributionTooltip(day),
			})
		}
	}

	weekdayLabels := []contributionWeekdayLabel{
		{Y: contribTopPad + int(time.Monday)*step + contribCellSize, Label: "Mon"},
		{Y: contribTopPad + int(time.Wednesday)*step + contribCellSize, Label: "Wed"},
		{Y: contribTopPad + int(time.Friday)*step + contribCellSize, Label: "Fri"},
	}

	return contributionGraph{
		Cells:         cells,
		MonthLabels:   monthLabels,
		WeekdayLabels: weekdayLabels,
		Total:         cal.Total,
		Width:         contribLeftPad + len(cal.Weeks)*step,
		Height:        contribTopPad + 7*step,
	}
}

func contributionTooltip(day contributionDay) string {
	count := "No contributions"
	switch day.Count {
	case 0:
	case 1:
		count = "1 contribution"
	default:
		count = fmt.Sprintf("%d contributions", day.Count)
	}
	return fmt.Sprintf("%s on %s", count, day.Date.Format("Jan 2, 2006"))
}

func renderHome(w io.Writer, now time.Time, repos []repo, pullRequests, issues, discussions []item, calendar contributionCalendar) error {
	tmpl := template.Must(template.ParseFS(templatesFS, "templates/layout.html", "templates/item.html", "templates/contributions.html", "templates/home.html"))
	data := struct {
		Title         string
		UpdatedAt     string
		Repos         []repo
		PullRequests  []item
		Issues        []item
		Discussions   []item
		Contributions contributionGraph
	}{
		Title:         siteTitle,
		UpdatedAt:     now.UTC().Format(timeFormat),
		Repos:         repos,
		PullRequests:  pullRequests[:min(len(pullRequests), sectionItemLimit)],
		Issues:        issues[:min(len(issues), sectionItemLimit)],
		Discussions:   discussions[:min(len(discussions), sectionItemLimit)],
		Contributions: buildContributionGraph(calendar),
	}
	return tmpl.ExecuteTemplate(w, "layout", data)
}

func renderSeeAll(w io.Writer, now time.Time, pageTitle string, items []item) error {
	tmpl := template.Must(template.ParseFS(templatesFS, "templates/layout.html", "templates/item.html", "templates/seeall.html"))
	data := struct {
		Title     string
		UpdatedAt string
		PageTitle string
		Groups    []yearGroup
	}{
		Title:     siteTitle + " - " + pageTitle,
		UpdatedAt: now.UTC().Format(timeFormat),
		PageTitle: pageTitle,
		Groups:    groupByYear(items),
	}
	return tmpl.ExecuteTemplate(w, "layout", data)
}

func groupByYear(items []item) []yearGroup {
	sorted := slices.Clone(items)
	slices.SortFunc(sorted, func(a, b item) int {
		return b.CreatedAt.Compare(a.CreatedAt)
	})

	var groups []yearGroup
	for _, it := range sorted {
		year := it.CreatedAt.UTC().Year()
		if len(groups) == 0 || groups[len(groups)-1].Year != year {
			groups = append(groups, yearGroup{Year: year})
		}
		groups[len(groups)-1].Items = append(groups[len(groups)-1].Items, it)
	}
	return groups
}
