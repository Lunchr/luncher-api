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
	Describe("GET /restaurants/:restaurantID/posts/:date", func() {
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
				facebookPageID := "fbpageid"
				restaurant := &model.Restaurant{
					ID:             restaurantID,
					FacebookPageID: facebookPageID,
				}
				user := &model.User{
					Session: model.UserSession{
						FacebookPageTokens: []model.FacebookPageToken{model.FacebookPageToken{
							PageID: facebookPageID,
						}},
					},
				}

				mockSessionManager.On("Get", mock.Anything).Return("session", nil)
				mockUsersCollection.On("GetSessionID", "session").Return(user, nil)
				mockRestaurantsCollection.On("GetID", restaurantID).Return(restaurant, nil)

				params = httprouter.Params{httprouter.Param{
					Key:   "date",
					Value: "2015-04-10",
				}, httprouter.Param{
					Key:   "restaurantID",
					Value: restaurantID.Hex(),
				}}
			})

			AfterEach(func() {
				mockSessionManager.AssertExpectations(GinkgoT())
				mockPostsCollection.AssertExpectations(GinkgoT())
				mockRestaurantsCollection.AssertExpectations(GinkgoT())
				mockUsersCollection.AssertExpectations(GinkgoT())
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

	Describe("POST /restaurants/:restaurantID/posts", func() {
		var (
			sessionManager        session.Manager
			postsCollection       db.OfferGroupPosts
			restaurantsCollection db.Restaurants
			usersCollection       db.Users
			facebookPost          *mocks.Post
			params                httprouter.Params
			handler               router.HandlerWithParams
		)

		JustBeforeEach(func() {
			handler = PostOfferGroupPost(postsCollection, sessionManager, usersCollection, restaurantsCollection, facebookPost)
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
				restaurantID              bson.ObjectId
				id                        bson.ObjectId
				user                      *model.User
				restaurant                *model.Restaurant
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
				facebookPost = new(mocks.Post)

				id = bson.NewObjectId()

				restaurantID = bson.NewObjectId()
				facebookPageID := "fbpageid"
				restaurant = &model.Restaurant{
					ID:             restaurantID,
					FacebookPageID: facebookPageID,
				}
				user = &model.User{
					Session: model.UserSession{
						FacebookPageTokens: []model.FacebookPageToken{model.FacebookPageToken{
							PageID: facebookPageID,
						}},
					},
				}

				mockSessionManager.On("Get", mock.Anything).Return("session", nil)
				mockUsersCollection.On("GetSessionID", "session").Return(user, nil).Once()
				mockRestaurantsCollection.On("GetID", restaurantID).Return(restaurant, nil).Once()

				requestMethod = "POST"
				params = httprouter.Params{httprouter.Param{
					Key:   "restaurantID",
					Value: restaurantID.Hex(),
				}}
			})

			AfterEach(func() {
				mockSessionManager.AssertExpectations(GinkgoT())
				mockPostsCollection.AssertExpectations(GinkgoT())
				mockRestaurantsCollection.AssertExpectations(GinkgoT())
				mockUsersCollection.AssertExpectations(GinkgoT())
				facebookPost.AssertExpectations(GinkgoT())
			})

			Context("with valid input", func() {
				BeforeEach(func() {
					requestData = map[string]interface{}{
						"date": "2115-04-18",
					}
				})

				Context("with DB update succeeding", func() {
					var date model.DateWithoutTime
					BeforeEach(func() {
						date = model.DateWithoutTime("2115-04-18")
						mockPostsCollection.On("Insert", mock.AnythingOfType("[]*model.OfferGroupPost")).Return([]*model.OfferGroupPost{
							&model.OfferGroupPost{
								ID:              id,
								Date:            date,
								MessageTemplate: "messagetemplate",
							},
						}, nil)
					})

					Context("with FB update succeeding", func() {
						BeforeEach(func() {
							facebookPost.On("Update", date, user, restaurant).Return(nil)
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

						It("should include the post with the new ID", func() {
							handler(responseRecorder, request, params)
							var post model.OfferGroupPost
							json.Unmarshal(responseRecorder.Body.Bytes(), &post)
							Expect(post.ID).To(Equal(id))
						})
					})

					Context("with FB update failing", func() {
						var err *router.HandlerError
						BeforeEach(func() {
							err = &router.HandlerError{Message: "a message"}
							facebookPost.On("Update", date, user, restaurant).Return(err)
						})

						It("should fail", func() {
							handlerErr := handler(responseRecorder, request, params)
							Expect(handlerErr).To(Equal(err))
						})
					})

					Context("with a restaurant with an associated FB page", func() {
						BeforeEach(func() {
							fbPageID := "fbpageid"
							restaurant := &model.Restaurant{
								ID:             restaurantID,
								FacebookPageID: fbPageID,
							}
							user := &model.User{
								RestaurantIDs: []bson.ObjectId{restaurant.ID},
							}
							mockRestaurantsCollection.GetID(restaurantID) // Best way I could think of getting rid of the previous mock
							mockRestaurantsCollection.On("GetID", restaurantID).Return(restaurant, nil)
							mockUsersCollection.GetSessionID("session")
							mockUsersCollection.On("GetSessionID", "session").Return(user, nil)
							facebookPost.On("Update", model.DateWithoutTime("2115-04-18"), user, restaurant).Return(nil)
						})

						It("should succeed", func() {
							err := handler(responseRecorder, request, params)
							Expect(err).To(BeNil())
						})
					})
				})

				Context("with DB update failing", func() {
					BeforeEach(func() {
						mockPostsCollection.On("Insert", mock.AnythingOfType("[]*model.OfferGroupPost")).Return(nil, errors.New("things"))
					})

					It("should fail", func() {
						err := handler(responseRecorder, request, params)
						Expect(err).NotTo(BeNil())
					})
				})
			})

			Context("with an invalid date", func() {
				BeforeEach(func() {
					requestData = map[string]interface{}{
						"date": "2115-74-18",
					}
				})

				It("should fail", func() {
					err := handler(responseRecorder, request, params)
					Expect(err).NotTo(BeNil())
				})
			})
		})
	})

	Describe("PUT /restaurants/:restaurantID/posts/:date", func() {
		var (
			sessionManager        session.Manager
			postsCollection       db.OfferGroupPosts
			restaurantsCollection db.Restaurants
			usersCollection       db.Users
			facebookPost          *mocks.Post
			handler               router.HandlerWithParams
		)

		JustBeforeEach(func() {
			handler = PutOfferGroupPost(postsCollection, sessionManager, usersCollection, restaurantsCollection,
				facebookPost)
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
				restaurantID              bson.ObjectId
				id                        bson.ObjectId
				params                    httprouter.Params
				user                      *model.User
				restaurant                *model.Restaurant
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
				facebookPost = new(mocks.Post)

				id = bson.NewObjectId()

				restaurantID = bson.NewObjectId()
				facebookPageID := "fbpageid"
				restaurant = &model.Restaurant{
					ID:             restaurantID,
					FacebookPageID: facebookPageID,
				}
				user = &model.User{
					Session: model.UserSession{
						FacebookPageTokens: []model.FacebookPageToken{model.FacebookPageToken{
							PageID: facebookPageID,
						}},
					},
				}

				mockSessionManager.On("Get", mock.Anything).Return("session", nil)
				mockUsersCollection.On("GetSessionID", "session").Return(user, nil)
				mockRestaurantsCollection.On("GetID", restaurantID).Return(restaurant, nil)

				requestMethod = "PUT"
				params = httprouter.Params{httprouter.Param{
					Key:   "date",
					Value: "2015-04-10",
				}, httprouter.Param{
					Key:   "restaurantID",
					Value: restaurantID.Hex(),
				}}
				mockPostsCollection.On("GetByDate", model.DateWithoutTime("2015-04-10"), restaurantID).Return(&model.OfferGroupPost{
					ID:              id,
					Date:            model.DateWithoutTime("2015-04-10"),
					RestaurantID:    restaurantID,
					MessageTemplate: "an old message template",
				}, nil)
			})

			AfterEach(func() {
				mockSessionManager.AssertExpectations(GinkgoT())
				mockPostsCollection.AssertExpectations(GinkgoT())
				mockRestaurantsCollection.AssertExpectations(GinkgoT())
				mockUsersCollection.AssertExpectations(GinkgoT())
				facebookPost.AssertExpectations(GinkgoT())
			})

			Context("with valid input", func() {
				BeforeEach(func() {
					requestData = map[string]interface{}{
						// id and date are garbage because they should be ignored
						"_id":              bson.ObjectId("123456789012"),
						"date":             "garbage",
						"message_template": "a message template",
					}
				})

				Context("with DB update succeeding", func() {
					var date model.DateWithoutTime
					BeforeEach(func() {
						date = model.DateWithoutTime("2015-04-10")

						mockPostsCollection.On("UpdateByID", id, &model.OfferGroupPost{
							ID:              id,
							Date:            date,
							RestaurantID:    restaurantID,
							MessageTemplate: "a message template",
						}).Return(nil)
					})

					Context("with FB update succeeding", func() {
						BeforeEach(func() {
							facebookPost.On("Update", date, user, restaurant).Return(nil)
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

						It("should return the post", func() {
							handler(responseRecorder, request, params)
							var post *model.OfferGroupPost
							json.Unmarshal(responseRecorder.Body.Bytes(), &post)
							Expect(post.ID).To(Equal(id))
						})
					})

					Context("with FB update failing", func() {
						var err *router.HandlerError
						BeforeEach(func() {
							err = &router.HandlerError{Message: "a message"}
							facebookPost.On("Update", date, user, restaurant).Return(err)
						})

						It("should fail", func() {
							handlerErr := handler(responseRecorder, request, params)
							Expect(handlerErr).To(Equal(err))
						})
					})
				})

				Context("with DB update failing", func() {
					BeforeEach(func() {
						mockPostsCollection.On("UpdateByID", id, mock.AnythingOfType("*model.OfferGroupPost")).Return(errors.New("things"))
					})

					It("should fail", func() {
						err := handler(responseRecorder, request, params)
						Expect(err).NotTo(BeNil())
					})
				})
			})
		})
	})
})
