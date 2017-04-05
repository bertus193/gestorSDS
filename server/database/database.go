package database

import (
	"encoding/base64"
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

// Lee el fichero
func before() {
	result := make(map[string]model.Usuario)

	bytesEntrada, err := ioutil.ReadFile("./server/database/bd.txt")
	error := false
	if err != nil {
		error = true
	}
	if error == true || len(string(bytesEntrada)) == 0 {
		//fileData := []byte("{}")
		ioutil.WriteFile("./server/database/bd.txt", []byte(""), 0644)
	} else {
		if err := json.Unmarshal(bytesEntrada, &result); err != nil {
			panic("Error al leer fichero de entrada")
		}
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

	salt, errSalt := utils.GenerateRandomBytes(64)
	if errSalt == nil {
		bytePass := []byte(pass)
		hashPass, _ := utils.DeriveKey(bytePass, salt)
		saltBase64 := base64.StdEncoding.EncodeToString(salt)

		gestor[email] = model.Usuario{
			MasterPassword:     string(hashPass),
			MasterPasswordSalt: saltBase64,
			Accounts:           make(map[string]model.Account)}
	}
}

// ExistsUser Comprueba que el usuario existe en la BD
func ExistsUser(userEmail string, userPass string) bool {
	user, ok := gestor[userEmail]
	if ok {
		// Comprobamos la contraseña
		// Recuperamos el salt del usuario
		salt, _ := base64.StdEncoding.DecodeString(user.MasterPasswordSalt)
		bytePass := []byte(userPass)
		// Regeneramos el hash
		hashPass, _ := utils.DeriveKey(bytePass, salt)

		// Comprobamos que sean iguales
		if user.MasterPassword == string(hashPass) {
			return true
		}
	}
	return false
}

// AddAccountToUser añade datos a un ya dado de alta
func AddAccountToUser(userEmail string, userPass string, serviceName string, serviceUser string, servicePass string) {
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

// GetJSONAllAccountsFromUser listado de cuentas asignadas a un usuario
func GetJSONAccountFromUser(usuario string, pass string, nombreServicio string) string {
	userAccount := gestor[usuario].Accounts[nombreServicio]
	// todo: comprobar y validar contraseña

	j, err := json.Marshal(userAccount)

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
