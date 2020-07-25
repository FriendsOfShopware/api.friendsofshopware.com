package _struct

type Token struct {
	Token  string `json:"token"`
	Expire struct {
		Date         string `json:"date"`
		TimezoneType int    `json:"timezone_type"`
		Timezone     string `json:"timezone"`
	} `json:"expire"`
	UserAccountID int  `json:"userAccountId"`
	UserID        int  `json:"userId"`
	LegacyLogin   bool `json:"legacyLogin"`
}

type LoginRequest struct {
	Email    string `json:"shopwareId"`
	Password string `json:"password"`
}
