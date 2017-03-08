package server

import (
	"fmt"
	"net/http"

	"github.com/bertus193/gestorSDS/server/database"
)

func handler(w http.ResponseWriter, r *http.Request) {
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

	http.HandleFunc("/hello", handler)
	err := http.ListenAndServeTLS(":8021", "cert.pem", "key.pem", nil)
	if err != nil {
		panic(err)
	}
}
