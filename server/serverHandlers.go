package server

import (
	"fmt"
	"net/http"
	"time"

	"github.com/bertus193/gestorSDS/server/database"
	"github.com/bertus193/gestorSDS/utils"
)

// función para escribir una respuesta del servidor
func response(w http.ResponseWriter, code int, payloadJSON string) {
	w.WriteHeader(code)
	fmt.Fprintf(w, payloadJSON)
}

var session = make(map[string]time.Time)

// Comprueba si existe un usuario en la BD
func loginUsuario(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	email := req.Form.Get("email")
	passw := req.Form.Get("pass")

	AddLog("loginUsuario: [" + email + ", " + passw + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	if database.ExistsUser(email, passw) {
		if tempUser, err := database.GetUserFromEmail(email); err != nil {
			response(w, 500, "")
		} else if tempUser.A2FEnabled == true {
			token, a2fcode := CreateUserSession(email, true)
			utils.Send2FACode(email, a2fcode)
			response(w, 250, token)
		} else {
			token, _ := CreateUserSession(email, false)
			response(w, 200, token)
		}
	} else {
		response(w, 500, "")
	}
}

// Añade un usuario a la BD
func registroUsuario(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	email := req.Form.Get("email")
	pass := req.Form.Get("pass")
	AddLog("registroUsuario: [" + email + ", " + pass + "]")

	database.AddUser(email, pass)

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")
	// Respondemos
	response(w, 201, "")
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

// Añade una cuenta de servicio a un usuario de la BD
func crearCuenta(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	token := req.Form.Get("token")
	nombreServicio := req.Form.Get("nombreServicio")
	usuarioServicio := req.Form.Get("usuarioServicio")
	passServicio := req.Form.Get("passServicio")
	AddLog("crearCuenta: [" + token + ", " + nombreServicio + ", " + usuarioServicio + ", " + passServicio + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	// Respondemos
	if email, err := GetUserFromSession(token); err != nil {
		// La sesión ha caducado o no es valida
		response(w, 401, "")
	} else {
		// La sesión sigue abierta, añadimos a la BD y respondemos
		database.AddAccountToUser(email, nombreServicio, usuarioServicio, passServicio)
		response(w, 201, "")
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

// Recupera las cuentas de servicio de un usuario de la BD
func listarCuentas(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	token := req.Form.Get("token")
	AddLog("listarCuentas: [" + token + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	// Respondemos
	if email, err := GetUserFromSession(token); err != nil {
		// La sesión ha caducado o no es valida
		response(w, 401, "")
	} else {
		// La sesión sigue abierta, devolvemos la información
		response(w, 200, database.GetJSONAllAccountsFromUser(email))
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
