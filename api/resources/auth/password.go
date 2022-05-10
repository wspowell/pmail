package auth

import (
	"crypto/sha512"
	"fmt"

	"golang.org/x/crypto/pbkdf2"
)

const (
	icHashPasswordAppPassword = "auth-password-1"
)

type appPassword struct {
	Password string `json:"password" envDefault:"password"`
}

func Sign(publicKey string, signature string) string {
	return fmt.Sprintf("%x", string(pbkdf2.Key([]byte(publicKey), []byte(signature), 4096, 128, sha512.New)))
}
