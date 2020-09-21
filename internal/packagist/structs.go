package packagist

import (
	"time"
)

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
