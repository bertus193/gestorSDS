package utils

import (
	"bytes"
	"compress/zlib"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"

	"golang.org/x/crypto/scrypt"
)

// DeriveKey Genera un hash a partir de una contraseña y un sal
func DeriveKey(pass, salt []byte) ([]byte, error) {

	key, err := scrypt.Key(pass, salt, 16384, 8, 1, 32)
	if err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf("%x", key)), nil
}

// GenerateRandomBytes Genera cadenas aleatorias
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// err == nil only if len(b) == n
	if err != nil {
		return nil, err
	}

	return b, nil
}

// función para comprobar errores (ahorra escritura)
func chk(e error) {
	if e != nil {
		panic(e)
	}
}

// función para cifrar (con AES en este caso), adjunta el IV al principio
func Encrypt(data, key []byte) (out []byte) {
	out = make([]byte, len(data)+16)    // reservamos espacio para el IV al principio
	rand.Read(out[:16])                 // generamos el IV
	blk, err := aes.NewCipher(key)      // cifrador en bloque (AES), usa key
	chk(err)                            // comprobamos el error
	ctr := cipher.NewCTR(blk, out[:16]) // cifrador en flujo: modo CTR, usa IV
	ctr.XORKeyStream(out[16:], data)    // ciframos los datos
	return
}

// función para descifrar (con AES en este caso)
func Decrypt(data, key []byte) (out []byte) {
	out = make([]byte, len(data)-16)     // la salida no va a tener el IV
	blk, err := aes.NewCipher(key)       // cifrador en bloque (AES), usa key
	chk(err)                             // comprobamos el error
	ctr := cipher.NewCTR(blk, data[:16]) // cifrador en flujo: modo CTR, usa IV
	ctr.XORKeyStream(out, data[16:])     // desciframos (doble cifrado) los datos
	return
}

// función para codificar de []bytes a string (Base64)
func Encode64(data []byte) string {
	return base64.StdEncoding.EncodeToString(data) // sólo utiliza caracteres "imprimibles"
}

// función para decodificar de string a []bytes (Base64)
func Decode64(s string) []byte {
	b, err := base64.StdEncoding.DecodeString(s) // recupera el formato original
	chk(err)                                     // comprobamos el error
	return b                                     // devolvemos los datos originales
}

// generar contraseñas
func GeneratePassword(tamano int, letras bool, numeros bool, simbolos bool) string {
	arrayLetras := "AaBbCcDdEeFfGgHhIiJjKkLlMmNnOoPpQqRrSsTtUuVvWwXxYyZz"
	arrayNumeros := "0123456789"
	arraySimbolos := "+-/&<%#@*{!"

	arrayBase := ""
	salida := ""

	if letras == true {
		arrayBase = arrayBase + arrayLetras
	}
	if numeros == true {
		arrayBase = arrayBase + arrayNumeros
	}
	if simbolos == true {
		arrayBase = arrayBase + arraySimbolos
	}

	if arrayBase != "" {
		for i := 0; i < tamano; i++ {
			salida += string(arrayBase[cryptoRandSecure(int64(len(arrayBase)))])
		}
	} else {
		return "error"
	}

	return string(salida)

}

func cryptoRandSecure(max int64) int64 {
	nBig, err := rand.Int(rand.Reader, big.NewInt(max))
	if err != nil {
		log.Println(err)
	}
	return nBig.Int64()
}

func Compress(text string) string {
	var b bytes.Buffer
	w := zlib.NewWriter(&b)
	w.Write([]byte(text))
	w.Close()
	return string(b.Bytes())
	// Output: [120 156 202 72 205 201 201 215 81 40 207 47 202 73 225 2 4 0 0 255 255 33 231 4 147]

}

func Decompress(buff []byte) string {
	b := bytes.NewReader(buff)
	r, err := zlib.NewReader(b)
	if err != nil {
		panic(err)
	}
	if b, err := ioutil.ReadAll(r); err == nil {
		return string(b)
	}
	return "error"
}
