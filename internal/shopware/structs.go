package shopware

type Ratings []Rating

type Rating struct {
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

type Sales []Sale

type Sale struct {
	CreationDate string `json:"creationDate"`
	OrderNumber  string `json:"orderNumber"`
	Price        int    `json:"price"`
	LicenseShop  struct {
		Company struct {
			CustomerNumber string `json:"customerNumber"`
			Name           string `json:"name"`
		} `json:"company"`
		Domain string `json:"domain"`
	} `json:"licenseShop"`
	Plugin struct {
		ID    int    `json:"id"`
		Name  string `json:"name"`
		Infos []struct {
			ID     int `json:"id"`
			Locale struct {
				ID   int    `json:"id"`
				Name string `json:"name"`
			} `json:"locale"`
			Name string `json:"name"`
		} `json:"infos"`
		Generation struct {
			ID          int    `json:"id"`
			Name        string `json:"name"`
			Description string `json:"description"`
		} `json:"generation"`
	} `json:"plugin"`
	ExpirationDate interface{} `json:"expirationDate"`
	Charging       struct {
		LastBookingDate interface{} `json:"lastBookingDate"`
		NextBookingDate interface{} `json:"nextBookingDate"`
	} `json:"charging"`
	Subscription struct {
		CreationDate   interface{} `json:"creationDate"`
		ExpirationDate interface{} `json:"expirationDate"`
	} `json:"subscription"`
	VariantType struct {
		Name string `json:"name"`
	} `json:"variantType"`
}

type Token struct {
	Token         string          `json:"token"`
	Expire        TokenExpiration `json:"expire"`
	UserAccountID int             `json:"userAccountId"`
	UserID        int             `json:"userId"`
	LegacyLogin   bool            `json:"legacyLogin"`
}

type TokenExpiration struct {
	Date         string `json:"date"`
	TimezoneType int    `json:"timezone_type"`
	Timezone     string `json:"timezone"`
}

type LoginRequest struct {
	Email    string `json:"shopwareId"`
	Password string `json:"password"`
}
