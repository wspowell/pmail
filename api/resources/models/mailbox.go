package models

import (
	"math/rand"
	"time"
)

const (
	addressCharLength = 12
	letterBytes       = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"
)

func randomAddress() string {
	randrand := rand.New(rand.NewSource(time.Now().Unix()))

	b := make([]byte, addressCharLength)
	for i := range b {
		// nolint:gosec // reason: no need for this to be secure
		b[i] = letterBytes[randrand.Intn(len(letterBytes))]
	}

	return string(b)
}

type Mailbox struct {
	// UserGuid of the owner.
	UserGuid string

	// Address of the mailbox to show the user.
	// Must be globally unique.
	// This is a code that can be used for others to send mail to you.
	Address string

	// Location in the world.
	Location Coordinate
}

func CreateMailbox(userGuid string, location Coordinate) Mailbox {
	return Mailbox{
		UserGuid: userGuid,
		// 599,555,620,984,320,000 permutations (36 character set, sets of 12)
		Address:  randomAddress(),
		Location: location,
	}
}

// FormatAddress as "AAAA-AAAA-AAAA"
func (self Mailbox) FormatAddress() string {
	formattedAddress := ""

	for i := 0; i < 4; i++ {
		formattedAddress += string(self.Address[i])
	}
	formattedAddress += "-"
	for i := 4; i < 8; i++ {
		formattedAddress += string(self.Address[i])
	}
	formattedAddress += "-"
	for i := 8; i < 12; i++ {
		formattedAddress += string(self.Address[i])
	}

	return formattedAddress
}

func (self Mailbox) IsNearby(location Coordinate, radiusMeters float32) bool {
	// TODO
	return true
}
