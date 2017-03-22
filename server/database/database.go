package database

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bertus193/gestorSDS/model"
	"github.com/bertus193/gestorSDS/utils"
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

func init() {
	before()
	// todo: borrar
	// Datos de relleno para probar el acceso (BORRAR)
	/*AddUser("demoEmail", "demoMasterPass")
	AddAccountToUser("demoEmail", "facebook", "facebookUser", "facebookPass")
	AddAccountToUser("demoEmail", "twitter", "twitterUser", "twitterPass")
	AddUser("demoEmail2", "demoMasterPass2")*/

}

func before() {
	bytesEntrada, err := ioutil.ReadFile("./server/database/bd.txt")
	if err != nil {
		panic(0)
	}

	result := make(map[string]model.Usuario)

	if err := json.Unmarshal(bytesEntrada, &result); err != nil {
		panic("Error al leer fichero de entrada")
	}

	gestor = result
}

// After Persistencia Base de Datos
func After() {
	salida, err := os.Create("./server/database/bd.txt")
	if err != nil {
		panic(0)
	}

	// todo: comprobar y validar contraseña

	j, err := json.Marshal(gestor)

	if err != nil {
		fmt.Println(err)
	}

	usuarios := string(j)

	salida.Write([]byte(usuarios))
}

// AddUser añade un usuarios al sistema
func AddUser(email string, pass string) {
	// todo: comprobar si el usuario ya existe

	// fmt.Print(GetAll())
	salt, err1 := utils.GenerateRandomBytes(64)
	if err1 == nil {
		bytePass := []byte(pass)
		hashpass, _ := utils.DeriveKey(bytePass, salt)

		gestor[email] = model.Usuario{MasterPassword: string(hashpass), MasterPasswordSalt: string(salt), Accounts: make(map[string]model.Account)}
	}

	fmt.Print(GetAll())
}

// AddAccountToUser añade datos a un ya dado de alta
func AddAccountToUser(userEmail string, serviceName string, serviceUser string, servicePass string) {
	// todo: comprobar que el usuario existe antes de asignar

	gestor[userEmail].Accounts[serviceName] = model.Account{User: serviceUser, Password: servicePass}
}

func ExistsUser(userEmail string, userPass string) bool {
	user, ok := gestor[userEmail]
	if ok {
		if userPass == user.MasterPassword {
			return true
		}
	}
	return false
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
func GetAll() string {
	j, err := json.Marshal(gestor)

	if err != nil {
		fmt.Println(err)
	}

	return string(j)
}
