package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	. "github.com/deiwin/luncher-api/handler"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RestaurantsHandler", func() {
	var (
		mockRestaurantsCollection db.Restaurants
		handler                   Handler
	)

	BeforeEach(func() {
		mockRestaurantsCollection = &mockRestaurants{}
	})

	JustBeforeEach(func() {
		handler = Restaurants(mockRestaurantsCollection)
	})

	Describe("Get", func() {
		It("should succeed", func(done Done) {
			defer close(done)
			handler(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
		})

		It("should return json", func(done Done) {
			defer close(done)
			handler(responseRecorder, request)
			contentTypes := responseRecorder.HeaderMap["Content-Type"]
			Expect(contentTypes).To(HaveLen(1))
			Expect(contentTypes[0]).To(Equal("application/json"))
			// TODO the header assertion could be made a custom matcher
		})

		Context("with simple mocked result from DB", func() {
			var (
				mockResult []*model.Restaurant
			)
			BeforeEach(func() {
				mockResult = []*model.Restaurant{&model.Restaurant{Name: "somerestaurant"}}
				mockRestaurantsCollection = &mockRestaurants{
					func() (restaurants []*model.Restaurant, err error) {
						restaurants = mockResult
						return
					},
				}
			})

			It("should write the returned data to responsewriter", func(done Done) {
				defer close(done)
				handler(responseRecorder, request)
				var result []*model.Restaurant
				json.Unmarshal(responseRecorder.Body.Bytes(), &result)
				Expect(result).To(HaveLen(1))
				Expect(result[0].Name).To(Equal(mockResult[0].Name))
			})
		})

		Context("with an error returned from the DB", func() {
			var dbErr = errors.New("DB stuff failed")

			BeforeEach(func() {
				mockRestaurantsCollection = &mockRestaurants{
					func() (restaurants []*model.Restaurant, err error) {
						err = dbErr
						return
					},
				}
			})

			It("should return error 500", func(done Done) {
				defer close(done)
				handler(responseRecorder, request)
				Expect(responseRecorder.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})

type mockRestaurants struct {
	getFunc func() ([]*model.Restaurant, error)
}

func (mock mockRestaurants) Insert(restaurantsToInsert ...*model.Restaurant) (err error) {
	return
}

func (mock mockRestaurants) Get() (restaurants []*model.Restaurant, err error) {
	if mock.getFunc != nil {
		restaurants, err = mock.getFunc()
	}
	return
}
