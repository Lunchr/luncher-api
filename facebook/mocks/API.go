package mocks

import "github.com/stretchr/testify/mock"

import "github.com/deiwin/facebook/model"

type API struct {
	mock.Mock
}

func (_m *API) Me() (*model.User, error) {
	ret := _m.Called()

	var r0 *model.User
	if rf, ok := ret.Get(0).(func() *model.User); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.User)
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
func (_m *API) Accounts() (*model.Accounts, error) {
	ret := _m.Called()

	var r0 *model.Accounts
	if rf, ok := ret.Get(0).(func() *model.Accounts); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Accounts)
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
func (_m *API) Page(pageID string) (*model.Page, error) {
	ret := _m.Called(pageID)

	var r0 *model.Page
	if rf, ok := ret.Get(0).(func(string) *model.Page); ok {
		r0 = rf(pageID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Page)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(pageID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *API) PagePublish(pageAccessToken string, pageID string, post *model.Post) (*model.Post, error) {
	ret := _m.Called(pageAccessToken, pageID, post)

	var r0 *model.Post
	if rf, ok := ret.Get(0).(func(string, string, *model.Post) *model.Post); ok {
		r0 = rf(pageAccessToken, pageID, post)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Post)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string, *model.Post) error); ok {
		r1 = rf(pageAccessToken, pageID, post)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *API) Post(pageAccessToken string, postID string) (*model.Post, error) {
	ret := _m.Called(pageAccessToken, postID)

	var r0 *model.Post
	if rf, ok := ret.Get(0).(func(string, string) *model.Post); ok {
		r0 = rf(pageAccessToken, postID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*model.Post)
		}
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(string, string) error); ok {
		r1 = rf(pageAccessToken, postID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
func (_m *API) PostUpdate(pageAccessToken string, postID string, post *model.Post) error {
	ret := _m.Called(pageAccessToken, postID, post)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string, *model.Post) error); ok {
		r0 = rf(pageAccessToken, postID, post)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
func (_m *API) PostDelete(pageAccessToken string, postID string) error {
	ret := _m.Called(pageAccessToken, postID)

	var r0 error
	if rf, ok := ret.Get(0).(func(string, string) error); ok {
		r0 = rf(pageAccessToken, postID)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}
