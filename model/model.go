package model

type Usuario struct {
	MasterPassword string
	Accounts       map[string]Account
}

type Account struct {
	User     string
	Password string
}
