package model

type Usuario struct {
	MasterPassword     string
	MasterPasswordSalt string
	A2FEnabled         bool
	Accounts           map[string]Account
}

type Account struct {
	User     string
	Password string
}
