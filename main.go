package main

import (
	"log"
	"net/http"
	"os"

	"github.com/getsentry/sentry-go"
	"github.com/julienschmidt/httprouter"

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

	log.Println("Go!")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", router))
}
