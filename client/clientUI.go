package client

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"

	"github.com/bertus193/gestorSDS/model"
	"github.com/bertus193/gestorSDS/utils"
)

var httpClient *http.Client
var clear map[string]func() //create a map for storing clear funcs

func checkErrors(errStr string) {
	switch errStr {
	case "errorSesion":
		uiLoginMaster("La sesión ha caducado")
	case "error":
		fmt.Printf("* Ha ocurrido un error\n")
	}
}

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		// Linux clear
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["darwin"] = func() {
		// Mac clear
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		// Windows clear
		cmd := exec.Command("cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func clearScreen() {
	// runtime.GOOS -> linux, windows
	value, ok := clear[runtime.GOOS]

	// if we defined a clear func for that platform:
	if ok {
		// we execute it
		value()
	} else {
		// unsupported platform
		fmt.Println("-----------------------------------------------------")
	}
}

func uiInicio(fromError string) {
	clearScreen()

	var inputSelectionStr string
	fmt.Printf("# Bienvenido\n\n")
	fmt.Println("1. Entrar")
	fmt.Println("2. Crear usuario")
	fmt.Println("0. Salir")

	if fromError != "" {
		fmt.Printf("\n* %s", fromError)
	}
	fmt.Printf("\nSeleccione una opción: ")
	fmt.Scanf("%s", &inputSelectionStr)

	switch {
	case inputSelectionStr == "1":
		uiLoginMaster("")
	case inputSelectionStr == "2":
		uiRegistroMaster("")
	case inputSelectionStr == "0":
		os.Exit(0)
	default:
		uiInicio("La opción elegida no es correcta")
	}
}

func uiLoginMaster(fromError string) {
	clearScreen()

	var inputUser string
	var inputPass string
	fmt.Printf("# Entrada de usuarios\n\n")
	if fromError != "" {
		fmt.Printf("* %s\n\n", fromError)
	}

	fmt.Print("Email: ")
	fmt.Scanf("%s", &inputUser)
	fmt.Print("Contraseña: ")
	fmt.Scanf("%s", &inputPass)

	userExists := loginUsuario(httpClient, inputUser, inputPass)
	if userExists == true {
		uiUserMainMenu("")
	} else {
		uiInicio("El usuario no existe")
	}
}

func uiRegistroMaster(fromError string) {
	clearScreen()

	var inputUser string
	var inputPass string
	fmt.Printf("# Registro de usuarios\n\n")
	if fromError != "" {
		fmt.Printf("* %s\n\n", fromError)
	}

	fmt.Print("Email: ")
	fmt.Scanf("%s", &inputUser)
	fmt.Print("Contraseña: ")
	fmt.Scanf("%s", &inputPass)

	registroUsuario(httpClient, inputUser, inputPass)
	uiInicio("")
}

func uiUserMainMenu(fromError string) {
	clearScreen()

	// Recuperamos las cuentas del usuarios
	cuentas, errStr := listarCuentas(httpClient)

	checkErrors(errStr)

	fmt.Printf("# Página de usuario\n\n")
	fmt.Printf("------ Listado de cuentas ------\n\n")

	if cuentas != nil && len(cuentas) != 0 {
		// Imprimimos los resultados
		for c := range cuentas {
			//tempAccount := cuentas[c]
			//tempPass := string(utils.Decrypt(utils.Decode64(tempAccount.Password), keyData))
			//fmt.Printf("[%s] -> (%s / %s)\n", c, tempAccount.User, tempPass)
			fmt.Printf("[%s] ", c)
		}
	} else {
		fmt.Printf("* No tienes ninguna cuenta guardada\n")
	}

	fmt.Printf("\n\n--------------------------------\n\n")

	var inputSelectionStr string

	fmt.Println("1. Añadir cuenta")
	fmt.Println("2. Ver detalle cuenta")
	fmt.Println("3. Eliminar mi usuario")
	fmt.Println("0. Salir")

	if fromError != "" {
		fmt.Printf("\n* %s", fromError)
	}
	fmt.Printf("\nSeleccione una opción: ")
	fmt.Scanf("%s", &inputSelectionStr)

	switch {
	case inputSelectionStr == "1":
		uiAddAccount("")
	case inputSelectionStr == "2":
		var inputAccountSelectionStr string
		fmt.Print("Elige la cuenta: ")
		fmt.Scanf("%s", &inputAccountSelectionStr)
		uiServiceMenu("", inputAccountSelectionStr)
	case inputSelectionStr == "3":
		uiDeleteUser("")
	case inputSelectionStr == "0":
		os.Exit(0)
	default:
		uiUserMainMenu("La opción elegida no es correcta")
	}
}

func uiAddAccount(fromError string) {
	clearScreen()

	var inputAccountType string
	var inputAccountUser string
	var inputAccountPass string
	fmt.Printf("# Añadir cuenta\n\n")
	if fromError != "" {
		fmt.Printf("* %s\n\n", fromError)
	}

	fmt.Print("Tipo de cuenta (twitter, facebook, etc): ")
	fmt.Scanf("%s", &inputAccountType)
	fmt.Print("Usuario: ")
	fmt.Scanf("%s", &inputAccountUser)
	fmt.Print("Contraseña: ")
	fmt.Scanf("%s", &inputAccountPass)

	crearCuenta(httpClient, inputAccountType, inputAccountUser, inputAccountPass)
	uiUserMainMenu("")
}

func startUI(c *http.Client) {
	httpClient = c
	uiInicio("")
}

func uiServiceMenu(fromError string, accountSelectionStr string) {
	clearScreen()

	accountDetails, errStr := detallesCuenta(httpClient, accountSelectionStr)
	checkErrors(errStr)

	fmt.Printf("# Detalles de cuenta [%s]\n\n", accountSelectionStr)
	fmt.Printf("--------------------------------\n\n")
	// Si los detalles de la cuenta están vacios
	if (model.Account{}) == accountDetails {
		// Volvemos al menú del usuario
		uiUserMainMenu("No existe la cuenta seleccionada")
	}
	// Desencriptamos la contraseña para mostrarla
	plainPass := string(utils.Decrypt(utils.Decode64(accountDetails.Password), keyData))
	// Mostramos los detalles de la cuenta
	fmt.Printf("[%s] -> (%s / %s)\n", accountSelectionStr, accountDetails.User, plainPass)
	fmt.Printf("\n--------------------------------\n\n")

	var inputSelectionStr string
	fmt.Println("1. Modificar usuario")
	fmt.Println("2. Borrar cuenta")
	fmt.Println("0. Volver")

	if fromError != "" {
		fmt.Printf("\n* %s", fromError)
	}
	fmt.Printf("\nSeleccione una opción: ")
	fmt.Scanf("%s", &inputSelectionStr)

	switch {
	case inputSelectionStr == "1":
		uiModifyAccount("", accountSelectionStr)
	case inputSelectionStr == "2":
		uiDeleteAccount("", accountSelectionStr)
	case inputSelectionStr == "0":
		uiUserMainMenu("")
	default:
		uiServiceMenu("La opción elegida no es correcta", accountSelectionStr)
	}
}

func uiModifyAccount(fromError string, nombreServicio string) {
	clearScreen()

	var inputAccountUser string
	var inputAccountPassword string
	fmt.Printf("# Editar usuario de cuenta\n\n")
	if fromError != "" {
		fmt.Printf("* %s\n\n", fromError)
	}

	fmt.Printf("Introduce el nombre de la cuenta  %s: ", nombreServicio)
	fmt.Scanf("%s", &inputAccountUser)

	fmt.Printf("Introduce la contraseña la cuenta %s: ", nombreServicio)
	fmt.Scanf("%s", &inputAccountPassword)

	modificarCuenta(httpClient, inputAccountUser, inputAccountPassword, nombreServicio)
	uiUserMainMenu("")
}

func uiDeleteAccount(fromError string, nombreServicio string) {

	var inputDecission string
	fmt.Print("¿Estás seguro? (si, no): ")
	fmt.Scanf("%s", &inputDecission)

	if inputDecission == "si" {
		eliminarCuenta(httpClient, nombreServicio)
	}

	uiUserMainMenu("")
}

func uiDeleteUser(fromError string) {

	var inputDecission string
	fmt.Print("¿Estás seguro? (si, no): ")
	fmt.Scanf("%s", &inputDecission)

	if inputDecission == "si" {
		_, _, errStr := eliminarUsuario(httpClient)
		checkErrors(errStr)
	}

	uiInicio("")
}

// Logout externo
func UIlogout() {
	os.Exit(0)
}
