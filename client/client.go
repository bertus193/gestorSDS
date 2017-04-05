package client

import (
	"crypto/tls"
	"net/http"

	"github.com/bertus193/gestorSDS/config"
)

var baseURL = config.SecureURL + config.SecureServerPort

//Start Inicio Cliente
func Start() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	// two different ways to execute an http.GET

	//response, err := client.PostForm(baseURL+"/hello", data)
	// listarCuentas(client, "demoEmail", "hash_del_pass")

	startUI(client)
}
