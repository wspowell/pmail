package auth

import (
	"crypto/sha512"
	"fmt"

	"github.com/wspowell/context"
	"github.com/wspowell/errors"

	"github.com/wspowell/snailmail/resources/aws"
)

const (
	icHashPasswordAppPassword = "auth-password-1"
)

type appPassword struct {
	Password string `json:"password" envDefault:"password"`
}

func Password(ctx context.Context, password string) (string, error) {
	var appSalt appPassword
	if err := aws.GetSecret(ctx, &appSalt); err != nil {
		return "", errors.Propagate(icHashPasswordAppPassword, err)
	}

	sum := sha512.Sum512([]byte(password + appSalt.Password))

	return fmt.Sprintf("%x", sum), nil
}
