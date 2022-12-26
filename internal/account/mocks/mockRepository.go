package mocks

import (
	"context"
	"waizly/models"

	"github.com/stretchr/testify/mock"
)

type AccountRepository struct {
	mock.Mock
}

func (_m *AccountRepository) Create(ctx context.Context, params models.Account) (int64, error) {
	ret := _m.Called(ctx, params)

	var r0 int64
	if rf, ok := ret.Get(0).(func(context.Context, models.Account) int64); ok {
		r0 = rf(ctx, params)
	} else {
		r0 = ret.Get(0).(int64)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(context.Context, models.Account) error); ok {
		r1 = rf(ctx, params)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *AccountRepository) FindByID(ctx context.Context, id int64) (models.Account, error) {
	ret := _m.Called(ctx, id)

	var r0 models.Account
	if rf, ok := ret.Get(0).(func(context.Context, int64) models.Account); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Get(0).(models.Account)
	}

	var r1 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r1 = rf(ctx, id)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *AccountRepository) FindByEmail(ctx context.Context, email string) (models.Account, error) {
	ret := _m.Called(ctx, email)

	var r0 models.Account
	if rf, ok := ret.Get(0).(func(context.Context, string) models.Account); ok {
		r0 = rf(ctx, email)
	} else {
		r0 = ret.Get(0).(models.Account)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(ctx context.Context, email string) error); ok {
		r1 = rf(ctx, email)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

func (_m *AccountRepository) Update(ctx context.Context, id int64, params models.Account) error {
	ret := _m.Called(ctx, id, params)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64, models.Account) error); ok {
		r0 = rf(ctx, id, params)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

func (_m *AccountRepository) Delete(ctx context.Context, id int64) error {
	ret := _m.Called(ctx, id)

	var r0 error
	if rf, ok := ret.Get(0).(func(context.Context, int64) error); ok {
		r0 = rf(ctx, id)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
