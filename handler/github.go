package handler

import (
	"encoding/json"
	"frosh-api/client"
	"github.com/google/go-github/v32/github"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"sort"
	"time"
)

const OrgName = "friendsofshopware"

var OrgCache = make(map[string][]*github.Repository)
var sortedContributors []*ContributionUser
var IssueCache = make(map[string][]*github.Issue)

type ContributionUser struct {
	User          string
	Name          string
	Contributions int
	Commits       int
	PullRequests  int
	AvatarURL     string
}

func init() {
	go func() {
		for {
			<-time.NewTicker(time.Hour).C
			refresh()
		}
	}()

	go func() {
		refresh()
	}()
}

func refresh() {
	log.Println("Refreshing Github Stats")
	OrgCache[OrgName] = client.AllRepos(OrgName)
	GetUserContributions()
	loadRepositoriesIssues()
	s := sortedContributors
	_ = s
	log.Println("Refreshed Github Stats")
}

func getRepositories() []*github.Repository {
	var repos []*github.Repository
	if entry, ok := OrgCache[OrgName]; ok {
		repos = entry
	} else {
		repos = client.AllRepos(OrgName)

		sort.Slice(repos, func(a, b int) bool {
			return repos[a].GetStargazersCount() > repos[b].GetStargazersCount()
		})

		OrgCache[OrgName] = repos
	}

	return repos
}

func ListRepositories(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	jData, err := json.Marshal(getRepositories())
	if err != nil {
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func GetUserContributions() []*ContributionUser {
	repos := getRepositories()

	var totalContributors = make(map[string]*ContributionUser)
	for _, repo := range repos {
		contributors, stats := client.GetContributors(repo.Owner.GetLogin(), repo.GetName())
		prs := client.GetPullRequests(repo.Owner.GetLogin(), repo.GetName())

		for _, c := range contributors {
			for _, s := range stats {
				if c.GetLogin() == s.Author.GetLogin() {
					username := c.GetLogin()
					entry, ok := totalContributors[username]
					if !ok {
						totalContributors[username] = &ContributionUser{
							User:      username,
							AvatarURL: s.Author.GetAvatarURL(),
						}
						entry = totalContributors[username]
					}

					entry.Commits += s.GetTotal()
					entry.Contributions += c.GetContributions()
				}
			}
		}

		for _, pr := range prs {
			entry, ok := totalContributors[pr.User.GetLogin()]

			if ok {
				entry.PullRequests++
			}
		}
	}

	for _, v := range totalContributors {
		v.Name = client.GetUser(v.User).GetName()
	}

	sortedContributors = make([]*ContributionUser, 0)
	for _, v := range totalContributors {
		sortedContributors = append(sortedContributors, v)
	}

	sort.Slice(sortedContributors, func(a, b int) bool {
		return sortedContributors[a].Contributions > sortedContributors[b].Contributions
	})

	return sortedContributors
}

func ListContributors(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	if len(sortedContributors) < 1 {
		GetUserContributions()
	}

	w.Header().Set("Content-Type", "application/json")
	jData, err := json.Marshal(sortedContributors)
	if err != nil {
	}

	w.Write(jData)
}

func loadRepositoriesIssues() {
	repos := getRepositories()

	for _, repo := range repos {
		IssueCache[repo.GetName()] = client.GetAllIssues(repo.Owner.GetLogin(), repo.GetName())
	}
}

func ListRepositoryIssues(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	issues, ok := IssueCache[ps.ByName("plugin")]

	if !ok {
		w.WriteHeader(404)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	jData, err := json.Marshal(issues)
	if err != nil {
	}

	w.Write(jData)
}
