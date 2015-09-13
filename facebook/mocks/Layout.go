package mocks

import "github.com/deiwin/picasso"
import "github.com/stretchr/testify/mock"

import "image"

type Layout struct {
	mock.Mock
}

func (_m *Layout) Compose(_a0 []image.Image) picasso.Node {
	ret := _m.Called(_a0)

	var r0 picasso.Node
	if rf, ok := ret.Get(0).(func([]image.Image) picasso.Node); ok {
		r0 = rf(_a0)
	} else {
		r0 = ret.Get(0).(picasso.Node)
	}

	return r0
}
