package handler

import (
	"encoding/json"
	"frosh-api/client"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"time"
)

var PackagesCache = make(map[string]*client.PackageStatistics)

func init() {
	go func() {
		for {
			<-time.NewTicker(time.Hour).C
			PackagesCache = client.GetPackageStatistics()
		}
	}()

	go func() {
		PackagesCache = client.GetPackageStatistics()
	}()
}

func ListPackages(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	jData, err := json.Marshal(PackagesCache)
	if err != nil {
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(jData)
}
