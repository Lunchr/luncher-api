package mocks

import "github.com/Lunchr/luncher-api/db"
import "github.com/stretchr/testify/mock"

import "github.com/Lunchr/luncher-api/db/model"

type Regions struct {
	mock.Mock
}

func (_m *Regions) Insert(_a0 ...*model.Region) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(...*model.Region) error); ok {
		r0 = rf(_a0...)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *Regions) GetName(_a0 string) (*model.Region, error) {
	ret := _m.Called(_a0)

	var r0 *model.Region
	if rf, ok := ret.Get(0).(func(string) *model.Region); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Region)
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
func (_m *Regions) GetAll() db.RegionIter {
	ret := _m.Called()

	var r0 db.RegionIter
	if rf, ok := ret.Get(0).(func() db.RegionIter); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(db.RegionIter)
	}

	return r0
}
func (_m *Regions) UpdateName(_a0 string, _a1 *model.Region) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, *model.Region) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
