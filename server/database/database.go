package database

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/bertus193/gestorSDS/config"
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

var gestor = make(map[string]*model.Usuario)

func init() {
	before()
}

// Lee el fichero
func before() {
	result := make(map[string]*model.Usuario)

	bytesEntrada, err := ioutil.ReadFile("./server/database/bd.txt")
	error := false
	if err != nil {
		error = true
	}
	if error == true || len(string(bytesEntrada)) == 0 {
		//fileData := []byte("{}")
		ioutil.WriteFile("./server/database/bd.txt", []byte(""), 0644)
	} else {
		decompress := []byte(utils.Decompress(bytesEntrada))
		decompress = utils.Decrypt(decompress, config.PassDBEncrypt)
		if err := json.Unmarshal(decompress, &result); err != nil {
			panic("Error al leer fichero de entrada")
		}
	}

	gestor = result
}

// CreateUser guarda un nuevo usuario en la BD
func CreateUser(email string, passw string) error {

	var errResult error

	// Comprobamos si existe el email en la BD
	if _, ok := gestor[email]; ok {
		// Si existe el email, no modificamos nada
		errResult = errors.New("user already exists")
	} else if salt, errSalt := utils.GenerateRandomBytes(64); errSalt != nil {
		// Error al generar "salt"
		errResult = errors.New("unable to save")
	} else {
		// Hash de la contraseña también en servidor
		bytePass := []byte(passw)
		hashPass, _ := utils.DeriveKey(bytePass, salt)
		saltBase64 := base64.StdEncoding.EncodeToString(salt)

		// Guardamos el nuevo usuario
		gestor[email] = &model.Usuario{
			UserPassword:     string(hashPass),
			UserPasswordSalt: saltBase64,
			A2FEnabled:       false,
			Vault:            make(map[string]model.VaultEntry)}
	}
	return errResult
}

// GetUser recupera un usuario de la BD que contenta el mismo
// email y contraseña que las indicads
func GetUser(email string, passw string) (*model.Usuario, error) {

	var userResult *model.Usuario
	var errResult error

	// Comprobamos si existe el email en la BD
	if user, ok := gestor[email]; !ok {
		// Si no existe el el usuario indicado
		errResult = errors.New("user not found")
	} else if salt, errSalt := base64.StdEncoding.DecodeString(user.UserPasswordSalt); errSalt != nil {
		// Error al recuperar el "salt"
		errResult = errors.New("unable to recover")
	} else {
		// Regeneramos el hash de servidor de la contraseña
		bytePass := []byte(passw)
		if hashPass, errHash := utils.DeriveKey(bytePass, salt); errHash != nil {
			// Error al regenerar el hash
			errResult = errors.New("unable to recover")
		} else if user.UserPassword != string(hashPass) {
			// Las contraseñas no coinciden
			errResult = errors.New("passwords do not match")
		} else {
			userResult = user
		}
	}
	return userResult, errResult
}

// GetEntries recupera la lista de entradas (sin detalles)
// de un usuario
func GetVaultEntries(email string) ([]string, error) {

	var entriesResult []string
	var errResult error

	// Comprobamos si existe el email en la BD
	if user, ok := gestor[email]; !ok {
		// Si no existe el el usuario indicado
		errResult = errors.New("user not found")
	} else {
		// Recuperamos solo el "título" de las entradas
		entriesResult = make([]string, len(user.Vault))
		for entry := range user.Vault {
			entriesResult = append(entriesResult, entry)
		}
	}
	return entriesResult, errResult
}

func CreateAccountVaultEntry(email string, entryTitle string, userAccount string, passwAccount string) error {
	var errResult error

	if user, okUser := gestor[email]; !okUser {
		// Si no existe el el usuario indicado, no modificamos nada
		errResult = errors.New("user not found")
	} else if _, okEntry := user.Vault[entryTitle]; okEntry {
		// Si ya existe una entrada con el mismo título
		errResult = errors.New("entry already exists")
	} else {
		user.Vault[entryTitle] = model.VaultEntry{
			Mode:     1, // Account
			User:     userAccount,
			Password: passwAccount,
		}
	}

	return errResult
}

/**/
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

	//usuarios := string(j)
	usuarios := string(utils.Encrypt(j, config.PassDBEncrypt)) //Encriptar
	usuarios = utils.Compress(usuarios)                        //Comprimir

	salida.Write([]byte(usuarios))
}

func GetUserFromEmail(userEmail string) (*model.Usuario, error) {
	var err error
	user, ok := gestor[userEmail]
	if !ok {
		err = errors.New("user not found")
	}
	return user, err
}

// AddAccountToUser añade datos a un ya dado de alta
func AddAccountToUser(userEmail string, serviceName string, serviceUser string, servicePass string) {
	// todo: comprobar que el usuario existe antes de asignar

	gestor[userEmail].Vault[serviceName] = model.VaultEntry{User: serviceUser, Password: servicePass}
}

// GetJSONAllAccountsFromUser listado de cuentas asignadas a un usuario
func GetJSONAccountFromUser(usuario string, nombreServicio string) string {
	userAccount := gestor[usuario].Vault[nombreServicio]
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

// SetAccountUser Modifica cuenta de usuario
func SetAccount(userEmail string, serviceName string, serviceUser string, servicePass string) {
	// todo: comprobar que el usuario existe antes de asignar
	gestor[userEmail].Vault[serviceName] = model.VaultEntry{User: serviceUser, Password: servicePass}
}

// deleteAccount Elimina cuenta de usuario
func DeleteAccount(userEmail string, serviceName string) {
	// todo: comprobar que el usuario existe antes de asignar
	delete(gestor[userEmail].Vault, serviceName)
}

// deleteAccount Elimina cuenta de usuario
func DeleteUser(userEmail string) {
	// todo: comprobar que el usuario existe antes de asignar
	delete(gestor, userEmail)
}

func ToggleA2f(userEmail string, status bool) {
	gestor[userEmail].A2FEnabled = status
}
