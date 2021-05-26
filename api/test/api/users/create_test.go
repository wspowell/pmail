package users

import (
	"net/http"
	"testing"

	"github.com/wspowell/snailmail/resources"
	"github.com/wspowell/snailmail/server"
	"github.com/wspowell/snailmail/test/resources/mocks"

	"github.com/stretchr/testify/mock"
	"github.com/wspowell/errors"
	"github.com/wspowell/spiderweb/server/restful/restfultest"
)

func Test_Create(t *testing.T) {
	snailmail := server.New()

	{
		userstore := &mocks.UserStore{}
		userAttributes := resources.UserAttributes{
			PineappleOnPizza: true,
		}
		userstore.On("CreateUser", mock.Anything, "wpowell", userAttributes).Return(&resources.User{
			Id:         123,
			Username:   "wpowell",
			Attributes: userAttributes,
		}, nil)

		restfultest.TestCase(t, snailmail, "Create new user").
			GivenRequest(http.MethodPost, "/users").
			WithResourceMock("userstore", userstore).
			WithRequestBody("application/json", []byte(`{
				"username": "wpowell",
				"pineapple_on_pizza": true
			}`)).
			ExpectResponse(http.StatusCreated).
			WithResponseBody("application/json", []byte(`{"user_id":123}`)).
			Run()
	}

	{
		userstore := &mocks.UserStore{}
		userstore.On("CreateUser", mock.Anything, "wpowell", resources.UserAttributes{
			PineappleOnPizza: true,
		}).Return(nil, resources.ErrUsernameConflict)

		restfultest.TestCase(t, snailmail, "Username already exists").
			GivenRequest(http.MethodPost, "/users").
			WithResourceMock("userstore", userstore).
			WithRequestBody("application/json", []byte(`{
						"username": "wpowell",
						"pineapple_on_pizza": true
					}`)).
			ExpectResponse(http.StatusConflict).
			WithResponseBody("application/json", []byte(`{"message":"[resources-userstore-2] username already exists"}`)).
			Run()
	}

	{
		userstore := &mocks.UserStore{}
		userstore.On("CreateUser", mock.Anything, "wpowell", resources.UserAttributes{
			PineappleOnPizza: true,
		}).Return(nil, errors.New("0000", "error"))

		restfultest.TestCase(t, snailmail, "Userstore failed").
			GivenRequest(http.MethodPost, "/users").
			WithResourceMock("userstore", userstore).
			WithRequestBody("application/json", []byte(`{
						"username": "wpowell",
						"pineapple_on_pizza": true
					}`)).
			ExpectResponse(http.StatusInternalServerError).
			WithResponseBody("application/json", []byte(`{"message":"[0000] error"}`)).
			Run()
	}
}
