package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/julienschmidt/httprouter"

	"frosh-api/internal/shopware"
)

var SalesCache shopware.Sales
var RatingsCache shopware.Ratings

type ListSales struct {
	Total   int            `json:"all"`
	Plugins map[string]int `json:"plugins"`
}

func init() {
	go func() {
		for {
			refreshShopware()
			time.Sleep(time.Hour)
		}
	}()
}

func refreshShopware() {
	log.Println("Refreshing Shopware API data")
	token, err := shopware.Login(&shopware.LoginRequest{
		Email:    os.Getenv("SHOPWARE_USER"),
		Password: os.Getenv("SHOPWARE_PASSWORD"),
	})
	if err != nil {
		log.Println(err)
		return
	}

	sales, err := shopware.GetAllPluginSales(token)
	if err != nil {
		log.Println(err)
		return
	}
	SalesCache = sales

	ratings, err := shopware.GetAllRatings(token)
	if err != nil {
		log.Println(err)
		return
	}
	RatingsCache = ratings

	log.Println("Refreshed Shopware API data")
}

func ListPluginBuys(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	var saleList ListSales
	saleList.Plugins = make(map[string]int)
	saleList.Total = len(SalesCache)

	for _, s := range SalesCache {
		name := s.Plugin.Name

		if _, ok := saleList.Plugins[name]; !ok {
			saleList.Plugins[name] = 0
		}
		saleList.Plugins[name]++
	}

	jData, err := json.Marshal(saleList)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}
