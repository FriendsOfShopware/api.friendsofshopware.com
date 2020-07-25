package _struct

type Ratings []struct {
	ID     int `json:"id"`
	Status struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	} `json:"status"`
	StoreRatingID  int    `json:"storeRatingId"`
	AuthorName     string `json:"authorName"`
	Headline       string `json:"headline"`
	Text           string `json:"text"`
	CreationDate   int    `json:"creationDate"`
	LastChangeDate int    `json:"lastChangeDate"`
	Rating         struct {
		ID    int `json:"id"`
		Value int `json:"value"`
	} `json:"rating"`
	Functionality struct {
		ID    int `json:"id"`
		Value int `json:"value"`
	} `json:"functionality"`
	Usability struct {
		ID    int `json:"id"`
		Value int `json:"value"`
	} `json:"usability"`
	Documentation struct {
		ID    int `json:"id"`
		Value int `json:"value"`
	} `json:"documentation"`
	Support interface{}   `json:"support"`
	Replies []interface{} `json:"replies"`
	History []interface{} `json:"history"`
	Plugin  struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Infos []struct {
			ID     int `json:"id"`
			Locale struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"locale"`
			Name               string        `json:"name"`
			Description        string        `json:"description"`
			InstallationManual string        `json:"installationManual"`
			ShortDescription   string        `json:"shortDescription"`
			Highlights         string        `json:"highlights"`
			Features           string        `json:"features"`
			Tags               []interface{} `json:"tags"`
			Videos             []interface{} `json:"videos"`
			Faqs               []interface{} `json:"faqs"`
		} `json:"infos"`
		Support               bool `json:"support"`
		SupportOnlyCommercial bool `json:"supportOnlyCommercial"`
		Generation            struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"generation"`
	} `json:"plugin"`
	Legacy bool `json:"legacy"`
}
