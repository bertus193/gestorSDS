package client

import (
	"crypto/tls"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/bertus193/gestorSDS/config"
)

var baseURL = config.SecureURL + config.SecureServerPort

func registroUsuario(client *http.Client, email string, pass string) (*http.Response, error) {
	data := url.Values{}
	data.Set("email", email)
	data.Set("pass", pass)

	return client.PostForm(baseURL+"/usuario/registro", data)

}

//Start Inicio Cliente
func Start() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	// two different ways to execute an http.GET

	//response, err := client.PostForm(baseURL+"/hello", data)
	response, err := registroUsuario(client, "alu@alu.ua.es", "hash_del_pass")
	if err != nil {
		log.Fatal(err)
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(contents))
	}
}
