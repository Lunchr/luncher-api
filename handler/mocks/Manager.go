package mocks

import "github.com/stretchr/testify/mock"

import "net/http"

type Manager struct {
	mock.Mock
}

func (_m *Manager) GetOrInit(_a0 http.ResponseWriter, _a1 *http.Request) string {
	ret := _m.Called(_a0, _a1)

	var r0 string
	if rf, ok := ret.Get(0).(func(http.ResponseWriter, *http.Request) string); ok {
		r0 = rf(_a0, _a1)
	} else {
		r0 = ret.Get(0).(string)
	}

	return r0
}
func (_m *Manager) Get(_a0 *http.Request) (string, error) {
	ret := _m.Called(_a0)

	var r0 string
	if rf, ok := ret.Get(0).(func(*http.Request) string); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(string)
	}

	var r1 error
	if rf, ok := ret.Get(1).(func(*http.Request) error); ok {
		r1 = rf(_a0)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}
