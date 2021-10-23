package db

import (
	"testing"

	"github.com/wspowell/snailmail/resources/user"
)

func Test_InMemory_user_Api(t *testing.T) {
	database := NewInMemory()
	userClient := user.NewClient(database)

	test.RunApiTestCases(t, userClient)
}
