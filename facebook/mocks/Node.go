package mocks

import "github.com/stretchr/testify/mock"

import "image"
import "image/color"

type Node struct {
	mock.Mock
}

func (_m *Node) Draw(width int, height int) image.Image {
	ret := _m.Called(width, height)

	var r0 image.Image
	if rf, ok := ret.Get(0).(func(int, int) image.Image); ok {
		r0 = rf(width, height)
	} else {
		r0 = ret.Get(0).(image.Image)
	}

	return r0
}
func (_m *Node) DrawWithBorder(width int, height int, borderColor color.Color, borderWidth int) image.Image {
	ret := _m.Called(width, height, borderColor, borderWidth)

	var r0 image.Image
	if rf, ok := ret.Get(0).(func(int, int, color.Color, int) image.Image); ok {
		r0 = rf(width, height, borderColor, borderWidth)
	} else {
		r0 = ret.Get(0).(image.Image)
	}

	return r0
}
