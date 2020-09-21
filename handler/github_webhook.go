package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"

	githubClient "frosh-api/internal/github"
)

type Webhook struct {
	Action     string `json:"action"`
	Repository struct {
		Name  string `json:"name"`
		Owner struct {
			Login string `json:"login"`
		} `json:"owner"`
	} `json:"repository"`
}

func GithubIssueWebhook(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	var d Webhook
	if err := json.NewDecoder(r.Body).Decode(&d); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	IssueCache[d.Repository.Name] = githubClient.GetAllIssues(d.Repository.Owner.Login, d.Repository.Name)
	fmt.Printf("Updated issues for %s", d.Repository.Name)

	w.WriteHeader(200)
}
