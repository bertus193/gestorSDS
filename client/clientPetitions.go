package client

import (
	"crypto/sha512"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"

	"fmt"

	"github.com/bertus193/gestorSDS/model"
	"github.com/bertus193/gestorSDS/utils"
)

var sessionToken string
var keyData []byte

func loginUsuario(client *http.Client, email string, pass string) bool {
	data := url.Values{}
	data.Set("email", email)

	keyClient := sha512.Sum512([]byte(pass))
	keyLogin := utils.Encode64(keyClient[0:31])
	keyData = keyClient[32:64]

	data.Set("pass", keyLogin)

	response, err := client.PostForm(baseURL+"/usuario/login", data)
	if err != nil {
		fmt.Println("* No se ha podido comunicar con el servidor")
		os.Exit(0)
	} else {
		// Cerramos la conexión
		defer response.Body.Close()

		if response.StatusCode == 200 {
			bodyBytes, _ := ioutil.ReadAll(response.Body)
			sessionToken = string(bodyBytes)
			return true
		}

		if response.StatusCode == 250 {
			bodyBytes, _ := ioutil.ReadAll(response.Body)
			sessionToken = string(bodyBytes)
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

func crearCuenta(client *http.Client, nombreServicio string, usuarioServicio string, passServicio string) (*http.Response, error, string) {
	data := url.Values{}
	data.Set("token", sessionToken)
	data.Set("nombreServicio", nombreServicio)
	data.Set("usuarioServicio", usuarioServicio)

	encryptPassServicio := utils.Encode64(utils.Encrypt([]byte(passServicio), keyData))
	data.Set("passServicio", encryptPassServicio)

	response, err := client.PostForm(baseURL+"/cuentas/nueva", data)

	if err != nil {
		log.Fatal(err)
	} else {
		// Cerramos la conexión
		defer response.Body.Close()
	}
	if response.StatusCode != 201 {
		return response, err, "errorSesion"
	}

	return response, err, ""
}

func eliminarUsuario(client *http.Client) (*http.Response, error, string) {
	data := url.Values{}
	data.Set("token", sessionToken)

	response, err := client.PostForm(baseURL+"/usuario/eliminar", data)

	if err != nil {
		log.Fatal(err)
	} else {
		// Cerramos la conexión
		defer response.Body.Close()
	}
	if response.StatusCode != 200 {
		return response, err, "errorSesion"
	}

	return response, err, ""
}

func modificarCuenta(client *http.Client, usuarioServicio string, passServicio string, nombreServicio string) (*http.Response, error, string) {
	data := url.Values{}
	data.Set("token", sessionToken)
	data.Set("nombreServicio", nombreServicio)
	data.Set("usuarioServicio", usuarioServicio)

	encryptPassServicio := utils.Encode64(utils.Encrypt([]byte(passServicio), keyData))
	data.Set("passServicio", encryptPassServicio)

	response, err := client.PostForm(baseURL+"/cuentas/modificar", data)

	if err != nil {
		log.Fatal(err)
	} else {
		// Cerramos la conexión
		defer response.Body.Close()
	}
	if response.StatusCode != 200 {
		return response, err, "errorSesion"
	}

	return response, err, ""
}

func eliminarCuenta(client *http.Client, nombreServicio string) (*http.Response, string) {
	data := url.Values{}
	data.Set("token", sessionToken)
	data.Set("nombreServicio", nombreServicio)

	response, err := client.PostForm(baseURL+"/cuentas/eliminar", data)

	if err != nil {
		log.Fatal(err)
	} else {
		// Cerramos la conexión
		defer response.Body.Close()
	}
	if response.StatusCode != 200 {
		return response, "errorSesion"
	}

	return response, ""
}

func listarCuentas(client *http.Client) (map[string]model.Account, string) {

	data := url.Values{}
	data.Set("token", sessionToken)

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

		result := make(map[string]model.Account)

		// Recuperamos el código http
		// fmt.Println(response.StatusCode)

		//to-do comprobar todos los codigos de error
		if response.StatusCode != 200 {
			return result, "errorSesion"
		}

		// Recuperamos el objeto del mensaje origianl
		if err := json.Unmarshal(contents, &result); err == nil {
			return result, ""
		}
	}

	return nil, "error"
}

func detallesCuenta(client *http.Client, nombreServicio string) (model.Account, string) {
	data := url.Values{}
	data.Set("token", sessionToken)
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

		result := model.Account{}

		// Recuperamos el código http
		// fmt.Println(response.StatusCode)

		//to-do comprobar todos los codigos de error
		if response.StatusCode != 200 {
			return result, "errorSesion"
		}

		// Recuperamos el objeto del mensaje origianl
		if err := json.Unmarshal(contents, &result); err == nil {
			return result, ""
		}
	}

	return model.Account{}, "error"
}
