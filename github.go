package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

const (
	graphqlEndpoint = "https://api.github.com/graphql"

	// The GitHub search API never returns more than 1000 results per query
	searchResultCap = 1000
)

const repoSearchQuery = `
query ($search: String!, $cursor: String) {
	search(type: REPOSITORY, query: $search, first: 100, after: $cursor) {
		pageInfo {
			hasNextPage
			endCursor
		}
		nodes {
			... on Repository {
				nameWithOwner
				url
				stargazerCount
				isPrivate
			}
		}
	}
}`

const issueSearchQuery = `
query ($search: String!, $cursor: String) {
	search(type: ISSUE, first: 100, query: $search, after: $cursor) {
		pageInfo {
			hasNextPage
			endCursor
		}
		nodes {
			... on Issue {
				number
				title
				url
				createdAt
				repository {
					nameWithOwner
					url
					isPrivate
				}
			}
			... on PullRequest {
				number
				title
				url
				createdAt
				state
				repository {
					nameWithOwner
					url
					isPrivate
				}
			}
		}
	}
}`

const discussionSearchQuery = `
query ($search: String!, $cursor: String) {
	search(type: DISCUSSION, first: 100, query: $search, after: $cursor) {
		pageInfo {
			hasNextPage
			endCursor
		}
		nodes {
			... on Discussion {
				number
				title
				url
				createdAt
				repository {
					nameWithOwner
					url
					isPrivate
				}
			}
		}
	}
}`

type client struct {
	token      string
	httpClient *http.Client
}

func newClient(token string) *client {
	return &client{
		token:      token,
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

func (c *client) famousRepos() ([]repo, error) {
	search := fmt.Sprintf(
		"user:%s stars:>=%d fork:false is:public sort:stars-desc",
		login,
		minRepoStars,
	)
	nodes, err := searchAll[repoNode](c, repoSearchQuery, search)
	if err != nil {
		return nil, err
	}
	repos := make([]repo, 0, len(nodes))
	for _, n := range nodes {
		if n.IsPrivate {
			continue
		}
		repos = append(repos, repo{Name: n.NameWithOwner, URL: n.URL, Stars: n.StargazerCount})
	}
	return repos, nil
}

func (c *client) createdPullRequests() ([]item, error) {
	search := fmt.Sprintf("is:pr is:public author:%s -user:%s sort:created-desc", login, login)
	return c.searchItems(issueSearchQuery, search)
}

func (c *client) participatedIssues() ([]item, error) {
	search := fmt.Sprintf("is:issue is:public involves:%s sort:updated-desc", login)
	return c.searchItems(issueSearchQuery, search)
}

func (c *client) participatedDiscussions() ([]item, error) {
	search := fmt.Sprintf("involves:%s sort:updated-desc", login)
	return c.searchItems(discussionSearchQuery, search)
}

func (c *client) searchItems(query, search string) ([]item, error) {
	nodes, err := searchAll[itemNode](c, query, search)
	if err != nil {
		return nil, err
	}
	items := make([]item, 0, len(nodes))
	for _, n := range nodes {
		// Guards against private results when running with a broadly scoped local token
		if n.Repository.IsPrivate {
			continue
		}
		items = append(items, item{
			Number:    n.Number,
			Title:     n.Title,
			URL:       n.URL,
			State:     strings.ToLower(n.State),
			CreatedAt: n.CreatedAt,
			RepoName:  n.Repository.NameWithOwner,
			RepoURL:   n.Repository.URL,
		})
	}
	return items, nil
}

type pageInfo struct {
	HasNextPage bool   `json:"hasNextPage"`
	EndCursor   string `json:"endCursor"`
}

type repoNode struct {
	NameWithOwner  string `json:"nameWithOwner"`
	URL            string `json:"url"`
	StargazerCount int    `json:"stargazerCount"`
	IsPrivate      bool   `json:"isPrivate"`
}

type itemNode struct {
	Number     int       `json:"number"`
	Title      string    `json:"title"`
	URL        string    `json:"url"`
	CreatedAt  time.Time `json:"createdAt"`
	State      string    `json:"state"`
	Repository struct {
		NameWithOwner string `json:"nameWithOwner"`
		URL           string `json:"url"`
		IsPrivate     bool   `json:"isPrivate"`
	} `json:"repository"`
}

func searchAll[T any](c *client, query, search string) ([]T, error) {
	var all []T
	var cursor *string
	for {
		var result struct {
			Search struct {
				PageInfo pageInfo `json:"pageInfo"`
				Nodes    []T      `json:"nodes"`
			} `json:"search"`
		}
		variables := map[string]any{"search": search, "cursor": cursor}
		if err := c.query(query, variables, &result); err != nil {
			return nil, err
		}
		all = append(all, result.Search.Nodes...)
		if !result.Search.PageInfo.HasNextPage || len(all) >= searchResultCap {
			return all, nil
		}
		cursor = &result.Search.PageInfo.EndCursor
	}
}

func (c *client) query(query string, variables map[string]any, out any) error {
	payload, err := json.Marshal(map[string]any{"query": query, "variables": variables})
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPost, graphqlEndpoint, bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+c.token)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("graphql request returned status %s", resp.Status)
	}

	var envelope struct {
		Data   json.RawMessage `json:"data"`
		Errors []struct {
			Message string `json:"message"`
		} `json:"errors"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&envelope); err != nil {
		return err
	}
	if len(envelope.Errors) > 0 {
		return fmt.Errorf("graphql error: %s", envelope.Errors[0].Message)
	}
	return json.Unmarshal(envelope.Data, out)
}
