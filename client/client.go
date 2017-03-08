package client

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/bertus193/gestorSDS/config"
	"github.com/bertus193/gestorSDS/model"
)

var baseURL = config.SecureURL + config.SecureServerPort

func registroUsuario(client *http.Client, email string, pass string) (*http.Response, error) {
	data := url.Values{}
	data.Set("email", email)
	data.Set("pass", pass)

	return client.PostForm(baseURL+"/usuario/registro", data)
}

func modificarUsuario(client *http.Client, email string, passAnterior string, passNuevo string) (*http.Response, error) {
	data := url.Values{}
	data.Set("email", email)
	data.Set("passAnterior", passAnterior)
	data.Set("passNuevo", passNuevo)

	return client.PostForm(baseURL+"/usuario/modificar", data)
}

func crearCuenta(client *http.Client, email string, pass string, nombreServicio string, usuarioServicio string, passServicio string) (*http.Response, error) {
	data := url.Values{}
	data.Set("email", email)
	data.Set("pass", pass)
	data.Set("nombreServicio", nombreServicio)
	data.Set("usuarioServicio", usuarioServicio)
	data.Set("passServicio", passServicio)

	return client.PostForm(baseURL+"/cuentas/nueva", data)
}

func modificarCuenta(client *http.Client, email string, pass string, nombreServicio string, usuarioServicio string, passServicioAnterior string, passServicioNueva string) (*http.Response, error) {
	data := url.Values{}
	data.Set("email", email)
	data.Set("pass", pass)
	data.Set("nombreServicio", nombreServicio)
	data.Set("usuarioServicio", usuarioServicio)
	data.Set("passServicioAnterior", passServicioAnterior)
	data.Set("passServicioNueva", passServicioAnterior)

	return client.PostForm(baseURL+"/cuentas/modificar", data)
}

func eliminarCuenta(client *http.Client, email string, pass string, nombreServicio string, usuarioServicio string) (*http.Response, error) {
	data := url.Values{}
	data.Set("email", email)
	data.Set("pass", pass)
	data.Set("nombreServicio", nombreServicio)
	data.Set("usuarioServicio", usuarioServicio)

	return client.PostForm(baseURL+"/cuentas/eliminar", data)
}

func listarCuentas(client *http.Client, email string, pass string) {
	data := url.Values{}
	data.Set("email", email)
	data.Set("pass", pass)

	response, err := client.PostForm(baseURL+"/cuentas", data)
	if err != nil {
		log.Fatal(err)
	} else {
		// Cerramos la conexión
		defer response.Body.Close()

		// Leemos la respuesta
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		// Recuperamos el código http
		fmt.Println(response.StatusCode)

		// Recuperamos el objeto del mensaje origianl
		result := make(map[string]model.Account)
		if err := json.Unmarshal(contents, &result); err != nil {
			fmt.Println("Error al recuperar el objeto")
		}

		// Imprimimos los resultados
		for k := range result {
			tempAccount := result[k]
			fmt.Println("[" + k + "]")
			fmt.Println("--> " + tempAccount.User)
			fmt.Println("--> " + tempAccount.Password)
		}
	}
}

func detallesCuenta(client *http.Client, email string, pass string, nombreServicio string, usuarioServicio string) (*http.Response, error) {
	data := url.Values{}
	data.Set("email", email)
	data.Set("pass", pass)
	data.Set("nombreServicio", nombreServicio)
	data.Set("usuarioServicio", usuarioServicio)

	return client.PostForm(baseURL+"/cuentas/detalles", data)
}

//Start Inicio Cliente
func Start() {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}
	// two different ways to execute an http.GET

	//response, err := client.PostForm(baseURL+"/hello", data)
	listarCuentas(client, "demoEmail", "hash_del_pass")
}
