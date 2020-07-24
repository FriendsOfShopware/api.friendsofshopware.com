package handler

import (
	"encoding/json"
	"frosh-api/client"
	"github.com/google/go-github/v32/github"
	"log"
	"net/http"
	"sort"
	"time"
)

const OrgName = "friendsofshopware"
var OrgCache = make(map[string][]*github.Repository)
var sortedContributors []*ContributionUser

type ContributionUser struct {
	User string
	Contributions int
	Commits int
	AvatarURL string
}

func init() {
	go func() {
		for {
			<-time.NewTicker(time.Hour).C
			refresh()
		}
	}()

	refresh()
}

func refresh() {
	log.Println("Refreshing shit")
	OrgCache[OrgName] = client.AllRepos(OrgName)
	GetUserContributions()
	s := sortedContributors
	_ = s
	log.Println("Refreshed shit")
}

func getRepositories() []*github.Repository {
	var repos []*github.Repository
	if entry, ok := OrgCache[OrgName]; ok {
		repos = entry
	} else {
		repos = client.AllRepos(OrgName)
		OrgCache[OrgName] = repos
	}

	return repos
}

func ListRepositories(w http.ResponseWriter, _ *http.Request) {
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
		contributors, stats := client.GetContributors(*repo.Owner.Login, *repo.Name)
		for _, c := range contributors {
			for _, s := range stats {
				if c.GetLogin() == s.Author.GetLogin() {
					username := c.GetLogin()
					entry, ok := totalContributors[username]
					if !ok {
						totalContributors[username] = &ContributionUser{
							User:          username,
							AvatarURL:     s.Author.GetAvatarURL(),
						}
						entry = totalContributors[username]
					}

					entry.Commits += s.GetTotal()
					entry.Contributions += c.GetContributions()
				}
			}
		}
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

func ListContributors(w http.ResponseWriter, _ *http.Request)  {
	if len(sortedContributors) < 1 {
		GetUserContributions()
	}

	w.Header().Set("Content-Type", "application/json")
	jData, err := json.Marshal(sortedContributors)
	if err != nil {
	}

	w.Write(jData)
}