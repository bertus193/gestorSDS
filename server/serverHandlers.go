package server

import (
	"log"
	"net/http"

	"github.com/bertus193/gestorSDS/server/database"
)

// Comprueba si existe un usuario en la BD
func loginUsuario(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	email := req.Form.Get("email")
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
	pass := req.Form.Get("pass")
	log.Println("registroUsuario: [" + email + ", " + pass + "]")

	database.AddUser(email, pass)

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")
	// Respondemos
	response(w, false, 201, "")
}

// Añade una cuenta de servicio a un usuario de la BD
func crearCuenta(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	email := req.Form.Get("email")
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

// Modifica una cuenta de servicio a un usuario de la BD
func modificarCuenta(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	email := req.Form.Get("email")
	pass := req.Form.Get("pass")
	nombreServicio := req.Form.Get("nombreServicio")
	usuarioServicio := req.Form.Get("usuarioServicio")
	passServicioAnterior := req.Form.Get("passServicioAnterior")
	passServicioNueva := req.Form.Get("passServicioNueva")
	log.Println("modificarCuenta: [" + email + ", " + pass + ", " + nombreServicio + ", " + usuarioServicio + ", " + passServicioAnterior + ", " + passServicioNueva + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")
	// Respondemos
	response(w, false, 501, "to-do")
}

// Elimina una cuenta de servicio a un usuario de la BD
func eliminarCuenta(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	email := req.Form.Get("email")
	pass := req.Form.Get("pass")
	nombreServicio := req.Form.Get("nombreServicio")
	usuarioServicio := req.Form.Get("usuarioServicio")
	log.Println("eliminarCuenta: [" + email + ", " + pass + ", " + nombreServicio + ", " + usuarioServicio + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")
	// Respondemos
	response(w, false, 501, "to-do")
}

// Recupera las cuentas de servicio de un usuario de la BD
func listarCuentas(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	email := req.Form.Get("email")
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
	pass := req.Form.Get("pass")
	nombreServicio := req.Form.Get("nombreServicio")
	log.Println("detallesCuenta: [" + email + ", " + pass + ", " + nombreServicio + "]")
	accountInfo := database.GetJSONAccountFromUser(email, pass, nombreServicio)

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")
	// Respondemos
	response(w, false, 200, accountInfo)
}
