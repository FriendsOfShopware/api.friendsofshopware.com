package handler

import (
	"github.com/julienschmidt/httprouter"
	"github.com/narqo/go-badge"
	"net/http"
	"strconv"
)

func GetStoreDownloadBadge(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	pluginDownloadList := make(map[string]int)

	for _, s := range SalesCache {
		_, ok := pluginDownloadList[s.Plugin.Name]

		if !ok {
			pluginDownloadList[s.Plugin.Name] = 0
		}
		pluginDownloadList[s.Plugin.Name]++
	}

	_, ok := pluginDownloadList[ps.ByName("plugin")]

	if !ok {
		w.WriteHeader(404)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml;charset=utf-8")

	badge, _ := badge.RenderBytes("Store Downloads", strconv.Itoa(pluginDownloadList[ps.ByName("plugin")])+" Downloads", "#189eff")

	w.Write(badge)
}
