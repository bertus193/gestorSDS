package client

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"runtime"
)

var httpClient *http.Client
var clear map[string]func() //create a map for storing clear funcs

var logguedUserEmail string
var logguedUserPass string

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
		logguedUserEmail = inputUser
		logguedUserPass = inputPass
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

	fmt.Printf("# Página de usuario\n\n")
	fmt.Printf("------ Listado de cuentas ------\n\n")
	// Recuperamos las cuentas del usuarios
	cuentas := listarCuentas(httpClient, logguedUserEmail, logguedUserPass)

	if cuentas != nil {
		// Imprimimos los resultados
		for c := range cuentas {
			tempAccount := cuentas[c]
			fmt.Printf("[%s] -> (%s / %s)\n", c, tempAccount.User, tempAccount.Password)
		}
	} else {
		fmt.Printf("No tienes ninguna cuenta guardada\n")
	}

	fmt.Printf("\n--------------------------------\n\n")

	var inputSelectionStr string
	fmt.Println("1. Añadir cuenta de servicio")
	fmt.Println("2. Ver detalle cuenta de servicio (to-do)")
	fmt.Println("3. Modificar mis datos (to-do)")
	fmt.Println("4. Eliminar mi cuenta (to-do)")
	fmt.Println("0. Salir")

	if fromError != "" {
		fmt.Printf("\n* %s", fromError)
	}
	fmt.Printf("\nSeleccione una opción: ")
	fmt.Scanf("%s", &inputSelectionStr)

	switch {
	case inputSelectionStr == "1":
		uiAddAccount("")
	// case inputSelectionStr == "2":
	// uiDetailAccountMenu("")
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
	fmt.Printf("# Añadir cuenta de servicio\n\n")
	if fromError != "" {
		fmt.Printf("* %s\n\n", fromError)
	}

	fmt.Print("Tipo de cuenta (twitter, facebook, etc): ")
	fmt.Scanf("%s", &inputAccountType)
	fmt.Print("Usuarios: ")
	fmt.Scanf("%s", &inputAccountUser)
	fmt.Print("Contraseña: ")
	fmt.Scanf("%s", &inputAccountPass)

	crearCuenta(httpClient, logguedUserEmail, logguedUserPass, inputAccountType, inputAccountUser, inputAccountPass)
	uiUserMainMenu("")
}

func startUi(c *http.Client) {
	httpClient = c
	uiInicio("")
}
