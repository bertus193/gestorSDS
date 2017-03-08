package database

import (
	"encoding/json"
	"fmt"

	"github.com/bertus193/gestorSDS/model"
)

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

var gestor = make(map[string]model.Usuario)

// AddUser añade un usuarios al sistema
func AddUser(email string, pass string) {
	// todo: comprobar si el usuario ya existe

	gestor[email] = model.Usuario{MasterPassword: pass, Accounts: make(map[string]model.Account)}
}

// AddAccountToUser añade datos a un ya dado de alta
func AddAccountToUser(userEmail string, serviceName string, serviceUser string, servicePass string) {
	// todo: comprobar que el usuario existe antes de asignar

	gestor[userEmail].Accounts[serviceName] = model.Account{User: serviceUser, Password: servicePass}
}

// GetJSONAllAccountsFromUser listado de cuentas asignadas a un usuario
func GetJSONAllAccountsFromUser(usuario string, pass string) string {
	userAccounts := gestor[usuario].Accounts
	// todo: comprobar y validar contraseña

	j, err := json.Marshal(userAccounts)

	if err != nil {
		fmt.Println(err)
	}

	return string(j)
}

// GetAll (Debug) Devuelve un string json con todos los datos
/*func GetAll() string {
	j, err := json.Marshal(gestor)

	if err != nil {
		fmt.Println(err)
	}

	return string(j)
}*/
