package client

import (
	"crypto/sha512"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"fmt"

	"github.com/bertus193/gestorSDS/model"
	"github.com/bertus193/gestorSDS/utils"
)

var sessionToken string
var keyData []byte

// Petición al servidor de creación de nuevo usuario
func registroUsuario(client *http.Client, email string, pass string) error {
	var errResult error

	data := url.Values{}

	// Generamos el hash de la contraseña introducida
	keyClient := sha512.Sum512([]byte(pass))
	// Usamos solo la primera parte para identificarnos
	keyRegister := keyClient[0:31]

	data.Set("email", email)
	data.Set("pass", utils.Encode64(keyRegister))

	// Realizamos la petición
	response, err := client.PostForm(baseURL+"/usuario/registro", data)
	if err == nil {
		// Si el código de estado recibido no es el esperado (201)
		if response.StatusCode != 201 {

			// Comprobamos el código de estado recibido
			switch response.StatusCode {
			case 409: // (409 - Conflict)
				errResult = errors.New("user already exists")
			default:
				errResult = errors.New("unknown")
			}
		}
	} else {
		// La petición al servidor no ha obtenido respuesta
		fmt.Println("* No se ha podido comunicar con el servidor")
		os.Exit(0)
	}
	// Cerramos la conexión
	defer response.Body.Close()

	return errResult
}

func loginUsuario(client *http.Client, email string, pass string) error {
	var errResult error

	data := url.Values{}
	data.Set("email", email)

	keyClient := sha512.Sum512([]byte(pass))
	keyLogin := utils.Encode64(keyClient[0:31])
	keyData = keyClient[32:64]

	data.Set("pass", keyLogin)

	// Realizamos la petición
	response, err := client.PostForm(baseURL+"/usuario/login", data)
	if err == nil {
		// Si el código de estado recibido no es el esperado (200)
		if response.StatusCode != 200 {

			// Comprobamos el código de estado recibido
			switch response.StatusCode {
			case 404: // (404 - Not found)
				errResult = errors.New("user not found")
			case 400: // (400 - Bad Request)
				errResult = errors.New("passwords do not match")
			case 250: // (250 - A2F required [custom])
				// Guardamos el token
				bodyBytes, _ := ioutil.ReadAll(response.Body)
				sessionToken = string(bodyBytes)
				// Solicitamos la resolución de A2f
				errResult = errors.New("a2f required")
			default:
				errResult = errors.New("unknown")
			}
		} else {
			// Guardamos el token
			bodyBytes, _ := ioutil.ReadAll(response.Body)
			sessionToken = string(bodyBytes)
		}
	} else {
		// La petición al servidor no ha obtenido respuesta
		fmt.Println("* No se ha podido comunicar con el servidor")
		os.Exit(0)
	}
	// Cerramos la conexión
	defer response.Body.Close()

	return errResult
}

func desbloquearA2F(client *http.Client, a2fcode string) error {

	var errResult error

	data := url.Values{}
	data.Set("token", sessionToken)
	data.Set("a2fcode", a2fcode)

	// Realizamos la petición
	response, err := client.PostForm(baseURL+"/a2f/desbloquear", data)
	if err == nil {
		// Si el código de estado recibido no es el esperado (200)
		if response.StatusCode != 200 {

			// Comprobamos el código de estado recibido
			switch response.StatusCode {
			case 404: // (404 - Not found)
				errResult = errors.New("session not found")
			case 304: // (304 - Not Modified)
				errResult = errors.New("2fa already resolved")
			case 408: // (408 - Request Timeout)
				errResult = errors.New("2fa expired")
			case 400: // (400 - Bad Request)
				errResult = errors.New("incorrect 2fa code")
			default:
				errResult = errors.New("unknown")
			}
		}

	} else {
		// La petición al servidor no ha obtenido respuesta
		fmt.Println("* No se ha podido comunicar con el servidor")
		os.Exit(0)
	}
	// Cerramos la conexión
	defer response.Body.Close()

	return errResult
}

func listarCuentas(client *http.Client) ([]string, error) {

	var entriesResult []string
	var errResult error

	data := url.Values{}
	data.Set("token", sessionToken)

	// Realizamos la petición
	response, err := client.PostForm(baseURL+"/cuentas", data)
	if err == nil {
		// Si el código de estado recibido no es el esperado (200)
		if response.StatusCode != 200 {

			// Comprobamos el código de estado recibido
			switch response.StatusCode {
			case 401: // (401 - Unauthorized)
				errResult = errors.New("unauthorized")
			default:
				errResult = errors.New("unknown")
			}

		} else {
			// Leemos la respuesta
			if contents, errRead := ioutil.ReadAll(response.Body); errRead != nil {
				errResult = errors.New("unable to read")
			} else {
				result := make([]string, 0)
				// Recuperamos el objeto del mensaje original
				if errJSON := json.Unmarshal(contents, &result); errJSON != nil {
					errResult = errors.New("unable to unmarshal")
				} else {
					entriesResult = result
				}
			}
		}

	} else {
		// La petición al servidor no ha obtenido respuesta
		fmt.Println("* No se ha podido comunicar con el servidor")
		os.Exit(0)
	}
	// Cerramos la conexión
	defer response.Body.Close()

	return entriesResult, errResult
}

func crearEntrada(client *http.Client, tituloEntrada string, usuario string, password string) error {

	var errResult error

	data := url.Values{}
	data.Set("token", sessionToken)
	data.Set("tituloEntrada", tituloEntrada)
	data.Set("usuarioCuenta", usuario)

	encryptPassServicio := utils.Encode64(utils.Encrypt([]byte(password), keyData))
	data.Set("passwordCuenta", encryptPassServicio)

	// Realizamos la petición
	response, err := client.PostForm(baseURL+"/vault/nueva", data)

	if err == nil {
		// Si el código de estado recibido no es el esperado (201)
		if response.StatusCode != 201 {

			// Comprobamos el código de estado recibido
			switch response.StatusCode {
			case 401: // (401 - Unauthorized)
				errResult = errors.New("unauthorized")
			case 404: // (404 - Not found)
				errResult = errors.New("user not found")
			case 409: // (409 - Conflict)
				errResult = errors.New("entry already exists")
			default:
				errResult = errors.New("unknown")
			}
		}

	} else {
		// La petición al servidor no ha obtenido respuesta
		fmt.Println("* No se ha podido comunicar con el servidor")
		os.Exit(0)
	}
	// Cerramos la conexión
	defer response.Body.Close()

	return errResult
}

func detallesEntrada(client *http.Client, tituloEntrada string) (model.VaultEntry, error) {

	var errResult error
	detailResult := model.VaultEntry{}

	data := url.Values{}
	data.Set("token", sessionToken)
	data.Set("tituloEntrada", tituloEntrada)

	// Realizamos la petición
	response, err := client.PostForm(baseURL+"/vault/detalles", data)

	if err == nil {
		// Si el código de estado recibido no es el esperado (201)
		if response.StatusCode != 201 {

			// Comprobamos el código de estado recibido
			switch response.StatusCode {
			case 401: // (401 - Unauthorized)
				errResult = errors.New("unauthorized")
			case 404: // (404 - Not found)
				errResult = errors.New("not found")
			default:
				errResult = errors.New("unknown")
			}
		} else {

			// Leemos la respuesta
			if contents, errRead := ioutil.ReadAll(response.Body); errRead != nil {
				errResult = errors.New("unable to read")
			} else {

				tempEntry := model.VaultEntry{}

				// Recuperamos el objeto del mensaje origianl
				if errJSON := json.Unmarshal(contents, &tempEntry); errJSON != nil {
					errResult = errors.New("unable to unmarshal")
				} else {
					// Desciframos la contraseña
					detailResult = model.VaultEntry{
						Mode:     1, // Account
						User:     tempEntry.User,
						Password: string(utils.Decrypt(utils.Decode64(tempEntry.Password), keyData)),
					}
				}
			}
		}

	} else {
		// La petición al servidor no ha obtenido respuesta
		fmt.Println("* No se ha podido comunicar con el servidor")
		os.Exit(0)
	}
	// Cerramos la conexión
	defer response.Body.Close()

	return detailResult, errResult
}

func eliminarEntrada(client *http.Client, tituloEntrada string) error {

	var errResult error

	data := url.Values{}
	data.Set("token", sessionToken)
	data.Set("tituloEntrada", tituloEntrada)

	// Realizamos la petición
	response, err := client.PostForm(baseURL+"/vault/eliminar", data)

	if err == nil {
		// Si el código de estado recibido no es el esperado (201)
		if response.StatusCode != 200 {

			// Comprobamos el código de estado recibido
			switch response.StatusCode {
			case 401: // (401 - Unauthorized)
				errResult = errors.New("unauthorized")
			case 404: // (404 - Not found)
				errResult = errors.New("not found")
			default:
				errResult = errors.New("unknown")
			}
		}

	} else {
		// La petición al servidor no ha obtenido respuesta
		fmt.Println("* No se ha podido comunicar con el servidor")
		os.Exit(0)
	}
	// Cerramos la conexión
	defer response.Body.Close()

	return errResult
}

func detallesUsuario(client *http.Client) (model.DetallesUsuario, error) {

	var errResult error
	detailResult := model.DetallesUsuario{}

	data := url.Values{}
	data.Set("token", sessionToken)

	// Realizamos la petición
	response, err := client.PostForm(baseURL+"/usuario/detalles", data)

	if err == nil {
		// Si el código de estado recibido no es el esperado (200)
		if response.StatusCode != 200 {

			// Comprobamos el código de estado recibido
			switch response.StatusCode {
			case 401: // (401 - Unauthorized)
				errResult = errors.New("unauthorized")
			case 404: // (404 - Not found)
				errResult = errors.New("user not found")
			default:
				errResult = errors.New("unknown")
			}
		} else {

			// Leemos la respuesta
			if contents, errRead := ioutil.ReadAll(response.Body); errRead != nil {
				errResult = errors.New("unable to read")
			} else {

				tempResult := model.DetallesUsuario{}

				// Recuperamos el objeto del mensaje origianl
				if errJSON := json.Unmarshal(contents, &tempResult); errJSON != nil {
					errResult = errors.New("unable to unmarshal")
				} else {
					detailResult = tempResult
				}
			}
		}

	} else {
		// La petición al servidor no ha obtenido respuesta
		fmt.Println("* No se ha podido comunicar con el servidor")
		os.Exit(0)
	}
	// Cerramos la conexión
	defer response.Body.Close()

	return detailResult, errResult
}

func updateA2F(client *http.Client, activar bool) error {

	var errResult error

	data := url.Values{}
	data.Set("token", sessionToken)

	var ruta string
	if activar {
		ruta = "/a2f/activar"
	} else {
		ruta = "/a2f/desactivar"
	}

	// Realizamos la petición
	response, err := client.PostForm(baseURL+ruta, data)

	if err == nil {
		// Si el código de estado recibido no es el esperado (200)
		if response.StatusCode != 200 {

			// Comprobamos el código de estado recibido
			switch response.StatusCode {
			case 401: // (401 - Unauthorized)
				errResult = errors.New("unauthorized")
			case 404: // (404 - Not found)
				errResult = errors.New("user not found")
			default:
				errResult = errors.New("unknown")
			}
		}

	} else {
		// La petición al servidor no ha obtenido respuesta
		fmt.Println("* No se ha podido comunicar con el servidor")
		os.Exit(0)
	}
	// Cerramos la conexión
	defer response.Body.Close()

	return errResult
}

func eliminarUsuario(client *http.Client) error {

	var errResult error

	data := url.Values{}
	data.Set("token", sessionToken)

	// Realizamos la petición
	response, err := client.PostForm(baseURL+"/usuario/eliminar", data)

	if err == nil {
		// Si el código de estado recibido no es el esperado (201)
		if response.StatusCode != 200 {

			// Comprobamos el código de estado recibido
			switch response.StatusCode {
			case 401: // (401 - Unauthorized)
				errResult = errors.New("unauthorized")
			case 404: // (404 - Not found)
				errResult = errors.New("user not found")
			default:
				errResult = errors.New("unknown")
			}
		}

	} else {
		// La petición al servidor no ha obtenido respuesta
		fmt.Println("* No se ha podido comunicar con el servidor")
		os.Exit(0)
	}
	// Cerramos la conexión
	defer response.Body.Close()

	return errResult
}
