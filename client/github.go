package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/google/go-github/v32/github"
	"golang.org/x/oauth2"
	"log"
	"net/http"
	"os"
)

var client *github.Client
var ctx = context.TODO()

func init()  {
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(context.Background(), ts)

	t := tc.Transport.(*oauth2.Transport)
	t.Base = &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}

	client = github.NewClient(tc)
	ctx = context.Background()
}

func AllRepos(organisation string) []*github.Repository {
	opt := &github.RepositoryListByOrgOptions{
		ListOptions: github.ListOptions{PerPage: 10},
	}
	// get all pages of results
	var allRepos []*github.Repository
	for {
		repos, resp, err := client.Repositories.ListByOrg(ctx, organisation, opt)
		if err != nil {
			log.Fatal(fmt.Errorf("error while getting all repos: %w", err))
			return nil
		}
		allRepos = append(allRepos, repos...)
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
		log.Fatal(fmt.Errorf("error while getting contributor stats: %w", err))
	}
	return list, stats
}