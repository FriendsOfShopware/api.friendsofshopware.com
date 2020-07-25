package client

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"
)

const packagesUrl = "https://packagist.org/"

type PackageList struct {
	PackageNames []string `json:"packageNames"`
}

type PackageDetail struct {
	Package struct {
		Name        string    `json:"name"`
		Description string    `json:"description"`
		Time        time.Time `json:"time"`
		Maintainers []struct {
			Name      string `json:"name"`
			AvatarURL string `json:"avatar_url"`
		} `json:"maintainers"`
		Versions struct {
			DevMaster struct {
				Name              string        `json:"name"`
				Description       string        `json:"description"`
				Keywords          []string      `json:"keywords"`
				Homepage          string        `json:"homepage"`
				Version           string        `json:"version"`
				VersionNormalized string        `json:"version_normalized"`
				License           []string      `json:"license"`
				Authors           []interface{} `json:"authors"`
				Source            struct {
					Type      string `json:"type"`
					URL       string `json:"url"`
					Reference string `json:"reference"`
				} `json:"source"`
				Dist struct {
					Type      string `json:"type"`
					URL       string `json:"url"`
					Reference string `json:"reference"`
					Shasum    string `json:"shasum"`
				} `json:"dist"`
				Type    string `json:"type"`
				Support struct {
					Source string `json:"source"`
					Issues string `json:"issues"`
				} `json:"support"`
				Time  time.Time `json:"time"`
				Extra struct {
					InstallerName       string `json:"installer-name"`
					ShopwarePluginClass string `json:"shopware-plugin-class"`
				} `json:"extra"`
				DefaultBranch bool `json:"default-branch"`
				Require       struct {
					Php                string `json:"php"`
					ComposerInstallers string `json:"composer/installers"`
				} `json:"require"`
			} `json:"dev-master"`
		} `json:"versions"`
		Type             string           `json:"type"`
		Repository       string           `json:"repository"`
		GithubStars      int              `json:"github_stars"`
		GithubWatchers   int              `json:"github_watchers"`
		GithubForks      int              `json:"github_forks"`
		GithubOpenIssues int              `json:"github_open_issues"`
		Language         string           `json:"language"`
		Dependents       int              `json:"dependents"`
		Suggesters       int              `json:"suggesters"`
		Downloads        PackageDownloads `json:"downloads"`
		Favers           int              `json:"favers"`
	} `json:"package"`
}

type PackageDownloads struct {
	Total   int `json:"total"`
	Monthly int `json:"monthly"`
	Daily   int `json:"daily"`
}

type PackageStatistics struct {
	Github struct {
		Stars    int
		Watchers int
		Forks    int
	}
	Downloads PackageDownloads
}

func request(url string) []byte {
	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

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
