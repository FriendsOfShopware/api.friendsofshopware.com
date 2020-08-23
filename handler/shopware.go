package handler

import (
	"encoding/json"
	"frosh-api/client/shopware"
	_struct "frosh-api/client/shopware/struct"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
	"time"
)

var SalesCache _struct.Sales
var RatingsCache _struct.Ratings

type ListSales struct {
	Total   int            `json:"all"`
	Plugins map[string]int `json:"plugins"`
}

func init() {
	go func() {
		for {
			<-time.NewTicker(time.Hour).C
			refreshShopware()
		}
	}()

	go func() {
		refreshShopware()
	}()
}

func refreshShopware() {
	log.Println("Refreshing Shopware API data")
	token := shopware.Login(_struct.LoginRequest{
		Email:    os.Getenv("SHOPWARE_USER"),
		Password: os.Getenv("SHOPWARE_PASSWORD"),
	})

	SalesCache = shopware.GetAllPluginSales(token)
	RatingsCache = shopware.GetAllRatings(token)
	log.Println("Refreshed Shopware API data")
}

func ListPluginBuys(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	var saleList ListSales
	saleList.Plugins = make(map[string]int)
	saleList.Total = len(SalesCache)

	for _, s := range SalesCache {
		_, ok := saleList.Plugins[s.Plugin.Name]

		if !ok {
			saleList.Plugins[s.Plugin.Name] = 0
		}
		saleList.Plugins[s.Plugin.Name]++
	}

	jData, err := json.Marshal(saleList)
	if err != nil {
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}

func ListPluginRatings(w http.ResponseWriter, _ *http.Request) {
	jData, err := json.Marshal(RatingsCache)
	if err != nil {
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}
