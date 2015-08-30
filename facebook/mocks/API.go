package mocks

import "github.com/stretchr/testify/mock"

import "github.com/deiwin/facebook/model"

type API struct {
	mock.Mock
}

func (m *API) Me() (*model.User, error) {
	ret := m.Called()

	var r0 *model.User
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.User)
	}
	r1 := ret.Error(1)

	return r0, r1
}
func (m *API) Accounts() (*model.Accounts, error) {
	ret := m.Called()

	var r0 *model.Accounts
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.Accounts)
	}
	r1 := ret.Error(1)

	return r0, r1
}
func (m *API) Page(pageID string) (*model.Page, error) {
	ret := m.Called(pageID)

	var r0 *model.Page
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.Page)
	}
	r1 := ret.Error(1)

	return r0, r1
}
func (m *API) PagePublish(pageAccessToken string, pageID string, message string) (*model.Post, error) {
	ret := m.Called(pageAccessToken, pageID, message)

	var r0 *model.Post
	if ret.Get(0) != nil {
		r0 = ret.Get(0).(*model.Post)
	}
	r1 := ret.Error(1)

	return r0, r1
}
func (m *API) PostDelete(pageAccessToken string, postID string) error {
	ret := m.Called(pageAccessToken, postID)

	r0 := ret.Error(0)

	return r0
}
