package _struct

type Sales []struct {
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
