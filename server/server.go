package server

import (
	"fmt"
	"net/http"

	"./database"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, database.GetAll())
}

// Launch lanza el servidor
func Launch() {
	// todo: borrar
	// Datos de relleno para probar el acceso (BORRAR)
	database.AddUser("demoEmail", "accountPass")
	database.AddToUser("demoEmail", "facebook", "12345")
	database.AddToUser("demoEmail", "twitter", "abcde")
	database.AddUser("demoEmail2", "accountPass2")

	http.HandleFunc("/hello", handler)
	err := http.ListenAndServeTLS(":8021", "cert.pem", "key.pem", nil)
	if err != nil {
		panic(err)
	}
}
