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

var logFile *os.File

func Init() {
	before()
}

// Launch lanza el servidor
func Launch() {

	Init()

	// suscripción SIGINT
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	mux := http.NewServeMux()
	mux.Handle("/usuario/login", http.HandlerFunc(loginUsuario))
	mux.Handle("/usuario/registro", http.HandlerFunc(registroUsuario))
	mux.Handle("/usuario/modificar", http.HandlerFunc(modificarUsuario))
	mux.Handle("/usuario/eliminar", http.HandlerFunc(eliminarUsuario))
	mux.Handle("/usuario/detalles", http.HandlerFunc(detallesUsuario))
	mux.Handle("/a2f/activar", http.HandlerFunc(activarA2F))
	mux.Handle("/a2f/desactivar", http.HandlerFunc(desactivarA2F))
	mux.Handle("/a2f/desbloquear", http.HandlerFunc(desbloquearA2F))
	mux.Handle("/cuentas", http.HandlerFunc(listarCuentas))
	mux.Handle("/vault/nueva", http.HandlerFunc(crearEntrada))
	mux.Handle("/cuentas/modificar", http.HandlerFunc(modificarCuenta))
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

func AddLog(logMessage string) {
	log.Println(logMessage)
	logMessage = time.Now().Format("2006-01-02 15:04:05") + " " + logMessage + "\n"
	logFile.Write([]byte(logMessage))
}

func before() {

	file, err := os.OpenFile("./server/logs.log", os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Printf("error opening file: %v", err)
	}

	logFile = file
}
