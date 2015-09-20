package mocks

import "github.com/Lunchr/luncher-api/db"
import "github.com/stretchr/testify/mock"

import "github.com/Lunchr/luncher-api/db/model"
import "golang.org/x/oauth2"

import "gopkg.in/mgo.v2/bson"

type Users struct {
	mock.Mock
}

func (_m *Users) Insert(_a0 ...*model.User) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(...*model.User) error); ok {
		r0 = rf(_a0...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *Users) GetFbID(_a0 string) (*model.User, error) {
	ret := _m.Called(_a0)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(string) *model.User); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Users) GetSessionID(_a0 string) (*model.User, error) {
	ret := _m.Called(_a0)

	var r0 *model.User
	if rf, ok := ret.Get(0).(func(string) *model.User); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Users) GetAll() db.UserIter {
	ret := _m.Called()

	var r0 db.UserIter
	if rf, ok := ret.Get(0).(func() db.UserIter); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(db.UserIter)
	}

	return r0
}
func (_m *Users) Update(_a0 string, _a1 *model.User) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *model.User) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *Users) SetAccessToken(_a0 string, _a1 oauth2.Token) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, oauth2.Token) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *Users) SetPageAccessTokens(_a0 string, _a1 []model.FacebookPageToken) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, []model.FacebookPageToken) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *Users) SetSessionID(_a0 bson.ObjectId, _a1 string) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(bson.ObjectId, string) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *Users) UnsetSessionID(_a0 bson.ObjectId) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(bson.ObjectId) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
