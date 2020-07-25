package main

import (
	"frosh-api/handler"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/v2/github/repositories", handler.ListRepositories)
	http.HandleFunc("/v2/github/contributors", handler.ListContributors)

	log.Println("Go!")
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}
