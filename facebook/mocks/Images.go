package mocks

import "github.com/stretchr/testify/mock"

import "image"

import "github.com/Lunchr/luncher-api/db/model"

type Images struct {
	mock.Mock
}

func (_m *Images) ChecksumDataURL(_a0 string) (string, error) {
	ret := _m.Called(_a0)

	var r0 string
	if rf, ok := ret.Get(0).(func(string) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Images) StoreDataURL(_a0 string) error {
	ret := _m.Called(_a0)

	var r0 error
	if rf, ok := ret.Get(0).(func(string) error); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *Images) PathsFor(checksum string) (*model.OfferImagePaths, error) {
	ret := _m.Called(checksum)

	var r0 *model.OfferImagePaths
	if rf, ok := ret.Get(0).(func(string) *model.OfferImagePaths); ok {
		r0 = rf(checksum)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.OfferImagePaths)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(checksum)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Images) HasChecksum(checksum string) (bool, error) {
	ret := _m.Called(checksum)

	var r0 bool
	if rf, ok := ret.Get(0).(func(string) bool); ok {
		r0 = rf(checksum)
	} else {
		r0 = ret.Get(0).(bool)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(checksum)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *Images) GetOriginal(checksum string) (image.Image, error) {
	ret := _m.Called(checksum)

	var r0 image.Image
	if rf, ok := ret.Get(0).(func(string) image.Image); ok {
		r0 = rf(checksum)
	} else {
		r0 = ret.Get(0).(image.Image)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(checksum)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
