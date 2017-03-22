package model

type Usuario struct {
	MasterPassword     string
	MasterPasswordSalt string
	Accounts           map[string]Account
}

type Account struct {
	User     string
	Password string
}
