package model

/* ----------- DATABASE ----------- */

type Usuario struct {
	UserPassword     string
	UserPasswordSalt string
	A2FEnabled       bool
	Vault            map[string]VaultEntry
}

type VaultEntry struct {
	Title string
	Mode  int
	// Mode 0 - Plain text
	Text string
	// Mode 1 - Account
	User     string
	Password string
	// Mode 2 - Credit card
	CreditCard string
}

/* -------------------------------- */

/*  ----- DETALLES DE USUARIO ----- */

type DetallesUsuario struct {
	Email      string
	A2FEnabled bool
	NumEntries int
}

/* -------------------------------- */
