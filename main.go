package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"

	"frosh-api/handler"
)

func main() {
	router := httprouter.New()
	router.GET("/v2/github/repositories", handler.ListRepositories)
	router.GET("/v2/github/contributors", handler.ListContributors)
	router.GET("/v2/github/issues/:plugin", handler.ListRepositoryIssues)
	router.GET("/v2/packagist/packages", handler.ListPackages)
	router.POST("/webhook/issue", handler.GithubIssueWebhook)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	// Insert the middleware
	handler := cors.Default().Handler(router)

	log.Println("Go!")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), handler))
}
