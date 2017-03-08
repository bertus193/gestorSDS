package server

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/bertus193/gestorSDS/config"
	"github.com/bertus193/gestorSDS/server/database"
)

func registroUsuario(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()
	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	// Recuperamos los datos
	email := req.Form.Get("email")
	pass := req.Form.Get("pass")
	log.Println("registroUsuario: [" + email + ", " + pass + "]")

	// Respondemos
	response(w, false, 501, "to-do")
}

func modificarUsuario(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()
	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	// Recuperamos los datos
	email := req.Form.Get("email")
	passAnterior := req.Form.Get("passAnterior")
	passNuevo := req.Form.Get("passNuevo")
	log.Println("modificarUsuario: [" + email + ", " + passAnterior + ", " + passNuevo + "]")

	// Respondemos
	response(w, false, 501, "to-do")
}

func crearCuenta(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()
	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	// Recuperamos los datos
	email := req.Form.Get("email")
	pass := req.Form.Get("pass")
	nombreServicio := req.Form.Get("nombreServicio")
	usuarioServicio := req.Form.Get("usuarioServicio")
	passServicio := req.Form.Get("passServicio")
	log.Println("crearCuenta: [" + email + ", " + pass + ", " + nombreServicio + ", " + usuarioServicio + ", " + passServicio + "]")

	// Respondemos
	response(w, false, 501, "to-do")
}

func modificarCuenta(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()
	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	// Recuperamos los datos
	email := req.Form.Get("email")
	pass := req.Form.Get("pass")
	nombreServicio := req.Form.Get("nombreServicio")
	usuarioServicio := req.Form.Get("usuarioServicio")
	passServicioAnterior := req.Form.Get("passServicioAnterior")
	passServicioNueva := req.Form.Get("passServicioNueva")
	log.Println("modificarCuenta: [" + email + ", " + pass + ", " + nombreServicio + ", " + usuarioServicio + ", " + passServicioAnterior + ", " + passServicioNueva + "]")

	// Respondemos
	response(w, false, 501, "to-do")
}

func eliminarCuenta(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()
	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	// Recuperamos los datos
	email := req.Form.Get("email")
	pass := req.Form.Get("pass")
	nombreServicio := req.Form.Get("nombreServicio")
	usuarioServicio := req.Form.Get("usuarioServicio")
	log.Println("eliminarCuenta: [" + email + ", " + pass + ", " + nombreServicio + ", " + usuarioServicio + "]")

	// Respondemos
	response(w, false, 501, "to-do")
}

func listarCuentas(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()
	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	// Recuperamos los datos
	email := req.Form.Get("email")
	pass := req.Form.Get("pass")
	log.Println("listarCuentas: [" + email + ", " + pass + "]")

	// Respondemos
	response(w, false, 201, database.GetJSONAllAccountsFromUser(email, pass))

}

func detallesCuenta(w http.ResponseWriter, req *http.Request) {
	// Parseamos el formulario
	req.ParseForm()
	// Cabecera estándar
	w.Header().Set("Content-Type", "text/plain")

	// Recuperamos los datos
	email := req.Form.Get("email")
	pass := req.Form.Get("pass")
	nombreServicio := req.Form.Get("nombreServicio")
	usuarioServicio := req.Form.Get("usuarioServicio")
	log.Println("detallesCuenta: [" + email + ", " + pass + ", " + nombreServicio + ", " + usuarioServicio + "]")

	// Respondemos
	response(w, false, 501, "to-do")
}

// función para escribir una respuesta del servidor
func response(w http.ResponseWriter, ok bool, code int, msgJSON string) {
	w.WriteHeader(code)
	fmt.Fprintf(w, msgJSON)
}

// Launch lanza el servidor
func Launch() {
	// todo: borrar
	// Datos de relleno para probar el acceso (BORRAR)
	database.AddUser("demoEmail", "demoMasterPass")
	database.AddAccountToUser("demoEmail", "facebook", "facebookUser", "facebookPass")
	database.AddAccountToUser("demoEmail", "twitter", "twitterUser", "twitterPass")
	database.AddUser("demoEmail2", "demoMasterPass2")

	// suscripción SIGINT
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	mux := http.NewServeMux()
	mux.Handle("/usuario/registro", http.HandlerFunc(registroUsuario))
	mux.Handle("/usuario/modificar", http.HandlerFunc(modificarUsuario))
	mux.Handle("/cuentas", http.HandlerFunc(listarCuentas))
	mux.Handle("/cuentas/nueva", http.HandlerFunc(crearCuenta))
	mux.Handle("/cuentas/modificar", http.HandlerFunc(modificarCuenta))
	mux.Handle("/cuentas/eliminar", http.HandlerFunc(eliminarCuenta))
	mux.Handle("/cuentas/detalles", http.HandlerFunc(detallesCuenta))

	//http.HandleFunc("/hello", handler)

	srv := &http.Server{Addr: config.SecureServerPort, Handler: mux}

	go func() {
		if err := srv.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	<-stopChan // espera señal SIGINT
	log.Println("Apagando servidor ...")

	// apagar servidor de forma segura
	ctx, fnc := context.WithTimeout(context.Background(), 5*time.Second)
	fnc()
	srv.Shutdown(ctx)

	log.Println("Servidor detenido correctamente")
}
