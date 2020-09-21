package shopware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

const ApiUrl = "https://api.shopware.com"

func Login(request *LoginRequest) (*Token, error) {
	s, err := json.Marshal(request)
	if err != nil {
		return nil, fmt.Errorf("login: %v", err)
	}

	resp, err := http.Post(ApiUrl+"/accesstokens", "application/json", bytes.NewBuffer(s))
	if err != nil {
		return nil, fmt.Errorf("login: %v", err)
	}

	defer resp.Body.Close()

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("login: %v", err)
	}

	var token Token
	if err := json.Unmarshal(data, &token); err != nil {
		return nil, fmt.Errorf("login: %v", err)
	}

	return &token, nil
}

func GetAllPluginSales(token *Token) (Sales, error) {
	getPage := func(limit, offset int) (Sales, error) {
		url := fmt.Sprintf("%s/producers/2287/sales?limit=%d&offset=%d&orderBy=creationDate&orderSequence=desc&search=&variantType=free", ApiUrl, limit, offset)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("getAllPluginSales: %v", err)
		}
		req.Header.Set("X-Shopware-Token", token.Token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("getAllPluginSales: %v", err)
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("getAllPluginSales: %v", err)
		}

		var sales Sales
		if err := json.Unmarshal(data, &sales); err != nil {
			return nil, fmt.Errorf("getAllPluginSales: %v", err)
		}

		return sales, nil
	}

	offset, limit := 0, 100
	var allSales Sales
	for {
		sales, err := getPage(limit, offset)
		if err != nil {
			return nil, err
		}

		for _, s := range sales {
			allSales = append(allSales, s)
		}

		if len(sales) == 0 {
			break
		}

		offset += limit
	}

	return allSales, nil
}

func GetAllRatings(token *Token) (Ratings, error) {
	getPage := func(limit, offset int) (Ratings, error) {
		url := fmt.Sprintf("%s/plugincomments?limit=%d&offset=%d&orderBy=creationDate&orderSequence=desc&search=&producerId=2287", ApiUrl, limit, offset)
		req, err := http.NewRequest("GET", url, nil)
		if err != nil {
			return nil, fmt.Errorf("getAllPluginSales: %v", err)
		}
		req.Header.Set("X-Shopware-Token", token.Token)

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, fmt.Errorf("getAllPluginSales: %v", err)
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("getAllPluginSales: %v", err)
		}

		var ratings Ratings
		if err := json.Unmarshal(data, &ratings); err != nil {
			return nil, fmt.Errorf("getAllPluginSales: %v", err)
		}

		return ratings, nil
	}

	offset, limit := 0, 100
	var allRatings Ratings
	for {
		ratings, err := getPage(limit, offset)
		if err != nil {
			return nil, err
		}

		for _, s := range ratings {
			allRatings = append(allRatings, s)
		}

		if len(ratings) == 0 {
			break
		}

		offset += limit
	}

	return allRatings, nil
}
