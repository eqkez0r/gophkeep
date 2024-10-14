package card

type Card struct {
	CardName       string `json:"card_name"`
	CardNumber     string `json:"card_number"`
	CardHolderName string `json:"card_holder_name"`
	ExpirationDate string `json:"expiration_date"`
	CVV            int32  `json:"cvv"`
}
