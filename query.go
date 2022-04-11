// Interact with the GitHub API.

package main

import (
	"log"
	"time"

	"github.com/cli/go-gh"
	"github.com/cli/go-gh/pkg/api"
	graphql "github.com/cli/shurcooL-graphql"
	"github.com/lithammer/fuzzysearch/fuzzy"
)

// The GitHub repository to query.
type Repository struct {
	Owner string
	Name  string
}

// Wrapper around an API query client in order to enable mocking.
type Querier interface {

	// Return the next set of results.
	Next() map[string]int

	// Returns true if there are no more matches.
	AtEnd() bool
}

// GitHub GraphQL API query client.
type GQLPRQuerier struct {
	client   api.GQLClient
	repo     Repository
	numfetch int
	cursor   string
	hasmore  bool
}

// Create a GitHub GraphQL PR Querier.
func CreateGQLPRQuerier(repo Repository, timeout time.Duration, numfetch int) *GQLPRQuerier {
	opts := &api.ClientOptions{
		EnableCache: true,
		Timeout:     timeout,
	}
	client, err := gh.GQLClient(opts)
	if err != nil {
		log.Fatal(err)
	}
	return &GQLPRQuerier{client, repo, numfetch, "", true}
}

func (q *GQLPRQuerier) Next() map[string]int {
	matchmap := make(map[string]int, q.numfetch)

	if q.cursor == "" {
		var query struct {
			Repository struct {
				PullRequests struct {
					Nodes []struct {
						Title  string
						Number int
					}
					PageInfo struct {
						StartCursor     string
						HasPreviousPage bool
					}
				} `graphql:"pullRequests(last: $last)"`
			} `graphql:"repository(owner: $owner, name: $name)"`
		}
		variables := map[string]interface{}{
			"last":  graphql.Int(q.numfetch),
			"owner": graphql.String(q.repo.Owner),
			"name":  graphql.String(q.repo.Name),
		}
		err := q.client.Query("RepositoryPRs", &query, variables)
		if err != nil {
			log.Fatal(err)
		}
		for _, node := range query.Repository.PullRequests.Nodes {
			matchmap[node.Title] = node.Number
		}
		q.cursor = query.Repository.PullRequests.PageInfo.StartCursor
		q.hasmore = query.Repository.PullRequests.PageInfo.HasPreviousPage
	} else {
		var query struct {
			Repository struct {
				PullRequests struct {
					Nodes []struct {
						Title  string
						Number int
					}
					PageInfo struct {
						StartCursor     string
						HasPreviousPage bool
					}
				} `graphql:"pullRequests(last: $last, before: $before)"`
			} `graphql:"repository(owner: $owner, name: $name)"`
		}
		variables := map[string]interface{}{
			"last":   graphql.Int(q.numfetch),
			"before": graphql.String(q.cursor),
			"owner":  graphql.String(q.repo.Owner),
			"name":   graphql.String(q.repo.Name),
		}
		err := q.client.Query("RepositoryPRs", &query, variables)
		if err != nil {
			log.Fatal(err)
		}

		for _, node := range query.Repository.PullRequests.Nodes {
			matchmap[node.Title] = node.Number
		}
		q.cursor = query.Repository.PullRequests.PageInfo.StartCursor
		q.hasmore = query.Repository.PullRequests.PageInfo.HasPreviousPage
	}

	return matchmap
}

func (q *GQLPRQuerier) AtEnd() bool {
	return !q.hasmore
}

func findMatches(search string, candidates []string) []string {
	return fuzzy.FindFold(search, candidates)
}

func filterMap(search string, matchmap map[string]int) map[string]int {
	candidates := make([]string, len(matchmap))
	i := 0
	for title := range matchmap {
		candidates[i] = title
		i++
	}
	tmatches := findMatches(search, candidates)

	ret := make(map[string]int, len(tmatches))

	for _, m := range tmatches {
		ret[m] = matchmap[m]
	}

	return ret
}

// Find the PRs which match a search.
func FindMatchingPRs(querier Querier, search string) map[string]int {
	for {
		results := querier.Next()
		matches := filterMap(search, results)
		if len(matches) > 0 {
			return matches
		} else if querier.AtEnd() {
			break
		}
	}

	return nil
}
