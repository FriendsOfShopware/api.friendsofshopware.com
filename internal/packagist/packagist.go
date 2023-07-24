package packagist

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"strings"
)

const packagesUrl = "https://packagist.org/"

func request(url string) []byte {
	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	return body
}

func getPluginNameFromPackage(detail PackageDetail) string {
	if len(detail.Package.Versions.DevMaster.Extra.ShopwarePluginClass) > 0 {
		split := strings.Split(detail.Package.Versions.DevMaster.Extra.ShopwarePluginClass, "\\")
		return split[len(split)-1]
	}

	return detail.Package.Versions.DevMaster.Extra.InstallerName
}

func GetPackageStatistics() map[string]*PackageStatistics {
	var Packages = make(map[string]*PackageStatistics)

	var packageList PackageList
	err := json.Unmarshal(request(packagesUrl+"packages/list.json?vendor=frosh"), &packageList)
	if err != nil {
		log.Println(err)
	}

	for _, name := range packageList.PackageNames {
		var pkg PackageDetail
		err := json.Unmarshal(request(packagesUrl+"packages/"+name+".json"), &pkg)
		if err != nil {
			log.Println(err)
		}

		pluginName := getPluginNameFromPackage(pkg)

		if len(pluginName) == 0 {
			continue
		}

		Packages[pluginName] = &PackageStatistics{
			Downloads: pkg.Package.Downloads,
			Github: struct {
				Stars    int
				Watchers int
				Forks    int
			}{Stars: pkg.Package.GithubStars, Watchers: pkg.Package.GithubWatchers, Forks: pkg.Package.GithubForks},
		}
	}

	return Packages
}
