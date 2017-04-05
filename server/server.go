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

// función para escribir una respuesta del servidor
func response(w http.ResponseWriter, ok bool, code int, msgJSON string) {
	w.WriteHeader(code)
	fmt.Fprintf(w, msgJSON)
}

// Launch lanza el servidor
func Launch() {

	// suscripción SIGINT
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	mux := http.NewServeMux()
	mux.Handle("/usuario/login", http.HandlerFunc(loginUsuario))
	mux.Handle("/usuario/registro", http.HandlerFunc(registroUsuario))
	mux.Handle("/usuario/modificar", http.HandlerFunc(modificarUsuario))
	mux.Handle("/cuentas", http.HandlerFunc(listarCuentas))
	mux.Handle("/cuentas/nueva", http.HandlerFunc(crearCuenta))
	mux.Handle("/cuentas/modificar/usuario", http.HandlerFunc(modificarUsuarioCuenta))
	mux.Handle("/cuentas/modificar/password", http.HandlerFunc(modificarPasswordCuenta))
	mux.Handle("/cuentas/eliminar", http.HandlerFunc(eliminarCuenta))
	mux.Handle("/cuentas/detalles", http.HandlerFunc(detallesCuenta))

	srv := &http.Server{Addr: config.SecureServerPort, Handler: mux}

	go func() {
		if err := srv.ListenAndServeTLS("cert.pem", "key.pem"); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	<-stopChan // espera señal SIGINT
	log.Println("Apagando servidor ...")

	// Guarda la información de la BD en un fichero
	database.After()

	// Apaga el servidor de forma segura
	ctx, fnc := context.WithTimeout(context.Background(), 5*time.Second)
	fnc()
	srv.Shutdown(ctx)

	log.Println("Servidor detenido correctamente")
}
