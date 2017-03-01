package main

import (
	"fmt"
	"os"

	"./client"
	"./server"
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
