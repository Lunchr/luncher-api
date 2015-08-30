package mocks

import "github.com/stretchr/testify/mock"

import "github.com/Lunchr/luncher-api/db/model"
import "github.com/Lunchr/luncher-api/router"

type Post struct {
	mock.Mock
}

func (_m *Post) Update(_a0 model.DateWithoutTime, _a1 *model.User, _a2 *model.Restaurant) *router.HandlerError {
	ret := _m.Called(_a0, _a1, _a2)

	var r0 *router.HandlerError
	if rf, ok := ret.Get(0).(func(model.DateWithoutTime, *model.User, *model.Restaurant) *router.HandlerError); ok {
		r0 = rf(_a0, _a1, _a2)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*router.HandlerError)
		}
	}

	return r0
}
