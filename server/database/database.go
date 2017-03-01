package database

import (
	"encoding/json"
	"fmt"
)

type usuario struct {
	Password string
	Datos    map[string]string
}

/* Demo estructura

"alu@alu.ua.es" : {
    password: "accoutPass"
    data: [
        "facebook": "12345"
        "twitter": "abcd"
    ]
}
*/

var gestor = make(map[string]usuario)

// AddUser añade un usuarios al sistema
func AddUser(email string, pass string) {
	// todo: comprobar si el usuario ya existe

	gestor[email] = usuario{Password: pass, Datos: make(map[string]string)}
}

// AddToUser añade datos a un ya dado de alta
func AddToUser(email string, service string, servicePass string) {
	// todo: comprobar que el usuario existe antes de asignar

	gestor[email].Datos[service] = servicePass
}

// GetAll (Debug) Devuelve un string json con todos los datos
func GetAll() string {
	j, err := json.Marshal(gestor)

	if err != nil {
		fmt.Println(err)
	}

	return string(j)
}
