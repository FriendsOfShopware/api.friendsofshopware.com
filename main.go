package main

import (
	"frosh-api/handler"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
)

func main() {
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
