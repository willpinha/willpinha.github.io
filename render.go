package main

import (
	"embed"
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
	Number    int
	Title     string
	URL       string
	State     string
	CreatedAt time.Time
	RepoName  string
	RepoURL   string
}

type yearGroup struct {
	Year  int
	Items []item
}

func renderHome(w io.Writer, now time.Time, repos []repo, pullRequests, issues, discussions []item) error {
	tmpl := template.Must(template.ParseFS(templatesFS, "templates/layout.html", "templates/item.html", "templates/home.html"))
	data := struct {
		Title        string
		UpdatedAt    string
		Repos        []repo
		PullRequests []item
		Issues       []item
		Discussions  []item
	}{
		Title:        siteTitle,
		UpdatedAt:    now.UTC().Format(timeFormat),
		Repos:        repos,
		PullRequests: pullRequests[:min(len(pullRequests), sectionItemLimit)],
		Issues:       issues[:min(len(issues), sectionItemLimit)],
		Discussions:  discussions[:min(len(discussions), sectionItemLimit)],
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
