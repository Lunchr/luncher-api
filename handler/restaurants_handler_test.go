package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/db/model"
	. "github.com/Lunchr/luncher-api/handler"
	"github.com/Lunchr/luncher-api/router"
	"github.com/Lunchr/luncher-api/session"
	"github.com/Lunchr/luncher-api/storage"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RestaurantsHandlers", func() {
	Describe("GET /restaurants", func() {
		var (
			mockRestaurantsCollection db.Restaurants
			handler                   router.Handler
		)

		BeforeEach(func() {
			mockRestaurantsCollection = &mockRestaurants{}
		})

		JustBeforeEach(func() {
			handler = Restaurants(mockRestaurantsCollection)
		})

		It("should succeed", func(done Done) {
			defer close(done)
			err := handler(responseRecorder, request)
			Expect(err).To(BeNil())
		})

		It("should return json", func(done Done) {
			defer close(done)
			handler(responseRecorder, request)
			contentTypes := responseRecorder.HeaderMap["Content-Type"]
			Expect(contentTypes).To(HaveLen(1))
			Expect(contentTypes[0]).To(Equal("application/json"))
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
					nil,
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
					nil,
				}
			})

			It("should return error 500", func(done Done) {
				defer close(done)
				err := handler(responseRecorder, request)
				Expect(err.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	Describe("GET /restaurant", func() {
		var (
			sessionManager            session.Manager
			mockRestaurantsCollection db.Restaurants
			mockUsersCollection       db.Users
			handler                   router.Handler
		)

		BeforeEach(func() {
			mockRestaurantsCollection = &mockRestaurants{}
		})

		JustBeforeEach(func() {
			handler = Restaurant(mockRestaurantsCollection, sessionManager, mockUsersCollection)
		})

		ExpectUserToBeLoggedIn(func() *router.HandlerError {
			return handler(responseRecorder, request)
		}, func(mgr session.Manager, users db.Users) {
			sessionManager = mgr
			mockUsersCollection = users
		})

		Context("with user logged in", func() {
			BeforeEach(func() {
				sessionManager = &mockSessionManager{isSet: true, id: "correctSession"}
				mockUsersCollection = mockUsers{}
			})

			It("should succeed", func(done Done) {
				defer close(done)
				err := handler(responseRecorder, request)
				Expect(err).To(BeNil())
			})

			It("should return json", func(done Done) {
				defer close(done)
				handler(responseRecorder, request)
				contentTypes := responseRecorder.HeaderMap["Content-Type"]
				Expect(contentTypes).To(HaveLen(1))
				Expect(contentTypes[0]).To(Equal("application/json"))
			})

			It("should include the updated offer in the response", func(done Done) {
				defer close(done)
				handler(responseRecorder, request)
				var restaurant *model.Restaurant
				json.Unmarshal(responseRecorder.Body.Bytes(), &restaurant)
				Expect(restaurant.Name).To(Equal("Asian Chef"))
			})
		})
	})

	Describe("GET /restaurant/offers", func() {
		var (
			sessionManager            session.Manager
			mockRestaurantsCollection db.Restaurants
			mockUsersCollection       db.Users
			handler                   router.Handler
			mockOffersCollection      db.Offers
			imageStorage              storage.Images
		)

		BeforeEach(func() {
			mockRestaurantsCollection = &mockRestaurants{}
			imageStorage = &mockImageStorage{}
		})

		JustBeforeEach(func() {
			handler = RestaurantOffers(mockRestaurantsCollection, sessionManager, mockUsersCollection, mockOffersCollection, imageStorage)
		})

		ExpectUserToBeLoggedIn(func() *router.HandlerError {
			return handler(responseRecorder, request)
		}, func(mgr session.Manager, users db.Users) {
			sessionManager = mgr
			mockUsersCollection = users
		})

		Context("with user logged in", func() {
			BeforeEach(func() {
				sessionManager = &mockSessionManager{isSet: true, id: "correctSession"}
				mockUsersCollection = mockUsers{}
				mockOffersCollection = mockOffers{}
			})

			It("should succeed", func(done Done) {
				defer close(done)
				err := handler(responseRecorder, request)
				Expect(err).To(BeNil())
			})

			It("should return json", func(done Done) {
				defer close(done)
				handler(responseRecorder, request)
				contentTypes := responseRecorder.HeaderMap["Content-Type"]
				Expect(contentTypes).To(HaveLen(1))
				Expect(contentTypes[0]).To(Equal("application/json"))
			})

			It("should include the offers in the response", func(done Done) {
				defer close(done)
				handler(responseRecorder, request)
				var result []*model.Offer
				json.Unmarshal(responseRecorder.Body.Bytes(), &result)
				Expect(result).To(HaveLen(2))
				Expect(result[0].Title).To(Equal("a"))
				Expect(result[1].Title).To(Equal("b"))
				Expect(result[0].Image).To(Equal("images/a large image path"))
				Expect(result[1].Image).To(Equal(""))
			})
		})
	})
})

type mockRestaurants struct {
	getFunc func() ([]*model.Restaurant, error)
	db.Restaurants
}

func (mock mockRestaurants) Get() (restaurants []*model.Restaurant, err error) {
	if mock.getFunc != nil {
		restaurants, err = mock.getFunc()
	}
	return
}

func (c mockOffers) GetForRestaurant(restaurant string, startTime time.Time) ([]*model.Offer, error) {
	Expect(restaurant).To(Equal("Asian Chef"))
	Expect(startTime.Sub(time.Now())).To(BeNumerically("~", 0, time.Second))
	return []*model.Offer{
		&model.Offer{
			Title: "a",
			Image: "image checksum",
		},
		&model.Offer{
			Title: "b",
		},
	}, nil
}
