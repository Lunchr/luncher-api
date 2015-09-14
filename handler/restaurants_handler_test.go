package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"gopkg.in/mgo.v2/bson"

	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/db/model"
	. "github.com/Lunchr/luncher-api/handler"
	"github.com/Lunchr/luncher-api/handler/mocks"
	"github.com/Lunchr/luncher-api/router"
	"github.com/Lunchr/luncher-api/session"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/mock"

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

		It("should succeed", func() {
			err := handler(responseRecorder, request)
			Expect(err).To(BeNil())
		})

		It("should return json", func() {
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

			It("should write the returned data to responsewriter", func() {
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

			It("should return error 500", func() {
				err := handler(responseRecorder, request)
				Expect(err.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	Describe("POST /restaurants", func() {
		var (
			sessionManager        session.Manager
			restaurantsCollection db.Restaurants
			usersCollection       db.Users
			handler               router.Handler
		)

		JustBeforeEach(func() {
			handler = PostRestaurants(restaurantsCollection, sessionManager, usersCollection)
		})

		ExpectUserToBeLoggedIn(func() *router.HandlerError {
			return handler(responseRecorder, request)
		}, func(mgr session.Manager, users db.Users) {
			sessionManager = mgr
			usersCollection = users
		})

		Context("with session set and a matching user in DB", func() {
			var (
				mockSessionManager        *mocks.Manager
				mockRestaurantsCollection *mocks.Restaurants
				mockUsersCollection       *mocks.Users
				user                      *model.User
				id                        bson.ObjectId
			)

			BeforeEach(func() {
				mockSessionManager = new(mocks.Manager)
				sessionManager = mockSessionManager
				mockRestaurantsCollection = new(mocks.Restaurants)
				restaurantsCollection = mockRestaurantsCollection
				mockUsersCollection = new(mocks.Users)
				usersCollection = mockUsersCollection
				user = &model.User{}
				id = bson.NewObjectId()

				mockSessionManager.On("Get", mock.Anything).Return("session", nil)
				mockUsersCollection.On("GetSessionID", "session").Return(user, nil)

				requestMethod = "POST"
				requestData = map[string]interface{}{
					"facebook_page_id": "1337",
					"name":             "A Restaurant Name",
					"address":          "Street 10, City, Country",
					"phone":            "+372 1234567890",
					"website":          "https://some.address.com/some/path",
					"email":            "an.email@address.com",
					"location": map[string]interface{}{
						"type":        "Point",
						"coordinates": []float64{12.34, 56.78},
					},
				}
				// Temorary:
				mockUsersCollection.On("UnsetSessionID", mock.AnythingOfType("bson.ObjectId")).Return(nil)
			})

			AfterEach(func() {
				mockSessionManager.AssertExpectations(GinkgoT())
				mockRestaurantsCollection.AssertExpectations(GinkgoT())
				mockUsersCollection.AssertExpectations(GinkgoT())
			})

			Context("with DB inserts succeeding", func() {
				BeforeEach(func() {
					mockRestaurantsCollection.On("Insert", mock.AnythingOfType("[]*model.Restaurant")).Return([]*model.Restaurant{
						&model.Restaurant{
							ID: id,
						},
					}, nil)
					mockUsersCollection.On("Update", mock.AnythingOfType("string"), mock.AnythingOfType("*model.User")).Return(nil)
				})

				It("should succeed", func() {
					err := handler(responseRecorder, request)
					Expect(err).To(BeNil())
				})

				It("should return json", func() {
					handler(responseRecorder, request)
					contentTypes := responseRecorder.HeaderMap["Content-Type"]
					Expect(contentTypes).To(HaveLen(1))
					Expect(contentTypes[0]).To(Equal("application/json"))
				})

				It("should include the restaurant with the new ID", func() {
					handler(responseRecorder, request)
					var restaurant *model.Restaurant
					json.Unmarshal(responseRecorder.Body.Bytes(), &restaurant)
					Expect(restaurant.ID).To(Equal(id))
				})
			})

			Context("the inserted restaurant", func() {
				var insertedRestaurant *model.Restaurant
				BeforeEach(func() {
					mockRestaurantsCollection.On("Insert", mock.AnythingOfType("[]*model.Restaurant")).Return([]*model.Restaurant{
						&model.Restaurant{
							ID: id,
						},
					}, nil).Run(func(args mock.Arguments) {
						insertedRestaurant = args.Get(0).([]*model.Restaurant)[0]
					})
					mockUsersCollection.On("Update", mock.AnythingOfType("string"), mock.AnythingOfType("*model.User")).Return(nil)
				})

				It("should correctly parse and insert the restaurant", func() {
					handler(responseRecorder, request)

					Expect(insertedRestaurant.FacebookPageID).To(Equal("1337"))
					Expect(insertedRestaurant.Name).To(Equal("A Restaurant Name"))
					Expect(insertedRestaurant.Address).To(Equal("Street 10, City, Country"))
					Expect(insertedRestaurant.Region).To(Equal("Tallinn"))
					Expect(insertedRestaurant.Phone).To(Equal("+372 1234567890"))
					Expect(insertedRestaurant.Website).To(Equal("https://some.address.com/some/path"))
					Expect(insertedRestaurant.Email).To(Equal("an.email@address.com"))
					Expect(insertedRestaurant.Location.Type).To(Equal("Point"))
					Expect(insertedRestaurant.Location.Coordinates[0]).To(Equal(12.34))
					Expect(insertedRestaurant.Location.Coordinates[1]).To(Equal(56.78))
				})
			})

			Context("the updated user", func() {
				var updatedUser *model.User
				BeforeEach(func() {
					mockRestaurantsCollection.On("Insert", mock.AnythingOfType("[]*model.Restaurant")).Return([]*model.Restaurant{
						&model.Restaurant{
							ID: id,
						},
					}, nil)
					mockUsersCollection.On("Update", mock.AnythingOfType("string"), mock.AnythingOfType("*model.User")).Return(nil).Run(func(args mock.Arguments) {
						updatedUser = args.Get(1).(*model.User)
					})
				})

				It("should update the user to include a reference to the restaurant", func() {
					handler(responseRecorder, request)
					Expect(updatedUser.RestaurantIDs[0]).To(Equal(id))
				})
			})
		})
	})

	Describe("GET /restaurants/:id", func() {
		var (
			sessionManager        session.Manager
			restaurantsCollection db.Restaurants
			usersCollection       db.Users
			params                httprouter.Params
			handler               router.HandlerWithParams
		)

		JustBeforeEach(func() {
			handler = Restaurant(restaurantsCollection, sessionManager, usersCollection)
		})

		ExpectUserToBeLoggedIn(func() *router.HandlerError {
			return handler(responseRecorder, request, nil)
		}, func(mgr session.Manager, users db.Users) {
			sessionManager = mgr
			usersCollection = users
		})

		Context("with user logged in", func() {
			var (
				mockSessionManager        *mocks.Manager
				mockRestaurantsCollection *mocks.Restaurants
				mockUsersCollection       *mocks.Users
				restaurant                *model.Restaurant
			)

			BeforeEach(func() {
				mockSessionManager = new(mocks.Manager)
				sessionManager = mockSessionManager
				mockRestaurantsCollection = new(mocks.Restaurants)
				restaurantsCollection = mockRestaurantsCollection
				mockUsersCollection = new(mocks.Users)
				usersCollection = mockUsersCollection

				restaurantID := bson.NewObjectId()
				restaurant = &model.Restaurant{
					ID:   restaurantID,
					Name: "restname",
				}

				mockSessionManager.On("Get", mock.Anything).Return("session", nil)
			})

			Context("with an invalid restaurant ID", func() {
				BeforeEach(func() {
					mockUsersCollection.On("GetSessionID", "session").Return(&model.User{}, nil)
					params = httprouter.Params{httprouter.Param{
						Key:   "id",
						Value: "gibberish",
					}}
				})

				It("fails", func() {
					err := handler(responseRecorder, request, params)
					Expect(err).To(HaveOccurred())
					Expect(err.Code).To(Equal(http.StatusBadRequest))
				})
			})

			Context("with a valid, but not authorized, restaurant ID", func() {
				BeforeEach(func() {
					user := &model.User{
						RestaurantIDs: []bson.ObjectId{},
					}
					mockUsersCollection.On("GetSessionID", "session").Return(user, nil)
					mockRestaurantsCollection.On("GetID", restaurant.ID).Return(restaurant, nil)
					params = httprouter.Params{httprouter.Param{
						Key:   "id",
						Value: restaurant.ID.Hex(),
					}}
				})

				It("fails", func() {
					err := handler(responseRecorder, request, params)
					Expect(err).To(HaveOccurred())
					Expect(err.Code).To(Equal(http.StatusForbidden))
				})
			})

			Context("with a valid restaurant ID", func() {
				BeforeEach(func() {
					user := &model.User{
						RestaurantIDs: []bson.ObjectId{restaurant.ID},
					}
					mockUsersCollection.On("GetSessionID", "session").Return(user, nil)
					mockRestaurantsCollection.On("GetID", restaurant.ID).Return(restaurant, nil)
					params = httprouter.Params{httprouter.Param{
						Key:   "id",
						Value: restaurant.ID.Hex(),
					}}
				})

				It("should succeed", func() {
					err := handler(responseRecorder, request, params)
					Expect(err).To(BeNil())
				})

				It("should return json", func() {
					handler(responseRecorder, request, params)
					contentTypes := responseRecorder.HeaderMap["Content-Type"]
					Expect(contentTypes).To(HaveLen(1))
					Expect(contentTypes[0]).To(Equal("application/json"))
				})

				It("should include the restaurant data in the response", func() {
					handler(responseRecorder, request, params)
					var response *model.Restaurant
					json.Unmarshal(responseRecorder.Body.Bytes(), &response)
					Expect(response.ID).To(Equal(restaurant.ID))
					Expect(response.Name).To(Equal("restname"))
				})
			})
		})
	})

	Describe("GET /restaurants/:id/offers", func() {
		var (
			sessionManager            session.Manager
			mockRestaurantsCollection db.Restaurants
			mockUsersCollection       db.Users
			params                    httprouter.Params
			handler                   router.HandlerWithParams
			mockOffersCollection      db.Offers
			imageStorage              *mocks.Images
			regionsCollection         *mocks.Regions
		)

		BeforeEach(func() {
			mockRestaurantsCollection = &mockRestaurants{}
			imageStorage = new(mocks.Images)
			imageStorage.On("PathsFor", "image checksum").Return(&model.OfferImagePaths{
				Large:     "images/a large image path",
				Thumbnail: "images/thumbnail",
			}, nil)
			imageStorage.On("PathsFor", "").Return(nil, nil)
		})

		JustBeforeEach(func() {
			handler = RestaurantOffers(mockRestaurantsCollection, sessionManager, mockUsersCollection,
				mockOffersCollection, imageStorage, regionsCollection)
		})

		ExpectUserToBeLoggedIn(func() *router.HandlerError {
			return handler(responseRecorder, request, nil)
		}, func(mgr session.Manager, users db.Users) {
			sessionManager = mgr
			mockUsersCollection = users
		})

		Context("with user logged in", func() {
			BeforeEach(func() {
				sessionManager = &mockSessionManager{isSet: true, id: "correctSession"}
				mockUsersCollection = mockUsers{}
				mockOffersCollection = mockOffers{}
				mockRestaurantsCollection = &mockRestaurants{}
				regionsCollection = new(mocks.Regions)
				regionsCollection.On("GetName", "Tartu").Return(&model.Region{
					Location: "Europe/Tallinn",
				}, nil)
				params = httprouter.Params{httprouter.Param{
					Key:   "id",
					Value: bson.ObjectId("12letrrestid").Hex(),
				}}
			})

			It("should succeed", func() {
				err := handler(responseRecorder, request, params)
				Expect(err).To(BeNil())
			})

			It("should return json", func() {
				handler(responseRecorder, request, params)
				contentTypes := responseRecorder.HeaderMap["Content-Type"]
				Expect(contentTypes).To(HaveLen(1))
				Expect(contentTypes[0]).To(Equal("application/json"))
			})

			It("should include the offers in the response", func() {
				handler(responseRecorder, request, params)
				var result []*model.OfferJSON
				json.Unmarshal(responseRecorder.Body.Bytes(), &result)
				Expect(result).To(HaveLen(2))
				Expect(result[0].Title).To(Equal("a"))
				Expect(result[1].Title).To(Equal("b"))
				Expect(result[0].Image.Large).To(Equal("images/a large image path"))
				Expect(result[1].Image).To(BeNil())
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

func (c mockOffers) GetForRestaurant(restaurantID bson.ObjectId, startTime time.Time) ([]*model.Offer, error) {
	Expect(restaurantID).To(Equal(bson.ObjectId("12letrrestid")))
	loc, err := time.LoadLocation("Europe/Tallinn")
	Expect(err).NotTo(HaveOccurred())
	now := time.Now()
	thisMorning := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
	Expect(startTime).To(Equal(thisMorning))
	return []*model.Offer{
		&model.Offer{
			CommonOfferFields: model.CommonOfferFields{
				Title: "a",
			},
			ImageChecksum: "image checksum",
		},
		&model.Offer{
			CommonOfferFields: model.CommonOfferFields{
				Title: "b",
			},
		},
	}, nil
}
