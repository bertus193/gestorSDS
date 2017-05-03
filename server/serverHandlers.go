package server

import (
	"log"
	"net/http"
	"time"

	"github.com/bertus193/gestorSDS/server/database"
)

var session = make(map[string]time.Time)

// Comprueba si existe un usuario en la BD
func loginUsuario(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	email := req.Form.Get("email")
	updateSession(email)
	pass := req.Form.Get("pass")

	log.Println("loginUsuario: [" + email + ", " + pass + "]")

	userExists := database.ExistsUser(email, pass)

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")
	// Respondemos
	if userExists {
		response(w, false, 200, "")
	} else {
		response(w, false, 500, "")
	}
}

// Añade un usuario a la BD
func registroUsuario(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	email := req.Form.Get("email")
	updateSession(email)
	pass := req.Form.Get("pass")
	log.Println("registroUsuario: [" + email + ", " + pass + "]")

	database.AddUser(email, pass)

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")
	// Respondemos
	response(w, false, 201, "")
}

// Modifica los datos de un usuario de la BD
func modificarUsuario(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	email := req.Form.Get("email")
	updateSession(email)
	passAnterior := req.Form.Get("passAnterior")
	passNuevo := req.Form.Get("passNuevo")
	log.Println("modificarUsuario: [" + email + ", " + passAnterior + ", " + passNuevo + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")
	// Respondemos
	response(w, false, 501, "to-do")
}

// Elimina un usuario
func eliminarUsuario(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	email := req.Form.Get("email")
	updateSession(email)
	pass := req.Form.Get("pass")
	log.Println("eliminarUsuario: [" + email + ", " + pass + "]")

	database.DeleteUser(email, pass)
	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")
	// Respondemos
	response(w, false, 200, "")
}

// Añade una cuenta de servicio a un usuario de la BD
func crearCuenta(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	email := req.Form.Get("email")
	updateSession(email)
	pass := req.Form.Get("pass")
	nombreServicio := req.Form.Get("nombreServicio")
	usuarioServicio := req.Form.Get("usuarioServicio")
	passServicio := req.Form.Get("passServicio")
	log.Println("crearCuenta: [" + email + ", " + pass + ", " + nombreServicio + ", " + usuarioServicio + ", " + passServicio + "]")

	// Añadimos el servicio a la BD
	database.AddAccountToUser(email, pass, nombreServicio, usuarioServicio, passServicio)

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")
	// Respondemos
	response(w, false, 501, "to-do")
}

// Modifica usuario de una cuenta (servicio)
func modificarCuenta(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos (servicio)
	email := req.Form.Get("email")
	updateSession(email)
	pass := req.Form.Get("pass")
	nombreServicio := req.Form.Get("nombreServicio")
	usuarioServicio := req.Form.Get("usuarioServicio")
	passServicio := req.Form.Get("passServicio")
	log.Println("modificarCuenta: [" + email + ", " + pass + ", " + nombreServicio + ", " + usuarioServicio + ", " + passServicio + " ]")

	database.SetAccount(email, pass, nombreServicio, usuarioServicio, passServicio)
	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")
	// Respondemos

	response(w, false, 200, "")
}

// Elimina una cuenta de servicio a un usuario de la BD
func eliminarCuenta(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	email := req.Form.Get("email")
	updateSession(email)
	pass := req.Form.Get("pass")
	nombreServicio := req.Form.Get("nombreServicio")
	log.Println("eliminarCuenta: [" + email + ", " + pass + ", " + nombreServicio + "]")

	database.DeleteAccount(email, pass, nombreServicio)
	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")
	// Respondemos
	response(w, false, 200, "")
}

// Recupera las cuentas de servicio de un usuario de la BD
func listarCuentas(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	email := req.Form.Get("email")
	updateSession(email)
	pass := req.Form.Get("pass")
	log.Println("listarCuentas: [" + email + ", " + pass + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")
	// Respondemos
	response(w, false, 201, database.GetJSONAllAccountsFromUser(email, pass))

}

// Recupera los detalles de una cuenta de servicio a un usuario de la BD
func detallesCuenta(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	email := req.Form.Get("email")
	updateSession(email)
	pass := req.Form.Get("pass")
	nombreServicio := req.Form.Get("nombreServicio")
	log.Println("detallesCuenta: [" + email + ", " + pass + ", " + nombreServicio + "]")
	accountInfo := database.GetJSONAccountFromUser(email, pass, nombreServicio)

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")
	// Respondemos
	response(w, false, 200, accountInfo)
}

func updateSession(email string) {
	if session[email].IsZero() {
		session[email] = time.Now()
	} else {
		duration := time.Now().Sub(session[email])
		if duration.Seconds() > 20 {
			log.Println("El usuario " + email + " ha superado el límite de sesión")
		} else {
			session[email] = time.Now()
		}
	}

}

func ClearSession(email string) {
	//session[email] = nil
}
