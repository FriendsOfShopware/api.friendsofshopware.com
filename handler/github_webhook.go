package handler

import (
	"encoding/json"
	"fmt"
	"frosh-api/client"
	"github.com/julienschmidt/httprouter"
	"net/http"
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
	err := json.NewDecoder(r.Body).Decode(&d)
	if err != nil {
		// handle error
	}

	IssueCache[d.Repository.Name] = client.GetAllIssues(d.Repository.Owner.Login, d.Repository.Name)
	fmt.Printf("Updated issues for %s", d.Repository.Name)

	w.WriteHeader(200)
}
