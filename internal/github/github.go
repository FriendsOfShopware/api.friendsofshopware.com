package github

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v53/github"
	"golang.org/x/oauth2"
)

var client = github.NewClient(
	oauth2.NewClient(
		context.Background(),
		oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
		),
	),
)

func AllRepos(organisation string) []*github.Repository {
	opt := &github.RepositoryListByOrgOptions{
		Type:        "public",
		ListOptions: github.ListOptions{PerPage: 100},
	}
	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(context.TODO(), organisation, opt)
		if err != nil {
			log.Println(fmt.Errorf("error while getting all repos: %w", err))
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
	list, _, err := client.Repositories.ListContributors(context.TODO(), owner, repository, nil)
	if err != nil {
		log.Println(fmt.Errorf("error while getting contributors: %w", err))
		return nil, nil
	}

	stats, resp, err := client.Repositories.ListContributorsStats(context.TODO(), owner, repository)
	if err != nil || resp == nil {
		log.Println(fmt.Errorf("error while getting contributor stats: %w", err))
		return nil, nil
	}

	if resp.StatusCode != http.StatusOK {
		if resp.StatusCode == http.StatusAccepted || strings.Contains(err.Error(), "job scheduled on GitHub side") {
			fmt.Println("Got job scheduled message error. Waiting some time to wait")
			time.Sleep(1 * time.Minute)
			return GetContributors(owner, repository)
		}

		log.Println(fmt.Errorf("error while getting contributor stats: %w", err))
		return nil, nil
	}

	return list, stats
}

func GetUser(login string) *github.User {
	user, _, err := client.Users.Get(context.TODO(), login)
	if err != nil {
		log.Println(fmt.Errorf("error while getting user: %w", err))
		return nil
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
		repos, resp, err := client.PullRequests.List(context.TODO(), owner, repository, opt)

		if err != nil {
			log.Println(fmt.Errorf("error while getting all repos: %w", err))
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
		issues, resp, err := client.Issues.ListByRepo(context.TODO(), owner, repository, opt)

		if err != nil {
			log.Println(fmt.Errorf("error while getting all issues: %w", err))
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
