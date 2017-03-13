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
	fmt.Println("2. Crear usuario (to-do)")

	if fromError != "" {
		fmt.Printf("\n* %s", fromError)
	}
	fmt.Printf("\nSeleccione una opción: ")
	fmt.Scanf("%s", &inputSelectionStr)

	switch {
	case inputSelectionStr == "1":
		uiLoginMaster("")
	/*
		case inputSelectionStr == "2":
			uiRegistroMaster("")
	*/
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
	//uiMainMenu()
}

func uiRegistroMaster(fromError string) {

}

func uiUserMainMenu(fromError string) {
	clearScreen()

	var inputSelectionStr string
	fmt.Printf("# Menú de usuario\n\n")
	fmt.Println("1. Ver cuentas de servicio (legacy)")
	fmt.Println("0. Salir")

	if fromError != "" {
		fmt.Printf("\n* %s", fromError)
	}
	fmt.Printf("\nSeleccione una opción: ")
	fmt.Scanf("%s", &inputSelectionStr)

	switch {
	case inputSelectionStr == "1":
		// to-do
		listarCuentas(httpClient, logguedUserEmail, logguedUserPass)
	case inputSelectionStr == "0":
		os.Exit(0)
	default:
		uiUserMainMenu("La opción elegida no es correcta")
	}
}

func startUi(c *http.Client) {
	httpClient = c
	uiInicio("")
}
