package mocks

import "github.com/Lunchr/luncher-api/db"
import "github.com/stretchr/testify/mock"

import "github.com/Lunchr/luncher-api/db/model"

import "gopkg.in/mgo.v2/bson"

type Restaurants struct {
	mock.Mock
}

func (_m *Restaurants) Insert(_a0 ...*model.Restaurant) ([]*model.Restaurant, error) {
	ret := _m.Called(_a0)

	var r0 []*model.Restaurant
	if rf, ok := ret.Get(0).(func(...*model.Restaurant) []*model.Restaurant); ok {
		r0 = rf(_a0...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Restaurant)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(...*model.Restaurant) error); ok {
		r1 = rf(_a0...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Restaurants) Get() ([]*model.Restaurant, error) {
	ret := _m.Called()

	var r0 []*model.Restaurant
	if rf, ok := ret.Get(0).(func() []*model.Restaurant); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Restaurant)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func() error); ok {
		r1 = rf()
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Restaurants) GetAll() db.RestaurantIter {
	ret := _m.Called()

	var r0 db.RestaurantIter
	if rf, ok := ret.Get(0).(func() db.RestaurantIter); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(db.RestaurantIter)
	}

	return r0
}
func (_m *Restaurants) GetID(_a0 bson.ObjectId) (*model.Restaurant, error) {
	ret := _m.Called(_a0)

	var r0 *model.Restaurant
	if rf, ok := ret.Get(0).(func(bson.ObjectId) *model.Restaurant); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Restaurant)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bson.ObjectId) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Restaurants) Exists(name string) (bool, error) {
	ret := _m.Called(name)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(name)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(name)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Restaurants) UpdateID(_a0 bson.ObjectId, _a1 *model.Restaurant) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(bson.ObjectId, *model.Restaurant) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
