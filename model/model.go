package model

/* ----------- DATABASE ----------- */

type Usuario struct {
	UserPassword     string
	UserPasswordSalt string
	A2FEnabled       bool
	Vault            map[string]VaultEntry
}

type VaultEntry struct {
	Mode int
	// Mode 0 - Plain text
	Text string
	// Mode 1 - Account
	User     string
	Password string
}

/* Demo estructura en json (sin cifrados)
"alu@alu.ua.es" : {
    "UserPassword": 	"accoutPass",
    "UserPasswordSalt": "accoutSalrPass",
    "Vault": [
        "memoria": {
			"Mode": "0""
			"Text": "texto de la entrada"
		},
        "twitter": {
			"Mode": "1"
			"User": "usuarioTwitter"
			"Password": "54321"
		}
    ]
}
*/

/* -------------------------------- */

/*  ----- DETALLES DE USUARIO ----- */

type DetallesUsuario struct {
	Email      string
	A2FEnabled bool
	NumEntries int
}

type ListaEntradas struct {
	Texts    []string
	Accounts []string
}

/* -------------------------------- */
