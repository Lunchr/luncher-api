package mocks

import "github.com/stretchr/testify/mock"

import "github.com/Lunchr/luncher-api/db/model"

type RegistrationAccessTokens struct {
	mock.Mock
}

func (_m *RegistrationAccessTokens) Insert(_a0 *model.RegistrationAccessToken) (*model.RegistrationAccessToken, error) {
	ret := _m.Called(_a0)

	var r0 *model.RegistrationAccessToken
	if rf, ok := ret.Get(0).(func(*model.RegistrationAccessToken) *model.RegistrationAccessToken); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.RegistrationAccessToken)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*model.RegistrationAccessToken) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *RegistrationAccessTokens) Exists(_a0 model.Token) (bool, error) {
	ret := _m.Called(_a0)

	var r0 bool
	if rf, ok := ret.Get(0).(func(model.Token) bool); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(model.Token) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
