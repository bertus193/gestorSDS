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

func listarCuentas(client *http.Client, email string, pass string) (*http.Response, error) {
	data := url.Values{}
	data.Set("email", email)
	data.Set("pass", pass)

	return client.PostForm(baseURL+"/cuentas", data)
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
