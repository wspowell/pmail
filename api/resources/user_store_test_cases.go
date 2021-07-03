package resources

import (
	"context"
	"reflect"
	"runtime"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	UserStoreTestCases = [...]func(*testing.T, UserStore){
		TestCase_CreateUser_new,
		TestCase_CreateUser_username_conflict,
		TestCase_GetUser_exists,
		TestCase_GetUser_not_exists,
		TestCase_DeleteUser_exists,
		TestCase_DeleteUser_not_exists,
		TestCase_UpdateUser_exists,
		TestCase_UpdateUser_not_exists,
	}
)

func functionName(i interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func RunUserStoreTestCases(t *testing.T, userStore UserStore) {
	for _, testCase := range UserStoreTestCases {
		t.Run(functionName(testCase), func(t *testing.T) {
			testCase(t, userStore)
		})
	}
}

func TestCase_CreateUser_new(t *testing.T, users UserStore) {
	ctx := context.Background()

	userAttributes := UserAttributes{
		Username:         uuid.New().String(),
		PineappleOnPizza: true,
	}
	user1, err := users.CreateUser(ctx, userAttributes)
	assert.Nil(t, err)
	// User Attributes should be the same.
	assert.Equal(t, user1.UserAttributes, userAttributes)

	userAttributes = UserAttributes{
		Username:         uuid.New().String(),
		PineappleOnPizza: true,
	}
	user2, err := users.CreateUser(ctx, userAttributes)
	assert.Nil(t, err)
	// User ID should be some random value other than an auto increment.
	assert.NotEqual(t, 1, user2.Id-user1.Id)
	assert.NotEqual(t, user1.Id, user2.Id)
}

func TestCase_CreateUser_username_conflict(t *testing.T, users UserStore) {
	ctx := context.Background()

	userAttributes := UserAttributes{
		Username:         uuid.New().String(),
		PineappleOnPizza: true,
	}
	_, err := users.CreateUser(ctx, userAttributes)
	assert.Nil(t, err)

	userAttributes = UserAttributes{
		Username:         uuid.New().String(),
		PineappleOnPizza: true,
	}
	user2, err := users.CreateUser(ctx, userAttributes)
	assert.Nil(t, user2)
	assert.ErrorIs(t, err, ErrCreateUserErrorUsernameConflict)
}

func TestCase_GetUser_exists(t *testing.T, users UserStore) {
	ctx := context.Background()

	userAttributes := UserAttributes{
		Username:         uuid.New().String(),
		PineappleOnPizza: true,
	}
	user, err := users.CreateUser(ctx, userAttributes)
	assert.Nil(t, err)

	foundUser, err := users.GetUser(ctx, user.Id)
	assert.Nil(t, err)
	assert.Equal(t, user.Id, foundUser.Id)
}

func TestCase_GetUser_not_exists(t *testing.T, users UserStore) {
	ctx := context.Background()

	foundUser, err := users.GetUser(ctx, UserId(uuid.New().ID()))
	assert.Nil(t, foundUser)
	assert.ErrorIs(t, err, ErrGetUserErrorUserNotFound)
}

func TestCase_DeleteUser_exists(t *testing.T, users UserStore) {
	ctx := context.Background()

	userAttributes := UserAttributes{
		Username:         uuid.New().String(),
		PineappleOnPizza: true,
	}
	user, err := users.CreateUser(ctx, userAttributes)
	assert.Nil(t, err)

	err = users.DeleteUser(ctx, user.Id)
	assert.Nil(t, err)

	foundUser, err := users.GetUser(ctx, user.Id)
	assert.Nil(t, err)
	assert.Nil(t, foundUser.Id)
}

func TestCase_DeleteUser_not_exists(t *testing.T, users UserStore) {
	ctx := context.Background()

	userId := UserId(uuid.New().ID())
	err := users.DeleteUser(ctx, userId)
	assert.Nil(t, err)

	foundUser, err := users.GetUser(ctx, userId)
	assert.Nil(t, err)
	assert.Nil(t, foundUser.Id)
}

func TestCase_UpdateUser_exists(t *testing.T, users UserStore) {
	ctx := context.Background()

	userAttributes := UserAttributes{
		Username:         uuid.New().String(),
		PineappleOnPizza: true,
	}
	user, err := users.CreateUser(ctx, userAttributes)
	assert.Nil(t, err)

	newUserAttributes := UserAttributes{
		Username:         uuid.New().String(),
		PineappleOnPizza: false,
	}
	err = users.UpdateUser(ctx, user.Id, newUserAttributes)
	assert.Nil(t, err)

	foundUser, err := users.GetUser(ctx, user.Id)
	assert.Nil(t, err)
	assert.Equal(t, user.Id, foundUser.Id)
	assert.Equal(t, newUserAttributes, foundUser.UserAttributes)
}

func TestCase_UpdateUser_not_exists(t *testing.T, users UserStore) {
	ctx := context.Background()

	newUserAttributes := UserAttributes{
		Username:         uuid.New().String(),
		PineappleOnPizza: false,
	}
	err := users.UpdateUser(ctx, UserId(uuid.New().ID()), newUserAttributes)
	assert.Nil(t, err)
}
