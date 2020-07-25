package shopware

import (
	"bytes"
	"encoding/json"
	"fmt"
	_struct "frosh-api/client/shopware/struct"
	"io/ioutil"
	"log"
	"net/http"
)

const ShopwareApiUrl = "https://api.shopware.com/"

func Login(request _struct.LoginRequest) _struct.Token {
	jsonStr, _ := json.Marshal(request)

	req, err := http.NewRequest("POST", ShopwareApiUrl+"accesstokens", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	var token _struct.Token

	json.Unmarshal(body, &token)

	return token
}

func GetAllPluginSales(token _struct.Token) _struct.Sales {
	offset := 0
	limit := 100

	client := &http.Client{}

	var allSales _struct.Sales

	for {
		url := fmt.Sprintf(ShopwareApiUrl+"producers/2287/sales?limit=%d&offset=%d&orderBy=creationDate&orderSequence=desc&search=&variantType=free", limit, offset)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("X-Shopware-Token", token.Token)

		resp, err := client.Do(req)

		if err != nil {
			log.Fatal(err)
		}

		body, _ := ioutil.ReadAll(resp.Body)

		var requestSales _struct.Sales

		json.Unmarshal(body, &requestSales)

		for _, s := range requestSales {
			allSales = append(allSales, s)
		}

		if len(requestSales) == 0 {
			break
		}

		offset += limit
	}

	return allSales
}

func GetAllRatings(token _struct.Token) _struct.Ratings {
	offset := 0
	limit := 100

	client := &http.Client{}

	var allRatings _struct.Ratings

	for {
		url := fmt.Sprintf(ShopwareApiUrl+"plugincomments?limit=%d&offset=%d&orderBy=creationDate&orderSequence=desc&search=&producerId=2287", limit, offset)
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("X-Shopware-Token", token.Token)

		if err != nil {

			resp, err := client.Do(req)
			log.Fatal(err)
		}

		body, _ := ioutil.ReadAll(resp.Body)

		var requestRatings _struct.Ratings

		json.Unmarshal(body, &requestRatings)

		for _, s := range requestRatings {
			allRatings = append(allRatings, s)
		}

		if len(requestRatings) == 0 {
			break
		}

		offset += limit
	}

	return allRatings
}
