package database

import (
	"encoding/json"
	"fmt"
)

/* Alternativa con un struct
type usuarioEmail struct {
	PasswordMaestra string
	Datos           map[string][string]string
}
*/

type usuario struct {
	MasterPassword string
	Accounts       map[string]account
}

type account struct {
	User     string
	Password string
}

/* Demo estructura en json
"alu@alu.ua.es" : {
    "MasterPassword": "accoutPass",
    "Accounts": [
        "facebook": {
			"User": "usuarioFacebook"
			"Password": "12345"
		},
        "twitter": {
			"User": "usuarioTwitter"
			"Password": "54321"
		}
    ]
}
*/

var gestor = make(map[string]usuario)

// AddUser añade un usuarios al sistema
func AddUser(email string, pass string) {
	// todo: comprobar si el usuario ya existe

	gestor[email] = usuario{MasterPassword: pass, Accounts: make(map[string]account)}
}

// AddToUser añade datos a un ya dado de alta
func AddAccountToUser(userEmail string, serviceName string, serviceUser string, servicePass string) {
	// todo: comprobar que el usuario existe antes de asignar

	gestor[userEmail].Accounts[serviceName] = account{User: serviceUser, Password: servicePass}
}

// GetAll (Debug) Devuelve un string json con todos los datos
func GetAll() string {
	j, err := json.Marshal(gestor)

	if err != nil {
		fmt.Println(err)
	}

	return string(j)
}
