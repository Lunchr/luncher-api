package handler_test

import (
	"encoding/json"
	"errors"

	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/db/model"
	. "github.com/Lunchr/luncher-api/handler"
	"github.com/Lunchr/luncher-api/handler/mocks"
	"github.com/Lunchr/luncher-api/router"
	"github.com/Lunchr/luncher-api/session"
	"github.com/julienschmidt/httprouter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ = Describe("OfferGroupPostHandlers", func() {
	Describe("GET /restaurant/post/:date", func() {
		var (
			sessionManager        session.Manager
			postsCollection       db.OfferGroupPosts
			restaurantsCollection db.Restaurants
			usersCollection       db.Users
			handler               router.HandlerWithParams
		)

		JustBeforeEach(func() {
			handler = OfferGroupPost(postsCollection, sessionManager, usersCollection, restaurantsCollection)
		})

		ExpectUserToBeLoggedIn(func() *router.HandlerError {
			return handler(responseRecorder, request, nil)
		}, func(mgr session.Manager, users db.Users) {
			sessionManager = mgr
			usersCollection = users
		})

		Context("with session set and a matching user in DB", func() {
			var (
				mockSessionManager        *mocks.Manager
				mockPostsCollection       *mocks.OfferGroupPosts
				mockRestaurantsCollection *mocks.Restaurants
				mockUsersCollection       *mocks.Users
				params                    httprouter.Params
				restaurantID              bson.ObjectId
			)

			BeforeEach(func() {
				mockSessionManager = new(mocks.Manager)
				sessionManager = mockSessionManager
				mockPostsCollection = new(mocks.OfferGroupPosts)
				postsCollection = mockPostsCollection
				mockRestaurantsCollection = new(mocks.Restaurants)
				restaurantsCollection = mockRestaurantsCollection
				mockUsersCollection = new(mocks.Users)
				usersCollection = mockUsersCollection

				restaurantID = bson.NewObjectId()
				restaurant := &model.Restaurant{
					ID: restaurantID,
				}
				user := &model.User{
					RestaurantIDs: []bson.ObjectId{restaurant.ID},
				}

				mockSessionManager.On("Get", mock.Anything).Return("session", nil)
				mockUsersCollection.On("GetSessionID", "session").Return(user, nil)
				mockRestaurantsCollection.On("GetID", restaurantID).Return(restaurant, nil)

				params = httprouter.Params{httprouter.Param{
					Key:   "date",
					Value: "2015-04-10",
				}}
			})

			Context("with db returning a proper result", func() {
				BeforeEach(func() {
					post := &model.OfferGroupPost{
						MessageTemplate: "this is a message template %%",
					}
					mockPostsCollection.On("GetByDate", model.DateWithoutTime("2015-04-10"), restaurantID).Return(post, nil)
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
					var result *model.OfferGroupPost
					json.Unmarshal(responseRecorder.Body.Bytes(), &result)
					Expect(result.MessageTemplate).To(Equal("this is a message template %%"))
				})
			})

			Context("with db returning an error", func() {
				BeforeEach(func() {
					mockPostsCollection.On("GetByDate", model.DateWithoutTime("2015-04-10"), restaurantID).Return(nil,
						errors.New("idk man, things happened"))
				})

				It("should fail", func() {
					err := handler(responseRecorder, request, params)
					Expect(err).NotTo(BeNil())
				})
			})

			Context("with db returning a NotFound error", func() {
				BeforeEach(func() {
					mockPostsCollection.On("GetByDate", model.DateWithoutTime("2015-04-10"), restaurantID).Return(nil,
						mgo.ErrNotFound)
				})

				It("should fail with 404", func() {
					err := handler(responseRecorder, request, params)
					Expect(err).NotTo(BeNil())
					Expect(err.Code).To(Equal(404))
				})
			})
		})
	})
})
