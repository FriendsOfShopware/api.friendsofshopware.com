package handler

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"

	"frosh-api/internal/packagist"
)

var PackagesCache = make(map[string]*packagist.PackageStatistics)

func init() {
	go func() {
		for {
			PackagesCache = packagist.GetPackageStatistics()
			time.Sleep(time.Hour)
		}
	}()
}

func ListPackages(w http.ResponseWriter, _ *http.Request, ps httprouter.Params) {
	jData, err := json.Marshal(PackagesCache)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write(jData)
}
