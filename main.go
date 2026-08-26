package main

import (
	"flag"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"
)

const (
	login        = "willpinha"
	minRepoStars = 10
	siteTitle    = "Willian Pinheiro"
)

type page struct {
	path   string
	render func(w io.Writer) error
}

func main() {
	outDir := flag.String("out", "dist", "path to the build output directory")
	server := flag.Bool("server", false, "serve the build output directory instead of generating it")
	port := flag.Int("port", 8080, "port to listen on when -server is set")
	mock := flag.Bool("mock", false, "generate the site with mock data instead of querying the GitHub API")
	flag.Parse()

	if *server {
		if err := runServer(*outDir, *port); err != nil {
			log.Fatalf("serve %s: %v", *outDir, err)
		}
		return
	}

	var repos []repo
	var pullRequests, issues, discussions []item

	if *mock {
		repos = mockRepos
		pullRequests = mockPullRequests
		issues = mockIssues
		discussions = mockDiscussions
	} else {
		token := os.Getenv("GITHUB_TOKEN")
		if token == "" {
			log.Fatal("GITHUB_TOKEN environment variable is not set")
		}

		client := newClient(token)

		var err error
		repos, err = client.famousRepos()
		if err != nil {
			log.Fatalf("fetch repositories: %v", err)
		}
		pullRequests, err = client.createdPullRequests()
		if err != nil {
			log.Fatalf("fetch pull requests: %v", err)
		}
		issues, err = client.participatedIssues()
		if err != nil {
			log.Fatalf("fetch issues: %v", err)
		}
		discussions, err = client.participatedDiscussions()
		if err != nil {
			log.Fatalf("fetch discussions: %v", err)
		}
	}

	if err := os.RemoveAll(*outDir); err != nil {
		log.Fatalf("clean %s: %v", *outDir, err)
	}
	if err := os.CopyFS(filepath.Join(*outDir, "assets"), os.DirFS("assets")); err != nil {
		log.Fatalf("copy assets: %v", err)
	}

	now := time.Now()

	pages := []page{
		{
			path: "index.html",
			render: func(w io.Writer) error {
				return renderHome(w, now, repos, pullRequests, issues, discussions)
			},
		},
		{
			path: filepath.Join("pull-requests", "index.html"),
			render: func(w io.Writer) error {
				return renderSeeAll(w, now, "Pull requests I created", pullRequests)
			},
		},
		{
			path: filepath.Join("issues", "index.html"),
			render: func(w io.Writer) error {
				return renderSeeAll(w, now, "Issues I participated in", issues)
			},
		},
		{
			path: filepath.Join("discussions", "index.html"),
			render: func(w io.Writer) error {
				return renderSeeAll(w, now, "Discussions I participated in", discussions)
			},
		},
	}

	for _, p := range pages {
		path := filepath.Join(*outDir, p.path)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			log.Fatalf("create directory for %s: %v", path, err)
		}
		f, err := os.Create(path)
		if err != nil {
			log.Fatalf("create %s: %v", path, err)
		}
		if err := p.render(f); err != nil {
			f.Close()
			log.Fatalf("render %s: %v", path, err)
		}
		if err := f.Close(); err != nil {
			log.Fatalf("close %s: %v", path, err)
		}
		log.Printf("wrote %s", path)
	}
}
