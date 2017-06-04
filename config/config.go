package config

// AppName contiene el nombre de la aplicación
var AppName = "Gestor SDS"

// SecureServerPort Puerto Seguro Cliente
var SecureServerPort = ":10443"

// SecureURL Url segura
var SecureURL = "https://127.0.0.1"

// MaxTimeSession es el tiempo máximo de sesión (segundos)
var MaxTimeSession = 60 * 30

// MaxA2FTime es el tiempo máximo de espera para resolver
// el reto de segundo factor de autenticación (segundos)
var MaxA2FTime = 60 * 5

// SizeA2FCode es el número de digitos que contendrá la clave
// que se envía al los usuario con A2F activado
var SizeA2FCode = 6

// PassDBEncrypt es la clave de cifrado del fichero de base de datos
var PassDBEncrypt = []byte("a very very very very secret key")

// Account2FA contiene los datos de la cuenta de correo
// encargada de enviar los códigos de inicio de sesión
var Account2FA = map[string]string{
	"email":      "gestorsds.ua@gmail.com",
	"passw":      "gestorsdspass",
	"smtpServer": "smtp.gmail.com",
	"smtpPort":   "587",
}

//CifrateLogs se encarga de indicar si se desea encriptar el log del servidor
var CifrateLogs = false

//PassCifrateLogs Password encriptacion Logs
var PassCifrateLogs = []byte("a really difficult logg password")
