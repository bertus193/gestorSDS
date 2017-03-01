package server

import (
	"fmt"
	"net/http"

	"./database"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, database.GetAll())
}

// redirectToHttps redirecciona a conexión segura
func redirectToHTTPS(w http.ResponseWriter, r *http.Request) {
	// Redirect the incoming HTTP request. Note that "127.0.0.1:8081" will only work if you are accessing the server from your local machine.
	http.Redirect(w, r, "https://127.0.0.1:8081"+r.RequestURI, http.StatusMovedPermanently)
}

// Launch lanza el servidor
func Launch() {
	// todo: borrar
	// Datos de relleno para probar el acceso (BORRAR)
	database.AddUser("demoEmail", "accountPass")
	database.AddToUser("demoEmail", "facebook", "12345")
	database.AddToUser("demoEmail", "twitter", "abcde")
	database.AddUser("demoEmail2", "accountPass2")

	http.HandleFunc("/", handler)
	// Start the HTTPS server in a goroutine
	go http.ListenAndServeTLS(":8081", "cert.pem", "key.pem", nil)
	// Start the HTTP server and redirect all incoming connections to HTTPS
	http.ListenAndServe(":8080", http.HandlerFunc(redirectToHTTPS))
}
