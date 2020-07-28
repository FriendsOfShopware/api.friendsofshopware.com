package main

import (
	"frosh-api/handler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/v2/github/repositories", handler.ListRepositories)
	http.HandleFunc("/v2/github/contributors", handler.ListContributors)
	http.HandleFunc("/v2/packagist/packages", handler.ListPackages)
	http.HandleFunc("/v2/shopware/sales", handler.ListPluginBuys)
	http.HandleFunc("/v2/shopware/badge", handler.GetStoreDownloadBadge)
	//http.HandleFunc("/", handler.ListPluginBuys)
	//http.HandleFunc("/v2/shopware/ratings", handler.ListPluginRatings)

	log.Println("Go!")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
