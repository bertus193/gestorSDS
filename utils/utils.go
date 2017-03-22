package utils

import (
	"crypto/rand"
	"fmt"

	"golang.org/x/crypto/scrypt"
)

// deriveKey generates a new NaCl key from a passphrase and salt.
func DeriveKey(pass, salt []byte) ([]byte, error) {

	key, err := scrypt.Key(pass, salt, 16384, 8, 1, 32)
	if err != nil {
		return nil, err
	}

	return []byte(fmt.Sprintf("%x", key)), nil
}

func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	// err == nil only if len(b) == n
	if err != nil {
		return nil, err
	}

	return b, nil
}
