package models

import (
	"math/rand"
	"time"

	"github.com/google/uuid"
)

const (
	signatureCharLength = 128
	signatureBytes      = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz!@#$%^&*()_-+=<>?"
)

func randomSignature() string {
	randrand := rand.New(rand.NewSource(time.Now().Unix()))

	b := make([]byte, signatureCharLength)
	for i := range b {
		// nolint:gosec // reason: no need for this to be secure
		b[i] = signatureBytes[randrand.Intn(len(signatureBytes))]
	}

	return string(b)
}

type User struct {
	Guid      string
	PublicKey string
	Signature string
	// CreatedOn timestamp.
	CreatedOn time.Time
	Mailbox   Mailbox
}

func CreateUser(publicKey string, location Coordinate) User {
	userGuid := uuid.New().String()

	return User{
		Guid:      userGuid,
		PublicKey: publicKey,
		Signature: randomSignature(),
		CreatedOn: time.Now().UTC(),
		Mailbox:   CreateMailbox(userGuid, location),
	}
}
