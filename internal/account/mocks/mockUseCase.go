package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"waizly/helpers/response"
	"waizly/models"
)

type MockAccount struct {
	mock.Mock
}

func (_m *MockAccount) Register(ctx context.Context, params models.RegisterRequest) response.Response {
	ret := _m.Called(ctx, params)

	var r0 response.Response

	if rf, ok := ret.Get(0).(func(context.Context, models.RegisterRequest) response.Response); ok {
		r0 = rf(ctx, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(response.Response)
		}
	}
	return r0
}

func (_m *MockAccount) Login(ctx context.Context, params models.LoginRequest) (response.Response, models.Token) {
	ret := _m.Called(ctx, params)

	var r0 response.Response
	if rf, ok := ret.Get(0).(func(context.Context, models.LoginRequest) response.Response); ok {
		r0 = rf(ctx, params)
	} else {
		r0 = ret.Get(0).(response.Response)
	}

	var r1 models.Token
	if rf, ok := ret.Get(1).(func(context.Context, models.LoginRequest) models.Token); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Get(1).(models.Token)
	}

	return r0, r1
}

// func (_m *MockAccount) DetailAccount(ctx context.Context, id int64) response.Response {
// 	ret := _m.Called(ctx, id)

// 	var r0 response.Response

// 	if rf, ok := ret.Get(0).(func(context.Context, int64) response.Response); ok {
// 		r0 = rf(ctx, id)
// 	} else {
// 		if ret.Get(0) != nil {
// 			r0 = ret.Get(0).(response.Response)
// 		}
// 	}

// 	return r0
// }

func (_m *MockAccount) DetailAccount(ctx context.Context, id int64) response.Response {
	ret := _m.Called(ctx, id)

	var r0 response.Response

	if rf, ok := ret.Get(0).(func(context.Context, int64) response.Response); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(response.Response)
		}
	}

	return r0
}

func (_m *MockAccount) UpdateAccount(ctx context.Context, id int64, params models.Account) response.Response {
	ret := _m.Called(ctx, id, params)

	var r0 response.Response
	if rf, ok := ret.Get(0).(func(context.Context, int64, models.Account) response.Response); ok {
		r0 = rf(ctx, id, params)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(response.Response)
		}
	}
	return r0
}

func (_m *MockAccount) DeleteAccount(ctx context.Context, id int64) response.Response {
	ret := _m.Called(ctx, id)

	var r0 response.Response
	if rf, ok := ret.Get(0).(func(context.Context, int64) response.Response); ok {
		r0 = rf(ctx, id)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(response.Response)
		}
	}
	return r0
}
