package server

import (
	"log"
	"net/http"

	"github.com/bertus193/gestorSDS/server/database"
)

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

func registroUsuario(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	email := req.Form.Get("email")
	pass := req.Form.Get("pass")
	log.Println("registroUsuario: [" + email + ", " + pass + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")
	// Respondemos
	response(w, false, 501, "to-do")
}

func modificarUsuario(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	email := req.Form.Get("email")
	passAnterior := req.Form.Get("passAnterior")
	passNuevo := req.Form.Get("passNuevo")
	log.Println("modificarUsuario: [" + email + ", " + passAnterior + ", " + passNuevo + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")
	// Respondemos
	response(w, false, 501, "to-do")
}

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

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")
	// Respondemos
	response(w, false, 501, "to-do")
}

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

func detallesCuenta(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()

	// Recuperamos los datos
	email := req.Form.Get("email")
	pass := req.Form.Get("pass")
	nombreServicio := req.Form.Get("nombreServicio")
	usuarioServicio := req.Form.Get("usuarioServicio")
	log.Println("detallesCuenta: [" + email + ", " + pass + ", " + nombreServicio + ", " + usuarioServicio + "]")

	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")
	// Respondemos
	response(w, false, 501, "to-do")
}
