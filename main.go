package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/julienschmidt/httprouter"
	"github.com/rs/cors"

	"frosh-api/handler"
)

func main() {
	err := sentry.Init(sentry.ClientOptions{
		Dsn: os.Getenv("SENTRY_URL"),
	})
	if err != nil {
		log.Fatalf("sentry.Init: %s", err)
	}

	router := httprouter.New()
	router.GET("/v2/github/repositories", handler.ListRepositories)
	router.GET("/v2/github/contributors", handler.ListContributors)
	router.GET("/v2/github/issues/:plugin", handler.ListRepositoryIssues)
	router.GET("/v2/packagist/packages", handler.ListPackages)
	router.GET("/v2/shopware/sales", handler.ListPluginBuys)
	router.GET("/v2/shopware/badge/:plugin", handler.GetStoreDownloadBadge)
	router.POST("/webhook/issue", handler.GithubIssueWebhook)
	router.GET("/", handler.ListPluginBuys)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	// Insert the middleware
	handler := cors.Default().Handler(router)

	log.Println("Go!")
	log.Fatal(http.ListenAndServe(fmt.Sprintf("0.0.0.0:%s", port), handler))
}
