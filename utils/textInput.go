package utils

import (
	"bufio"
	"os"

	"golang.org/x/crypto/ssh/terminal"
)

// CustomScanf es la alternativa a "Scanf" que nos permite
// leer entradas de texto por teclado que contienen espacios
func CustomScanf() string {
	var line = ""
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		line = scanner.Text()
	}
	return line
}

// GetPassw permite leer entradas de texto por terminal
// sin que se muestre el texto introducido en pantalla
func GetPassw() string {
	passwordStr := ""
	password, err := terminal.ReadPassword(0)
	if err == nil {
		passwordStr = string(password)
	}
	return passwordStr
}
