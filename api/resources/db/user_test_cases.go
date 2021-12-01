package db

// import (
// 	"reflect"
// 	"runtime"
// 	"testing"

// 	"github.com/google/uuid"
// 	"github.com/stretchr/testify/assert"
// 	"github.com/wspowell/context"
// 	"github.com/wspowell/snailmail/resources/models/user"
// )

// var (
// 	apiTestCases = [...]func(*testing.T, user.Api){
// 		TestCase_CreateUser_new,
// 		TestCase_CreateUser_username_conflict,
// 		TestCase_GetUser_exists,
// 		TestCase_GetUser_not_exists,
// 		TestCase_DeleteUser_exists,
// 		TestCase_DeleteUser_not_exists,
// 		TestCase_UpdateUser_exists,
// 		TestCase_UpdateUser_not_exists,
// 	}
// )

// func functionName(i interface{}) string {
// 	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
// }

// // RunApiTestCases for any given user store implementation.
// // Runs all unit tests that must pass regardless of implementation.
// func RunApiTestCases(t *testing.T, api user.Api) {
// 	t.Parallel()
// 	for _, testCase := range apiTestCases {
// 		t.Run(functionName(testCase), func(t *testing.T) {
// 			t.Parallel()
// 			testCase(t, api)
// 		})
// 	}
// }

// func TestCase_CreateUser_new(t *testing.T, users user.Api) {
// 	ctx := context.Background()

// 	userAttributes := user.Attributes{
// 		Username:         uuid.New().String(),
// 		PineappleOnPizza: true,
// 	}
// 	user1, err := users.CreateUser(ctx, userAttributes)
// 	assert.Nil(t, err)
// 	// User Attributes should be the same.
// 	assert.Equal(t, user1.Attributes, userAttributes)

// 	userAttributes = user.Attributes{
// 		Username:         uuid.New().String(),
// 		PineappleOnPizza: true,
// 	}

// 	user2, err := users.CreateUser(ctx, userAttributes)
// 	assert.Nil(t, err)
// 	assert.NotEqual(t, user1.UserGuid, user2.UserGuid)
// }

// func TestCase_CreateUser_username_conflict(t *testing.T, users user.Api) {
// 	ctx := context.Background()

// 	userAttributes := user.Attributes{
// 		Username:         uuid.New().String(),
// 		PineappleOnPizza: true,
// 	}
// 	_, err := users.CreateUser(ctx, userAttributes)
// 	assert.Nil(t, err)

// 	user2, err := users.CreateUser(ctx, userAttributes)
// 	assert.Nil(t, user2)
// 	assert.ErrorIs(t, err, user.ErrUsernameExists)
// }

// func TestCase_GetUser_exists(t *testing.T, users user.Api) {
// 	ctx := context.Background()

// 	userAttributes := user.Attributes{
// 		Username:         uuid.New().String(),
// 		PineappleOnPizza: true,
// 	}
// 	user, err := users.CreateUser(ctx, userAttributes)
// 	assert.Nil(t, err)

// 	foundUser, err := users.GetUser(ctx, user.UserGuid)
// 	assert.Nil(t, err)
// 	assert.Equal(t, user.UserGuid, foundUser.UserGuid)
// }

// func TestCase_GetUser_not_exists(t *testing.T, users user.Api) {
// 	ctx := context.Background()

// 	foundUser, err := users.GetUser(ctx, user.Guid(uuid.New().String()))
// 	assert.Nil(t, foundUser)
// 	assert.ErrorIs(t, err, user.ErrNotFound)
// }

// func TestCase_DeleteUser_exists(t *testing.T, users user.Api) {
// 	ctx := context.Background()

// 	userAttributes := user.Attributes{
// 		Username:         uuid.New().String(),
// 		PineappleOnPizza: true,
// 	}
// 	createdUser, err := users.CreateUser(ctx, userAttributes)
// 	assert.Nil(t, err)

// 	err = users.DeleteUser(ctx, createdUser.UserGuid)
// 	assert.Nil(t, err)

// 	foundUser, err := users.GetUser(ctx, createdUser.UserGuid)
// 	assert.ErrorIs(t, err, user.ErrNotFound)
// 	assert.Nil(t, foundUser)
// }

// func TestCase_DeleteUser_not_exists(t *testing.T, users user.Api) {
// 	ctx := context.Background()

// 	userId := user.Guid(uuid.New().String())
// 	err := users.DeleteUser(ctx, userId)
// 	assert.Nil(t, err)

// 	foundUser, err := users.GetUser(ctx, userId)
// 	assert.ErrorIs(t, err, user.ErrNotFound)
// 	assert.Nil(t, foundUser)
// }

// func TestCase_UpdateUser_exists(t *testing.T, users user.Api) {
// 	ctx := context.Background()

// 	userAttributes := user.Attributes{
// 		Username:         uuid.New().String(),
// 		PineappleOnPizza: true,
// 	}
// 	userCreated, err := users.CreateUser(ctx, userAttributes)
// 	assert.Nil(t, err)

// 	newUserAttributes := user.Attributes{
// 		Username:         uuid.New().String(),
// 		PineappleOnPizza: false,
// 	}
// 	err = users.UpdateUser(ctx, userCreated.UserGuid, newUserAttributes)
// 	assert.Nil(t, err)

// 	foundUser, err := users.GetUser(ctx, userCreated.UserGuid)
// 	assert.Nil(t, err)
// 	assert.Equal(t, userCreated.UserGuid, foundUser.UserGuid)
// 	assert.Equal(t, newUserAttributes, foundUser.Attributes)
// }

// func TestCase_UpdateUser_not_exists(t *testing.T, users user.Api) {
// 	ctx := context.Background()

// 	newUserAttributes := user.Attributes{
// 		Username:         uuid.New().String(),
// 		PineappleOnPizza: false,
// 	}
// 	err := users.UpdateUser(ctx, user.Guid(uuid.New().String()), newUserAttributes)
// 	assert.ErrorIs(t, err, user.ErrNotFound)
// }
