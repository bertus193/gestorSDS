package main

import (
	"encoding/json"
	"fmt"
)

type usuario struct {
	IDUser   string
	NameUser string
	Datos    map[string]string
}

func addUser() {}

func addToUser() {}

func main() {

	gestor := make(map[string]usuario)
	gestor["usuarioDemo"] = usuario{IDUser: "idValue", NameUser: "userValue", Datos: make(map[string]string)}
	gestor["usuarioDemo2"] = usuario{IDUser: "idValue2", NameUser: "userValue2", Datos: make(map[string]string)}

	gestor["usuarioDemo"].Datos["facebook"] = "12345"
	gestor["usuarioDemo2"].Datos["facebook"] = "54321"

	fmt.Println(gestor)

	j, err := json.Marshal(gestor)
	fmt.Println(string(j))
	fmt.Println(err)

}
