package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/bertus193/gestorSDS/model"
	"github.com/bertus193/gestorSDS/server/database"
	"github.com/bertus193/gestorSDS/utils"
)

// función para escribir una respuesta del servidor
func response(w http.ResponseWriter, code int, payloadJSON string) {
	w.WriteHeader(code)
	fmt.Fprintf(w, payloadJSON)
}

// Añade un usuario a la BD
func registroUsuario(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	email := req.Form.Get("email")
	pass := req.Form.Get("pass")

	// Logs
	AddLog("registroUsuario: [" + email + ", " + pass + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	// Añadimos el usuario a la base de datos
	if err := database.CreateUser(email, pass); err != nil {

		// Si ha ocurrido un error al añadir el usuario, comprobamos
		// el error y respondemos con el código http adecuado
		switch err.Error() {
		case "user already exists":
			response(w, 409, "") // (409 - Conflict)
		default:
			response(w, 500, "") // (500 - Internal Server Error)
		}

	} else {
		// Si la inserción se ha realizado correctamente
		response(w, 201, "")
	}
}

// Comprueba si existe un usuario en la BD
func loginUsuario(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	email := req.Form.Get("email")
	passw := req.Form.Get("pass")

	// Logs
	AddLog("loginUsuario: [" + email + ", " + passw + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	// Añadimos el usuario a la base de datos
	if user, err := database.GetUser(email, passw); err != nil {

		// Si ha ocurrido un error al recuperar el usuario, comprobamos
		// el error y respondemos con el código http adecuado
		switch err.Error() {
		case "user not found":
			response(w, 404, "") // (404 - Not found)
		case "passwords do not match":
			response(w, 400, "") // (400 - Bad Request)
		default:
			response(w, 500, "") // (500 - Internal Server Error)
		}

	} else if user.A2FEnabled == true {
		// Si el usuario existe pero tiene A2F activado
		// Creamos la sesión con activación vía A2F
		token, a2fcode := CreateUserSession(email, true)
		// Enviamos el código de A2F por correo
		utils.Send2FACode(email, a2fcode)
		// Respondemos con el token e informando
		response(w, 250, token) // (250 - A2F required [custom])
	} else {
		// Si el usuario existe y no tiene A2F activado
		token, _ := CreateUserSession(email, false)
		response(w, 200, token)
	}
}

func desbloquearA2F(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	token := req.Form.Get("token")
	a2fcode := req.Form.Get("a2fcode")

	// Logs
	AddLog("desbloquearA2F: [" + token + ", " + a2fcode + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	if err := UnlockSessionWith2FA(token, a2fcode); err != nil {

		// Si ha ocurrido un error al recuperar el usuario, comprobamos
		// el error y respondemos con el código http adecuado
		switch err.Error() {
		case "session not found":
			response(w, 404, "") // (404 - Not found)
		case "2fa already resolved":
			response(w, 304, "") // (304 - Not Modified)
		case "2fa expired":
			response(w, 408, "") // (408 - Request Timeout)
		case "incorrect 2fa code":
			response(w, 400, "") // (400 - Bad Request)
		default:
			response(w, 500, "") // (500 - Internal Server Error)
		}
	} else {
		// La sesión se ha desbloqueado correctamente
		response(w, 200, "")
	}
}

// Recupera las cuentas de servicio de un usuario de la BD
func listarCuentas(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	token := req.Form.Get("token")

	// Logs
	AddLog("listarCuentas: [" + token + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	// Respondemos
	if email, errSession := GetUserFromSession(token); errSession != nil {
		// La sesión ha caducado o no es valida
		response(w, 401, "") // (401 - Unauthorized)
	} else if entries, errEntries := database.GetVaultEntries(email); errEntries != nil {
		response(w, 500, "") // (500 - Internal Server Error)
	} else if entriesJSON, errJSON := json.Marshal(entries); errJSON != nil {
		response(w, 500, "") // (500 - Internal Server Error)
	} else {
		// Devolvemos la información
		response(w, 200, string(entriesJSON))
	}
}

// Añade una cuenta de servicio a un usuario de la BD
func crearEntrada(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	token := req.Form.Get("token")
	nombreServicio := req.Form.Get("tituloEntrada")
	usuarioServicio := req.Form.Get("usuarioCuenta")
	passServicio := req.Form.Get("passwordCuenta")

	// Logs
	AddLog("crearCuenta: [" + token + ", " + nombreServicio + ", " + usuarioServicio + ", " + passServicio + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	// Respondemos
	if email, errSession := GetUserFromSession(token); errSession != nil {
		// La sesión ha caducado o no es valida
		response(w, 401, "") // (401 - Unauthorized)
	} else if errCreate := database.CreateAccountVaultEntry(email, nombreServicio, usuarioServicio, passServicio); errCreate != nil {

		// Si ha ocurrido un error al insetar, comprobamos
		// el error y respondemos con el código http adecuado
		switch errCreate.Error() {
		case "user not found":
			response(w, 404, "") // (404 - Not found)
		case "entry already exists":
			response(w, 409, "") // (409 - Conflict)
		default:
			response(w, 500, "") // (500 - Internal Server Error)
		}

	} else {
		// Devolvemos la información
		response(w, 201, "")
	}
}

// Modifica los datos de un usuario de la BD
func modificarUsuario(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	token := req.Form.Get("token")
	passAnterior := req.Form.Get("passAnterior")
	passNuevo := req.Form.Get("passNuevo")
	AddLog("modificarUsuario: [" + token + ", " + passAnterior + ", " + passNuevo + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	// Respondemos
	if _, err := GetUserFromSession(token); err != nil {
		// La sesión ha caducado o no es valida
		response(w, 401, "")
	} else {
		// La sesión sigue abierta, eliminamos en la BD y respondemos
		// to-do
		response(w, 501, "to-do")
	}
}

// Elimina un usuario
func eliminarUsuario(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	token := req.Form.Get("token")
	AddLog("eliminarUsuario: [" + token + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	// Respondemos
	if email, err := GetUserFromSession(token); err != nil {
		// La sesión ha caducado o no es valida
		response(w, 401, "")
	} else {
		// La sesión sigue abierta, eliminamos en la BD y respondemos
		database.DeleteUser(email)
		response(w, 200, "")
	}
}

// Modifica usuario de una cuenta (servicio)
func modificarCuenta(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos (servicio)
	token := req.Form.Get("token")
	nombreServicio := req.Form.Get("nombreServicio")
	usuarioServicio := req.Form.Get("usuarioServicio")
	passServicio := req.Form.Get("passServicio")
	AddLog("modificarCuenta: [" + token + ", " + nombreServicio + ", " + usuarioServicio + ", " + passServicio + " ]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	// Respondemos
	if email, err := GetUserFromSession(token); err != nil {
		// La sesión ha caducado o no es valida
		response(w, 401, "")
	} else {
		// La sesión sigue abierta, modificamos en la BD y respondemos
		database.SetAccount(email, nombreServicio, usuarioServicio, passServicio)
		response(w, 200, "")
	}
}

// Elimina una cuenta de servicio a un usuario de la BD
func eliminarCuenta(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	token := req.Form.Get("token")
	nombreServicio := req.Form.Get("nombreServicio")
	AddLog("eliminarCuenta: [" + token + ", " + nombreServicio + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	// Respondemos
	if email, err := GetUserFromSession(token); err != nil {
		// La sesión ha caducado o no es valida
		response(w, 401, "")
	} else {
		// La sesión sigue abierta, eliminamos en la BD y respondemos
		database.DeleteAccount(email, nombreServicio)
		response(w, 200, "")
	}
}

// Recupera los detalles de una cuenta de servicio a un usuario de la BD
func detallesCuenta(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	token := req.Form.Get("token")
	nombreServicio := req.Form.Get("nombreServicio")
	AddLog("detallesCuenta: [" + token + ", " + nombreServicio + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	// Respondemos
	if email, err := GetUserFromSession(token); err != nil {
		// La sesión ha caducado o no es valida
		response(w, 401, "")
	} else {
		// La sesión sigue abierta, devolvemos la información
		accountInfo := database.GetJSONAccountFromUser(email, nombreServicio)
		response(w, 200, accountInfo)
	}
}

func detallesUsuario(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	token := req.Form.Get("token")
	AddLog("detallesUsuario: [" + token + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	// Respondemos
	if email, err := GetUserFromSession(token); err != nil {
		// La sesión ha caducado o no es valida
		response(w, 401, "")
	} else {
		// La sesión sigue abierta, devolvemos la información
		user, _ := database.GetUserFromEmail(email)
		log.Println(user.A2FEnabled)

		details := model.DetallesUsuario{
			Email:      email,
			A2FEnabled: user.A2FEnabled,
			NumEntries: len(user.Vault),
		}
		jsonString, _ := json.Marshal(details)
		response(w, 200, string(jsonString))
	}
}

func activarA2F(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	token := req.Form.Get("token")
	AddLog("activarA2f: [" + token + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	// Respondemos
	if email, err := GetUserFromSession(token); err != nil {
		// La sesión ha caducado o no es valida
		response(w, 401, "")
	} else {
		// La sesión sigue abierta, devolvemos la información
		database.ToggleA2f(email, true)
		response(w, 200, "")
	}
}

func desactivarA2F(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	token := req.Form.Get("token")
	AddLog("desactivarA2f: [" + token + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	// Respondemos
	if email, err := GetUserFromSession(token); err != nil {
		// La sesión ha caducado o no es valida
		response(w, 401, "")
	} else {
		// La sesión sigue abierta, devolvemos la información
		database.ToggleA2f(email, false)
		response(w, 200, "")
	}
}
