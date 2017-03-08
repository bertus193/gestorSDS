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

func handler(w http.ResponseWriter, req *http.Request) {
	fmt.Printf("run")
	fmt.Fprintf(w, database.GetAll())
}

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
	fmt.Fprintf(w, database.GetAll())
}

func modificarUsuario(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, database.GetAll())
}

func crearCuenta(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, database.GetAll())
}

func modificarCuenta(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, database.GetAll())
}

func eliminarCuenta(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, database.GetAll())
}

func listarCuentas(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, database.GetAll())
}

func detallesCuenta(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, database.GetAll())
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
	mux.Handle("/hello", http.HandlerFunc(handler))
	mux.Handle("/usuario/registro", http.HandlerFunc(registroUsuario))
	mux.Handle("/usuario/modificar", http.HandlerFunc(modificarUsuario))
	mux.Handle("/cuentas/nueva", http.HandlerFunc(crearCuenta))
	mux.Handle("/cuentas/modificar", http.HandlerFunc(modificarCuenta))
	mux.Handle("/cuentas/eliminar", http.HandlerFunc(eliminarCuenta))
	mux.Handle("/cuentas", http.HandlerFunc(listarCuentas))
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
