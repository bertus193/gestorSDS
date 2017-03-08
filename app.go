package main

import (
	"fmt"
	"os"

	"github.com/bertus193/gestorSDS/client"
	"github.com/bertus193/gestorSDS/server"
)

func main() {

	// Recogemos el valor de los argumentos
	if len(os.Args) == 2 {

		argMode := os.Args[1]

		if argMode == "client" {
			client.Start()
		} else if argMode == "server" {
			server.Launch()
		}

	} else {
		fmt.Printf("El número de parametros indicado no es correcto\n")
	}
}
