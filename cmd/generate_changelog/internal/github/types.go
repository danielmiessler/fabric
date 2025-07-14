package github

import "time"

type PR struct {
	Number      int
	Title       string
	Body        string
	Author      string
	AuthorURL   string
	AuthorType  string // "user", "organization", or "bot"
	URL         string
	MergedAt    time.Time
	Commits     []PRCommit
	MergeCommit string
}

type PRCommit struct {
	SHA     string
	Message string
	Author  string
}

// GraphQL query structures for hasura client
type PullRequestsQuery struct {
	Repository struct {
		PullRequests struct {
			PageInfo struct {
				HasNextPage bool
				EndCursor   string
			}
			Nodes []struct {
				Number   int
				Title    string
				Body     string
				URL      string
				MergedAt time.Time
				Author   *struct {
					Typename string `graphql:"__typename"`
					Login    string `graphql:"login"`
					URL      string `graphql:"url"`
				}
				Commits struct {
					Nodes []struct {
						Commit struct {
							OID     string `graphql:"oid"`
							Message string
							Author  struct {
								Name string
							}
						}
					}
				} `graphql:"commits(first: 250)"`
			}
		} `graphql:"pullRequests(first: 100, after: $after, states: MERGED, orderBy: {field: UPDATED_AT, direction: DESC})"`
	} `graphql:"repository(owner: $owner, name: $repo)"`
}
