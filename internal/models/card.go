package models

type Card struct {
	CardName       string `json:"card_name"`
	CardNumber     string `json:"card_number"`
	CardHolderName string `json:"card_holder_name"`
	ExpirationDate string `json:"expiration_date"`
	CVV            string `json:"cvv"`
}
