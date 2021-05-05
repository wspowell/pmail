package users

import (
	"net/http"
	"testing"

	"github.com/wspowell/errors"
	"github.com/wspowell/snailmail/resources"
	"github.com/wspowell/snailmail/server"
	"github.com/wspowell/snailmail/test/resources/mocks"
	"github.com/wspowell/spiderweb/server/servertest"
)

func Test_Create(t *testing.T) {
	t.Parallel()

	snailmail := server.New()

	{
		userstore := &mocks.UserStore{}
		userstore.On("CreateUser", "wpowell", resources.UserAttributes{
			PineappleOnPizza: true,
		}).Return(uint32(123), nil)

		servertest.TestRequest(t, snailmail, "Create new user",
			servertest.GivenRequest(http.MethodPost, "/users").
				WithResourceMock("userstore", userstore).
				WithRequestBody("application/json", []byte(`{
				"username": "wpowell",
				"pineapple_on_pizza": true
			}`)).
				ExpectResponse(http.StatusCreated).
				WithResponseBody("application/json", []byte(`{"user_id":123}`)))
	}

	{
		userstore := &mocks.UserStore{}
		userstore.On("CreateUser", "wpowell", resources.UserAttributes{
			PineappleOnPizza: true,
		}).Return(uint32(0), resources.ErrUsernameConflict)

		servertest.TestRequest(t, snailmail, "Username already exists",
			servertest.GivenRequest(http.MethodPost, "/users").
				WithResourceMock("userstore", userstore).
				WithRequestBody("application/json", []byte(`{
						"username": "wpowell",
						"pineapple_on_pizza": true
					}`)).
				ExpectResponse(http.StatusConflict).
				WithResponseBody("application/json", []byte(`{"message":"[resources-userstore-2] username already exists"}`)))
	}

	{
		userstore := &mocks.UserStore{}
		userstore.On("CreateUser", "wpowell", resources.UserAttributes{
			PineappleOnPizza: true,
		}).Return(uint32(0), errors.New("0000", "error"))

		servertest.TestRequest(t, snailmail, "Userstore failed",
			servertest.GivenRequest(http.MethodPost, "/users").
				WithResourceMock("userstore", userstore).
				WithRequestBody("application/json", []byte(`{
						"username": "wpowell",
						"pineapple_on_pizza": true
					}`)).
				ExpectResponse(http.StatusInternalServerError).
				WithResponseBody("application/json", []byte(`{"message":"[0000] error"}`)))
	}
}
