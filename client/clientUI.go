package client

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/bertus193/gestorSDS/config"
	"github.com/bertus193/gestorSDS/utils"
	"github.com/fatih/color"
)

var httpClient *http.Client

/*func checkErrors(errStr string) {
	switch errStr {
	case "errorSesion":
		uiLoginMaster("La sesión ha caducado")
	case "error":
		fmt.Printf("* Ha ocurrido un error\n")
	}
}*/

func startUI(c *http.Client) {
	httpClient = c
	uiInicio("")
}

// Pantalla de bienvenida con las opciones de
// login, registro y cerrar aplicación
func uiInicio(fromError string) {

	// Limpiamos la pantalla
	utils.ClearScreen()

	// Título de la pantalla
	fmt.Printf("# Bienvenido a %s\n\n", config.AppName)

	// Opciones
	fmt.Println("1. Entrar")
	fmt.Println("2. Crear usuario")
	fmt.Println("0. Salir")

	// Mensaje de error en caso de existir
	if fromError != "" {
		color.Red("\n* %s", fromError)
	}

	// Lectura de opción elegida
	fmt.Printf("\nSeleccione una opción: ")
	inputSelectionStr := utils.CustomScanf()

	// Ejecución de la opción elegida
	switch {
	case inputSelectionStr == "1": // Login
		uiLoginUser("")
	case inputSelectionStr == "2": // Registro
		uiRegistroUsuario("")
	case inputSelectionStr == "0": // Salir
		os.Exit(0)
	default:
		uiInicio("La opción elegida no es correcta")
	}
}

// Pantalla de creación de usuario
func uiRegistroUsuario(fromError string) {

	// Limpiamos la pantalla
	utils.ClearScreen()

	// Título de la pantalla
	fmt.Printf("# Registro de usuarios\n\n")

	// Mensaje de error en caso de existir
	if fromError != "" {
		fmt.Printf("* %s\n\n", fromError)
	}

	// Lectura de datos del nuevo usuario
	fmt.Print("Email: ")
	inputUser := utils.CustomScanf()
	fmt.Print("Contraseña: ")
	inputPass := utils.GetPassw()

	// Petición al servidor
	if err := registroUsuario(httpClient, inputUser, inputPass); err != nil {
		// Si hay un error, mostramos el mensaje de error adecuado
		switch err.Error() {
		case "user already exists":
			uiRegistroUsuario("Ya existe un usuario con ese correo.")
		default:
			uiRegistroUsuario("Ocurrio un error al crear el usuario.")
		}
	} else {
		// Registro completado, volvemos a la pantalla de inicio
		uiInicio("")
	}
}

// Pantalla de entrada de usuarios
func uiLoginUser(fromError string) {

	// Limpiamos la pantalla
	utils.ClearScreen()

	// Título de la pantalla
	fmt.Printf("# Entrada de usuarios\n\n")

	// Mensaje de error en caso de existir
	if fromError != "" {
		fmt.Printf("* %s\n\n", fromError)
	}

	// Lectura de datos del usuario
	fmt.Print("Email: ")
	inputUser := utils.CustomScanf()
	fmt.Print("Contraseña: ")
	inputPass := utils.GetPassw()

	// Petición al servidor
	if err := loginUsuario(httpClient, inputUser, inputPass); err != nil {

		// Si hay un error, mostramos el mensaje de error adecuado
		switch err.Error() {
		case "user not found":
			uiLoginUser("No exite ningún usuario con esos datos.")
		case "passwords do not match":
			// No damos información detallada del error en este caso
			uiLoginUser("No exite ningún usuario con esos datos.")
		case "a2f required":
			uiUnlockA2F("")
		default:
			uiLoginUser("Ocurrio un error al realizar el login.")
		}

	} else {
		// Login completado, vamos a la pantalla principal del usuario
		uiUserMainMenu("")
	}
}

func uiUnlockA2F(fromError string) {
	// Limpiamos la pantalla
	utils.ClearScreen()

	// Título de la pantalla
	fmt.Printf("# Verificación en 2 pasos\n\n")

	// Mensaje de información
	fmt.Printf("Te hemos enviado un correo a tu cuenta con el código que debes introducir para iniciar sesión.\n\n")

	// Mensaje de error en caso de existir
	if fromError != "" {
		fmt.Printf("* %s\n\n", fromError)
	}

	// Lectura de opción elegida
	fmt.Print("Código: ")
	inputA2Fcode := utils.CustomScanf()

	// Petición al servidor
	if err := desbloquearA2F(httpClient, inputA2Fcode); err != nil {
		// Si hay un error, mostramos el mensaje de error adecuado
		switch err.Error() {
		case "session not found":
			uiLoginUser("La sesión de usuario ha cadudado.")
		case "2fa already resolved":
			uiUserMainMenu("")
		case "2fa expired":
			uiLoginUser("El código de verificación en dos pasos ha caducado.")
		case "incorrect 2fa code":
			uiUnlockA2F("El código introducido no es valido.")
		default:
			uiUnlockA2F("Ocurrio un error verificar el código.")
		}
	} else {
		uiUserMainMenu("")
	}
}

func uiUserMainMenu(fromError string) {

	// Limpiamos la pantalla
	utils.ClearScreen()

	// Título de la pantalla
	fmt.Printf("# Página de usuario\n\n")

	// Recuperamos las cuentas del usuario
	fmt.Printf("------ Listado de cuentas ------\n\n")
	// Petición al servidor
	entradas, err := listarCuentas(httpClient)
	if err != nil {
		// Si hay un error, mostramos el mensaje de error adecuado
		switch err.Error() {
		case "unauthorized":
			uiLoginUser("La sesión de usuario ha cadudado.")
		default:
			fmt.Println("Ocurrio un error al recuperar las entradas." + err.Error())
		}
	} else {
		if entradas != nil && len(entradas) != 0 {
			// Imprimimos los resultados
			for c := range entradas {
				fmt.Printf("[%s]\n", entradas[c])
			}
		} else {
			fmt.Printf("* No hay nada guardado todavía\n")
		}
	}
	fmt.Printf("\n--------------------------------\n\n")

	// Opciones
	fmt.Println("1. Añadir cuenta")
	fmt.Println("2. Ver detalle cuenta")
	fmt.Println("3. Configuración de mi cuenta")
	fmt.Println("0. Salir")

	// Mensaje de error en caso de existir
	if fromError != "" {
		fmt.Printf("\n* %s", fromError)
	}

	// Lectura de opción elegida
	fmt.Printf("\nSeleccione una opción: ")
	inputSelectionStr := utils.CustomScanf()

	// Ejecución de la opción elegida
	switch {
	case inputSelectionStr == "1":
		uiAddNewEntry("")
	case inputSelectionStr == "2":
		// var inputAccountSelectionStr string
		// fmt.Print("Elige la cuenta: ")
		// inputAccountSelectionStr = utils.CustomScanf()
		// uiServiceMenu("", inputAccountSelectionStr)
	case inputSelectionStr == "3":
		uiUserConfiguration("")
	case inputSelectionStr == "0":
		os.Exit(0)
	default:
		uiUserMainMenu("La opción elegida no es correcta")
	}
}

func uiAddNewEntry(fromError string) {

	// Limpiamos la pantalla
	utils.ClearScreen()

	// Título de la pantalla
	fmt.Printf("# Añadir nueva entrada\n\n")

	// Mensaje de error en caso de existir
	if fromError != "" {
		fmt.Printf("* %s\n\n", fromError)
	}

	// Lectura de los datos de la nueva entrada
	fmt.Print("Título de la entrada (twitter, facebook, etc): ")
	inputAccountType := utils.CustomScanf()
	fmt.Println(inputAccountType)
	fmt.Print("Usuario: ")
	inputAccountUser := utils.CustomScanf()
	fmt.Println(inputAccountUser)
	fmt.Print("¿Deseas generar una contraseña? (si, no): ")
	inputGeneratePassw := utils.CustomScanf()
	var finalPassw string
	if inputGeneratePassw == "si" || inputGeneratePassw == "s" {
		// Solicitamos infromación de como se desea generar la contraseña

		// Tamaño de la contraseña
		var genLenght int
		for {
			fmt.Print("¿Que tamaño de contraseña deseas? ")
			inputLenght := utils.CustomScanf()
			if convLenght, err := strconv.Atoi(inputLenght); err == nil {
				genLenght = convLenght
				break
			}
		}

		// La contraseña generada puede tener números
		fmt.Print("¿Deseas que tenga números? (si, no): ")
		inputWithNums := utils.CustomScanf()
		genWithNums := inputWithNums == "si" || inputWithNums == "s"

		// La contraseña generada puede tener simbolos
		fmt.Print("¿Deseas que tenga símbolos? (si, no): ")
		inputWithSymbols := utils.CustomScanf()
		genWithSymbols := inputWithSymbols == "si" || inputWithSymbols == "s"

		finalPassw = utils.GeneratePassword(genLenght, true, genWithNums, genWithSymbols)

	} else {
		fmt.Print("Contraseña: ")
		finalPassw = utils.GetPassw()
	}

	// Petición al servidor
	if err := crearEntrada(httpClient, inputAccountType, inputAccountUser, finalPassw); err != nil {
		// Si hay un error, mostramos el mensaje de error adecuado
		switch err.Error() {
		case "unauthorized":
			uiLoginUser("La sesión de usuario ha cadudado.")
		case "user not found":
			uiLoginUser("Ha ocurrido un error al guardar la entrada en tu cuenta.")
		case "2fa expired":
			uiAddNewEntry("Ya existe una entrada con ese título.")
		default:
			uiUserMainMenu("Ocurrio un error al añadir la entrada el código.")
		}
	} else {
		uiUserMainMenu("")
	}
}

func uiUserConfiguration(fromError string) {

	// Limpiamos la pantalla
	utils.ClearScreen()

	// Título de la pantalla
	fmt.Printf("# Página de configuración de usuario\n\n")

	userDetails, _ := detallesUsuario(httpClient)
	fmt.Printf("------ Información de la cuenta ------\n\n")
	fmt.Printf("Correo electrónico: %s\n", userDetails.Email)
	fmt.Printf("Número de cuentas guardadas: %d\n", userDetails.NumEntries)
	fmt.Printf("Segundo factor de autenticación: ")
	if userDetails.A2FEnabled {
		fmt.Println("Activado")
	} else {
		fmt.Println("Desactivado")
	}
	fmt.Printf("\n\n----------------------------------\n\n")

	var inputSelectionStr string

	fmt.Println("1. Modificar contraseña")
	fmt.Println("2. Eliminar mi usuario")
	if userDetails.A2FEnabled {
		fmt.Println("3. Desactivar 2FA")
	} else {
		fmt.Println("3. Activar 2FA")
	}

	fmt.Println("0. Volver")

	if fromError != "" {
		fmt.Printf("\n* %s", fromError)
	}
	fmt.Printf("\nSeleccione una opción: ")
	inputSelectionStr = utils.CustomScanf()

	switch {
	case inputSelectionStr == "1":
		uiUserConfiguration("")
	case inputSelectionStr == "2":
		uiUserConfiguration("")
	case inputSelectionStr == "3":
		toggleA2f(httpClient, !userDetails.A2FEnabled)
		uiUserConfiguration("")
	case inputSelectionStr == "0":
		uiUserMainMenu("")
	default:
		uiUserConfiguration("* La opción elegida no es correcta")
	}
}

/*




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
	inputSelectionStr = utils.CustomScanf()

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
	inputAccountUser = utils.CustomScanf()

	fmt.Printf("Introduce la contraseña la cuenta %s: ", nombreServicio)
	inputAccountPassword = utils.CustomScanf()

	modificarCuenta(httpClient, inputAccountUser, inputAccountPassword, nombreServicio)
	uiUserMainMenu("")
}

func uiDeleteAccount(fromError string, nombreServicio string) {

	var inputDecission string
	fmt.Print("¿Estás seguro? (si, no): ")
	inputDecission = utils.CustomScanf()

	if inputDecission == "si" {
		eliminarCuenta(httpClient, nombreServicio)
	}

	uiUserMainMenu("")
}

func uiDeleteUser(fromError string) {

	var inputDecission string
	fmt.Print("¿Estás seguro? (si, no): ")
	inputDecission = utils.CustomScanf()

	if inputDecission == "si" {
		_, _, errStr := eliminarUsuario(httpClient)
		checkErrors(errStr)
	}

	uiInicio("")
}*/

// Logout externo
func UIlogout() {
	os.Exit(0)
}
