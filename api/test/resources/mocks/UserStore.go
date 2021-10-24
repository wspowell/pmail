// Code generated by mockery v0.0.0-dev. DO NOT EDIT.

package mocks

import (
	context "context"

	mock "github.com/stretchr/testify/mock"
)

// UserStore is an autogenerated mock type for the UserStore type
type UserStore struct {
	mock.Mock
}

// CreateUser provides a mock function with given fields: ctx, username, attributes
func (_m *UserStore) CreateUser(ctx context.Context, username string, attributes resources.UserAttributes) (*resources.User, error) {
	ret := _m.Called(ctx, username, attributes)

	var r0 *resources.User
	if rf, ok := ret.Get(0).(func(context.Context, string, resources.UserAttributes) *resources.User); ok {
		r0 = rf(ctx, username, attributes)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*resources.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, string, resources.UserAttributes) error); ok {
		r1 = rf(ctx, username, attributes)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// DeleteUser provides a mock function with given fields: ctx, userId
func (_m *UserStore) DeleteUser(ctx context.Context, userId uint32) error {
	ret := _m.Called(ctx, userId)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32) error); ok {
		r0 = rf(ctx, userId)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// GetUser provides a mock function with given fields: ctx, userId
func (_m *UserStore) GetUser(ctx context.Context, userId uint32) (*resources.User, error) {
	ret := _m.Called(ctx, userId)

	var r0 *resources.User
	if rf, ok := ret.Get(0).(func(context.Context, uint32) *resources.User); ok {
		r0 = rf(ctx, userId)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*resources.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, uint32) error); ok {
		r1 = rf(ctx, userId)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// UpdateUser provides a mock function with given fields: ctx, userId, newAttributes
func (_m *UserStore) UpdateUser(ctx context.Context, userId uint32, newAttributes resources.UserAttributes) error {
	ret := _m.Called(ctx, userId, newAttributes)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, uint32, resources.UserAttributes) error); ok {
		r0 = rf(ctx, userId, newAttributes)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
