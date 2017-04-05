package client

import (
	"crypto/sha512"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	"github.com/bertus193/gestorSDS/model"
	"github.com/bertus193/gestorSDS/utils"
)

var userLogin string
var keyLogin string
var keyData []byte

func loginUsuario(client *http.Client, email string, pass string) bool {
	data := url.Values{}
	data.Set("email", email)

	keyClient := sha512.Sum512([]byte(pass))
	keyData = keyClient[32:64]
	keyLogin = utils.Encode64(keyClient[0:31])
	userLogin = email

	data.Set("pass", keyLogin)

	response, err := client.PostForm(baseURL+"/usuario/login", data)
	if err != nil {
		log.Fatal(err)
	} else {
		// Cerramos la conexión
		defer response.Body.Close()

		if response.StatusCode == 200 {
			return true
		}
	}
	return false
}

func registroUsuario(client *http.Client, email string, pass string) (*http.Response, error) {
	data := url.Values{}
	data.Set("email", email)

	keyClient := sha512.Sum512([]byte(pass))
	keyRegister := keyClient[0:31]

	data.Set("pass", utils.Encode64(keyRegister))

	return client.PostForm(baseURL+"/usuario/registro", data)
}

/*func modificarUsuario(client *http.Client, email string, passAnterior string, passNuevo string) (*http.Response, error) {
	data := url.Values{}
	data.Set("email", email)
	data.Set("passAnterior", passAnterior)
	data.Set("passNuevo", passNuevo)

	return client.PostForm(baseURL+"/usuario/modificar", data)
}*/

func crearCuenta(client *http.Client, nombreServicio string, usuarioServicio string, passServicio string) (*http.Response, error) {
	data := url.Values{}
	data.Set("email", userLogin)
	data.Set("pass", keyLogin)
	data.Set("nombreServicio", nombreServicio)
	data.Set("usuarioServicio", usuarioServicio)

	encryptPassServicio := utils.Encode64(utils.Encrypt([]byte(passServicio), keyData))
	data.Set("passServicio", encryptPassServicio)

	return client.PostForm(baseURL+"/cuentas/nueva", data)
}

func eliminarCuenta(client *http.Client, email string, pass string, nombreServicio string, usuarioServicio string) (*http.Response, error) {
	data := url.Values{}
	data.Set("email", email)
	data.Set("pass", pass)
	data.Set("nombreServicio", nombreServicio)
	data.Set("usuarioServicio", usuarioServicio)

	return client.PostForm(baseURL+"/cuentas/eliminar", data)
}

func listarCuentas(client *http.Client) map[string]model.Account {
	data := url.Values{}
	data.Set("email", userLogin)
	data.Set("pass", keyLogin)

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
		// fmt.Println(response.StatusCode)

		// Recuperamos el objeto del mensaje origianl
		result := make(map[string]model.Account)
		if err := json.Unmarshal(contents, &result); err == nil {
			return result
		}
	}

	return nil
}

func detallesCuenta(client *http.Client, nombreServicio string) model.Account {
	data := url.Values{}
	data.Set("email", userLogin)
	data.Set("pass", keyLogin)
	data.Set("nombreServicio", nombreServicio)

	response, err := client.PostForm(baseURL+"/cuentas/detalles", data)
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
		// fmt.Println(response.StatusCode)

		// Recuperamos el objeto del mensaje origianl
		result := model.Account{}
		if err := json.Unmarshal(contents, &result); err == nil {
			return result
		}
	}

	return model.Account{}
}

func modificarCuenta(client *http.Client, usuarioServicio string, passServicio string, nombreServicio string) (*http.Response, error) {
	data := url.Values{}
	data.Set("email", userLogin)
	data.Set("pass", keyLogin)
	data.Set("nombreServicio", nombreServicio)
	data.Set("usuarioServicio", usuarioServicio)

	encryptPassServicio := utils.Encode64(utils.Encrypt([]byte(passServicio), keyData))
	data.Set("passServicio", encryptPassServicio)

	return client.PostForm(baseURL+"/cuentas/modificar", data)
}
