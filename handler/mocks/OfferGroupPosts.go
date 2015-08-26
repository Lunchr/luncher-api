package mocks

import "github.com/stretchr/testify/mock"

import "github.com/Lunchr/luncher-api/db/model"

import "gopkg.in/mgo.v2/bson"

type OfferGroupPosts struct {
	mock.Mock
}

func (_m *OfferGroupPosts) Insert(_a0 ...*model.OfferGroupPost) ([]*model.OfferGroupPost, error) {
	ret := _m.Called(_a0)

	var r0 []*model.OfferGroupPost
	if rf, ok := ret.Get(0).(func(...*model.OfferGroupPost) []*model.OfferGroupPost); ok {
		r0 = rf(_a0...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.OfferGroupPost)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(...*model.OfferGroupPost) error); ok {
		r1 = rf(_a0...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *OfferGroupPosts) UpdateByID(_a0 bson.ObjectId, _a1 *model.OfferGroupPost) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(bson.ObjectId, *model.OfferGroupPost) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *OfferGroupPosts) GetByID(_a0 bson.ObjectId) (*model.OfferGroupPost, error) {
	ret := _m.Called(_a0)

	var r0 *model.OfferGroupPost
	if rf, ok := ret.Get(0).(func(bson.ObjectId) *model.OfferGroupPost); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.OfferGroupPost)
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
func (_m *OfferGroupPosts) GetByDate(_a0 model.DateWithoutTime, _a1 bson.ObjectId) (*model.OfferGroupPost, error) {
	ret := _m.Called(_a0, _a1)

	var r0 *model.OfferGroupPost
	if rf, ok := ret.Get(0).(func(model.DateWithoutTime, bson.ObjectId) *model.OfferGroupPost); ok {
		r0 = rf(_a0, _a1)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.OfferGroupPost)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(model.DateWithoutTime, bson.ObjectId) error); ok {
		r1 = rf(_a0, _a1)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
