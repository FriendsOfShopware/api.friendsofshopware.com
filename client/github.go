package client

import (
	"context"
	"fmt"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
	"log"
	"os"
	"strings"
	"time"
)

var client *github.Client
var ctx = context.TODO()

func init() {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	client = github.NewClient(tc)
	ctx = context.Background()
}

func AllRepos(organisation string) []*github.Repository {
	opt := &github.RepositoryListByOrgOptions{
		Type:        "public",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, organisation, opt)
		if err != nil {
			log.Fatal(fmt.Errorf("error while getting all repos: %w", err))
			return nil
		}

		for _, repo := range repos {
			if repo.GetArchived() {
				continue
			}

			allRepos = append(allRepos, repo)
		}

		if resp.NextPage == 0 {
			break
		}
		opt.Page = resp.NextPage
	}

	return allRepos
}

func GetContributors(owner, repository string) ([]*github.Contributor, []*github.ContributorStats) {
	list, _, err := client.Repositories.ListContributors(ctx, owner, repository, nil)
	if err != nil {
		log.Fatal(fmt.Errorf("error while getting contributors: %w", err))
	}

	stats, _, err := client.Repositories.ListContributorsStats(ctx, owner, repository)
	if err != nil {
		errorMsg := fmt.Sprintf("%w", err)

		if strings.Contains(errorMsg, "job scheduled on GitHub side") {
			fmt.Println("Got job scheduled message error. Waiting some time to wait")
			time.Sleep(30 * time.Second)
			return GetContributors(owner, repository)
		}

		log.Fatal(fmt.Errorf("error while getting contributor stats: %w", err))
	}
	return list, stats
}

func GetUser(login string) *github.User {
	user, _, err := client.Users.Get(ctx, login)

	if err != nil {
		log.Fatal(fmt.Errorf("error while getting user: %w", err))
	}

	return user
}

func GetPullRequests(owner, repository string) []*github.PullRequest {
	opt := &github.PullRequestListOptions{
		State:       "all",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	// get all pages of results
	var allPullRequests []*github.PullRequest

	for {
		repos, resp, err := client.PullRequests.List(ctx, owner, repository, opt)

		if err != nil {
			log.Fatal(fmt.Errorf("error while getting all repos: %w", err))
			return allPullRequests
		}

		allPullRequests = append(allPullRequests, repos...)

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return allPullRequests
}

func GetAllIssues(owner, repository string) []*github.Issue {
	opt := &github.IssueListByRepoOptions{
		State:       "open",
		ListOptions: github.ListOptions{PerPage: 100},
	}

	var allIssues []*github.Issue

	for {
		issues, resp, err := client.Issues.ListByRepo(ctx, owner, repository, opt)

		if err != nil {
			log.Fatal(fmt.Errorf("error while getting all issues: %w", err))
			return allIssues
		}

		allIssues = append(allIssues, issues...)

		if resp.NextPage == 0 {
			break
		}

		opt.Page = resp.NextPage
	}

	return allIssues
}
