package mocks

import "github.com/stretchr/testify/mock"

import "time"
import "github.com/Lunchr/luncher-api/db/model"
import "github.com/Lunchr/luncher-api/geo"

import "gopkg.in/mgo.v2/bson"

type Offers struct {
	mock.Mock
}

func (_m *Offers) Insert(_a0 ...*model.Offer) ([]*model.Offer, error) {
	ret := _m.Called(_a0)

	var r0 []*model.Offer
	if rf, ok := ret.Get(0).(func(...*model.Offer) []*model.Offer); ok {
		r0 = rf(_a0...)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Offer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(...*model.Offer) error); ok {
		r1 = rf(_a0...)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Offers) GetForRegion(region string, startTime time.Time, endTime time.Time) ([]*model.Offer, error) {
	ret := _m.Called(region, startTime, endTime)

	var r0 []*model.Offer
	if rf, ok := ret.Get(0).(func(string, time.Time, time.Time) []*model.Offer); ok {
		r0 = rf(region, startTime, endTime)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Offer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, time.Time, time.Time) error); ok {
		r1 = rf(region, startTime, endTime)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Offers) GetNear(loc geo.Location, startTime time.Time, endTime time.Time) ([]*model.OfferWithDistance, error) {
	ret := _m.Called(loc, startTime, endTime)

	var r0 []*model.OfferWithDistance
	if rf, ok := ret.Get(0).(func(geo.Location, time.Time, time.Time) []*model.OfferWithDistance); ok {
		r0 = rf(loc, startTime, endTime)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.OfferWithDistance)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(geo.Location, time.Time, time.Time) error); ok {
		r1 = rf(loc, startTime, endTime)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Offers) GetForRestaurant(restaurantID bson.ObjectId, startTime time.Time) ([]*model.Offer, error) {
	ret := _m.Called(restaurantID, startTime)

	var r0 []*model.Offer
	if rf, ok := ret.Get(0).(func(bson.ObjectId, time.Time) []*model.Offer); ok {
		r0 = rf(restaurantID, startTime)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Offer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bson.ObjectId, time.Time) error); ok {
		r1 = rf(restaurantID, startTime)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Offers) GetSimilarTitlesForRestaurant(restaurantID bson.ObjectId, partialTitle string) ([]string, error) {
	ret := _m.Called(restaurantID, partialTitle)

	var r0 []string
	if rf, ok := ret.Get(0).(func(bson.ObjectId, string) []string); ok {
		r0 = rf(restaurantID, partialTitle)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]string)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bson.ObjectId, string) error); ok {
		r1 = rf(restaurantID, partialTitle)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Offers) GetForRestaurantByTitle(restaurantID bson.ObjectId, title string) (*model.Offer, error) {
	ret := _m.Called(restaurantID, title)

	var r0 *model.Offer
	if rf, ok := ret.Get(0).(func(bson.ObjectId, string) *model.Offer); ok {
		r0 = rf(restaurantID, title)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Offer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bson.ObjectId, string) error); ok {
		r1 = rf(restaurantID, title)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Offers) GetForRestaurantWithinTimeBounds(restaurantID bson.ObjectId, startTime time.Time, endTime time.Time) ([]*model.Offer, error) {
	ret := _m.Called(restaurantID, startTime, endTime)

	var r0 []*model.Offer
	if rf, ok := ret.Get(0).(func(bson.ObjectId, time.Time, time.Time) []*model.Offer); ok {
		r0 = rf(restaurantID, startTime, endTime)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]*model.Offer)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(bson.ObjectId, time.Time, time.Time) error); ok {
		r1 = rf(restaurantID, startTime, endTime)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Offers) UpdateID(_a0 bson.ObjectId, _a1 *model.Offer) error {
	ret := _m.Called(_a0, _a1)

	var r0 error
	if rf, ok := ret.Get(0).(func(bson.ObjectId, *model.Offer) error); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *Offers) GetID(_a0 bson.ObjectId) (*model.Offer, error) {
	ret := _m.Called(_a0)

	var r0 *model.Offer
	if rf, ok := ret.Get(0).(func(bson.ObjectId) *model.Offer); ok {
		r0 = rf(_a0)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Offer)
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
func (_m *Offers) RemoveID(_a0 bson.ObjectId) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(bson.ObjectId) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
