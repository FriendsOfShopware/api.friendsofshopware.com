package handler

import (
	"github.com/narqo/go-badge"
	"net/http"
	"strconv"
)

func GetStoreDownloadBadge(w http.ResponseWriter, r *http.Request) {
	keys, ok := r.URL.Query()["plugin"]

	if !ok || len(keys[0]) < 1 {
		w.WriteHeader(204)
		return
	}

	pluginDownloadList := make(map[string]int)

	for _, s := range SalesCache {
		_, ok := pluginDownloadList[s.Plugin.Name]

		if !ok {
			pluginDownloadList[s.Plugin.Name] = 0
		}
		pluginDownloadList[s.Plugin.Name]++
	}

	_, ok = pluginDownloadList[keys[0]]

	if !ok {
		w.WriteHeader(204)
		return
	}

	w.Header().Set("Content-Type", "image/svg+xml;charset=utf-8")

	badge, _ := badge.RenderBytes("Store Downloads", strconv.Itoa(pluginDownloadList[keys[0]])+" Downloads", "#189eff")

	w.Write(badge)
}
