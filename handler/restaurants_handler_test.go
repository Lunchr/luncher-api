package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"golang.org/x/oauth2"

	"gopkg.in/mgo.v2"
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
	Describe("GET /user/restaurants", func() {
		var (
			sessionManager session.Manager
			users          db.Users

			restaurants *mocks.Restaurants
			handler     router.Handler
		)

		JustBeforeEach(func() {
			handler = UserRestaurants(restaurants, sessionManager, users)
		})

		ExpectUserToBeLoggedIn(func() *router.HandlerError {
			return handler(responseRecorder, request)
		}, func(mgr session.Manager, _users_ db.Users) {
			sessionManager = mgr
			users = _users_
		})

		Context("with the user logged in", func() {
			var (
				mockSessionManager *mocks.Manager
				mockUsers          *mocks.Users
			)

			BeforeEach(func() {
				mockSessionManager = new(mocks.Manager)
				sessionManager = mockSessionManager
				mockUsers = new(mocks.Users)
				users = mockUsers
				restaurants = new(mocks.Restaurants)

				allRestaurants := []*model.Restaurant{&model.Restaurant{
					ID: bson.NewObjectId(),
				}, &model.Restaurant{
					FacebookPageID: "fbpageid1",
				}, &model.Restaurant{
					ID: bson.NewObjectId(),
				}, &model.Restaurant{
					FacebookPageID: "fbpageid2",
				}}
				fbUserToken := oauth2.Token{
					AccessToken: "usertoken",
				}
				user := &model.User{
					RestaurantIDs: []bson.ObjectId{allRestaurants[0].ID, allRestaurants[2].ID},
					Session: model.UserSession{
						FacebookUserToken: fbUserToken,
						FacebookPageTokens: []model.FacebookPageToken{model.FacebookPageToken{
							PageID: "fbpageid1",
						}, model.FacebookPageToken{
							PageID: "fbpageid2",
						}},
					},
				}

				mockSessionManager.On("Get", request).Return("session", nil)
				mockUsers.On("GetSessionID", "session").Return(user, nil)
				restaurants.On("GetByIDs", user.RestaurantIDs).Return([]*model.Restaurant{allRestaurants[0], allRestaurants[2]}, nil)
				restaurants.On("GetByFacebookPageIDs", []string{"fbpageid1", "fbpageid2"}).Return([]*model.Restaurant{allRestaurants[1], allRestaurants[3]}, nil)
			})

			It("succeeds", func() {
				err := handler(responseRecorder, request)
				Expect(err).To(BeNil())
			})

			It("returns json", func() {
				handler(responseRecorder, request)
				contentTypes := responseRecorder.HeaderMap["Content-Type"]
				Expect(contentTypes).To(HaveLen(1))
				Expect(contentTypes[0]).To(Equal("application/json"))
			})

			It("includes all restaurants", func() {
				handler(responseRecorder, request)
				var response []*model.Restaurant
				json.Unmarshal(responseRecorder.Body.Bytes(), &response)
				Expect(response).To(HaveLen(4))
			})
		})
	})

	Describe("POST /restaurants", func() {
		var (
			sessionManager        session.Manager
			restaurantsCollection db.Restaurants
			usersCollection       db.Users
			handler               router.Handler
			fbAuth                *mocks.Authenticator
		)

		JustBeforeEach(func() {
			handler = PostRestaurants(restaurantsCollection, sessionManager, usersCollection, fbAuth)
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
				fbAuth = new(mocks.Authenticator)

				user = &model.User{
					Session: model.UserSession{
						FacebookUserToken: oauth2.Token{
							AccessToken: "athing",
						},
					},
				}
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

			Describe("the updated user", func() {
				var updatedUser *model.User
				Context("with restaurant not being attached to a FB page", func() {
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

				Context("with restaurant being attached to a FB page", func() {
					var (
						facebookPageID  = "something"
						pageAccessToken = "access token"
					)

					BeforeEach(func() {
						mockRestaurantsCollection.On("Insert", mock.AnythingOfType("[]*model.Restaurant")).Return([]*model.Restaurant{
							&model.Restaurant{
								ID:             id,
								FacebookPageID: facebookPageID,
							},
						}, nil)
						mockUsersCollection.On("Update", mock.AnythingOfType("string"), mock.AnythingOfType("*model.User")).Return(nil).Run(func(args mock.Arguments) {
							updatedUser = args.Get(1).(*model.User)
						})
						fbAuth.On("PageAccessToken", &user.Session.FacebookUserToken, facebookPageID).Return(pageAccessToken, nil)
					})

					It("should update the user to add a page access token for the restaurant", func() {
						err := handler(responseRecorder, request)
						Expect(err).To(BeNil())
						insertedToken := updatedUser.Session.FacebookPageTokens[0]
						Expect(insertedToken.PageID).To(Equal(facebookPageID))
						Expect(insertedToken.Token).To(Equal(pageAccessToken))
					})
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
				facebookPageID            = "a facebook page ID"
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
					ID:             restaurantID,
					Name:           "restname",
					FacebookPageID: facebookPageID,
				}

				mockSessionManager.On("Get", mock.Anything).Return("session", nil)
			})

			Context("with an invalid restaurant ID", func() {
				BeforeEach(func() {
					mockUsersCollection.On("GetSessionID", "session").Return(&model.User{}, nil)
					params = httprouter.Params{httprouter.Param{
						Key:   "restaurantID",
						Value: "gibberish",
					}}
				})

				It("fails", func() {
					err := handler(responseRecorder, request, params)
					Expect(err).To(HaveOccurred())
					Expect(err.Code).To(Equal(http.StatusBadRequest))
				})
			})

			Context("with a valid restaurant ID", func() {
				BeforeEach(func() {
					mockRestaurantsCollection.On("GetID", restaurant.ID).Return(restaurant, nil)
					params = httprouter.Params{httprouter.Param{
						Key:   "restaurantID",
						Value: restaurant.ID.Hex(),
					}}
				})

				Context("but not authorized", func() {
					BeforeEach(func() {
						user := &model.User{}
						mockUsersCollection.On("GetSessionID", "session").Return(user, nil)
					})

					It("fails", func() {
						err := handler(responseRecorder, request, params)
						Expect(err).To(HaveOccurred())
						Expect(err.Code).To(Equal(http.StatusForbidden))
					})
				})

				Context("having a FB page access token for the restaurant's page", func() {
					BeforeEach(func() {
						user := &model.User{
							RestaurantIDs: []bson.ObjectId{},
							Session: model.UserSession{
								FacebookPageTokens: []model.FacebookPageToken{model.FacebookPageToken{
									PageID: facebookPageID,
								}},
							},
						}
						mockUsersCollection.On("GetSessionID", "session").Return(user, nil)
					})

					It("succeeds", func() {
						err := handler(responseRecorder, request, params)
						Expect(err).To(BeNil())
					})
				})

				Context("having direct access", func() {
					BeforeEach(func() {
						user := &model.User{
							RestaurantIDs: []bson.ObjectId{restaurant.ID},
						}
						mockUsersCollection.On("GetSessionID", "session").Return(user, nil)
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
	})

	Describe("GET /restaurants/:id/offers", func() {
		var (
			sessionManager            session.Manager
			mockRestaurantsCollection db.Restaurants
			mockUsersCollection       db.Users
			params                    httprouter.Params
			handler                   router.HandlerWithParams
			offersCollection          db.Offers
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
				offersCollection, imageStorage, regionsCollection)
		})

		ExpectUserToBeLoggedIn(func() *router.HandlerError {
			return handler(responseRecorder, request, nil)
		}, func(mgr session.Manager, users db.Users) {
			sessionManager = mgr
			mockUsersCollection = users
		})

		Context("with user logged in", func() {
			var restaurantID = bson.ObjectId("12letrrestid")

			BeforeEach(func() {
				sessionManager = &mockSessionManager{isSet: true, id: "correctSession"}
				mockUsersCollection = mockUsers{}
				mockRestaurantsCollection = &mockRestaurants{}
				regionsCollection = new(mocks.Regions)
				regionsCollection.On("GetName", "Tartu").Return(&model.Region{
					Location: "Europe/Tallinn",
				}, nil)
				params = httprouter.Params{httprouter.Param{
					Key:   "restaurantID",
					Value: restaurantID.Hex(),
				}}
			})

			Context("with title not specified", func() {
				BeforeEach(func() {
					requestQuery = url.Values{}
					offersCollection = mockOffers{}
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

			Context("with title specified", func() {
				var (
					title                = "a title"
					mockOffersCollection *mocks.Offers
				)

				BeforeEach(func() {
					mockOffersCollection = new(mocks.Offers)
					offersCollection = mockOffersCollection
					requestQuery = url.Values{
						"title": {url.QueryEscape(title)},
					}
				})

				Context("with no matching offer found", func() {
					BeforeEach(func() {
						mockOffersCollection.On("GetForRestaurantByTitle", restaurantID, title).Return(nil, mgo.ErrNotFound)
					})

					It("should respond with StatusNotFound", func() {
						err := handler(responseRecorder, request, params)
						Expect(err.Code).To(Equal(http.StatusNotFound))
					})
				})

				Context("with DB request failing", func() {
					BeforeEach(func() {
						mockOffersCollection.On("GetForRestaurantByTitle", restaurantID, title).Return(nil, errors.New("something went wrong"))
					})

					It("should respond with StatusInternalServerError", func() {
						err := handler(responseRecorder, request, params)
						Expect(err.Code).To(Equal(http.StatusInternalServerError))
					})
				})

				Context("with DB returning an offer", func() {
					BeforeEach(func() {
						offer := &model.Offer{
							CommonOfferFields: model.CommonOfferFields{
								Title: title,
							},
							ImageChecksum: "image checksum",
						}
						mockOffersCollection.On("GetForRestaurantByTitle", restaurantID, title).Return(offer, nil)
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

					It("should respond with the offer", func() {
						handler(responseRecorder, request, params)
						var result *model.OfferJSON
						json.Unmarshal(responseRecorder.Body.Bytes(), &result)
						Expect(result.Title).To(Equal(title))
						Expect(result.Image.Large).To(Equal("images/a large image path"))
					})
				})
			})
		})
	})

	Describe("GET /restaurants/:id/offer_suggestions", func() {
		var (
			sessionManager            session.Manager
			mockRestaurantsCollection db.Restaurants
			mockUsersCollection       db.Users
			params                    httprouter.Params
			handler                   router.HandlerWithParams
			offersCollection          *mocks.Offers
		)

		BeforeEach(func() {
			mockRestaurantsCollection = &mockRestaurants{}
		})

		JustBeforeEach(func() {
			handler = RestaurantOfferSuggestions(mockRestaurantsCollection, sessionManager, mockUsersCollection,
				offersCollection)
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
				offersCollection = new(mocks.Offers)
				mockRestaurantsCollection = &mockRestaurants{}
				params = httprouter.Params{httprouter.Param{
					Key:   "restaurantID",
					Value: bson.ObjectId("12letrrestid").Hex(),
				}}
			})

			Context("without the 'title' parameter", func() {
				BeforeEach(func() {
					requestQuery = url.Values{}
					offersCollection.On("GetSimilarTitlesForRestaurant", bson.ObjectId("12letrrestid"), "rubbish").Return([]string{}, nil)
				})

				It("fails with StatusBadRequest", func() {
					err := handler(responseRecorder, request, params)
					Expect(err.Code).To(Equal(http.StatusBadRequest))
				})
			})

			Context("with title specified, but no matching results", func() {
				BeforeEach(func() {
					requestQuery = url.Values{
						"title": {"rubbish"},
					}
					offersCollection.On("GetSimilarTitlesForRestaurant", bson.ObjectId("12letrrestid"), "rubbish").Return([]string{}, nil)
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

				It("should return an empty list", func() {
					handler(responseRecorder, request, params)
					var result []string
					json.Unmarshal(responseRecorder.Body.Bytes(), &result)
					Expect(result).To(BeEmpty())
				})
			})

			Context("with title specified and matching results", func() {
				var (
					title1 = "Sweet & Sour Chicken"
					title2 = "Chicken Soup"
				)
				BeforeEach(func() {
					requestQuery = url.Values{
						"title": {"chicken"},
					}
					offersCollection.On("GetSimilarTitlesForRestaurant", bson.ObjectId("12letrrestid"),
						"chicken").Return([]string{title1, title2}, nil)
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

				It("should return the list of matching offer titles", func() {
					handler(responseRecorder, request, params)
					var result []string
					json.Unmarshal(responseRecorder.Body.Bytes(), &result)
					Expect(result).To(HaveLen(2))
					Expect(result[0]).To(Equal(title1))
					Expect(result[1]).To(Equal(title2))
				})
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
