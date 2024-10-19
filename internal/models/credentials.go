package models

type Credentials struct {
	CredentialsName string `json:"credentials_name"`
	Login           string `json:"login"`
	Password        string `json:"password"`
}
