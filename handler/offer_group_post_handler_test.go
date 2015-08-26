package handler_test

import (
	"encoding/json"
	"errors"
	"time"

	"golang.org/x/oauth2"

	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/db/model"
	. "github.com/Lunchr/luncher-api/handler"
	"github.com/Lunchr/luncher-api/handler/mocks"
	"github.com/Lunchr/luncher-api/router"
	"github.com/Lunchr/luncher-api/session"
	"github.com/deiwin/facebook"
	fbmodel "github.com/deiwin/facebook/model"
	"github.com/julienschmidt/httprouter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

var _ = Describe("OfferGroupPostHandlers", func() {
	Describe("GET /restaurant/posts/:date", func() {
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

	Describe("POST /restaurant/posts", func() {
		var (
			sessionManager        session.Manager
			postsCollection       db.OfferGroupPosts
			restaurantsCollection db.Restaurants
			offersCollection      db.Offers
			regionsCollection     db.Regions
			fbAuth                facebook.Authenticator
			usersCollection       db.Users
			handler               router.Handler
		)

		JustBeforeEach(func() {
			handler = PostOfferGroupPost(postsCollection, sessionManager, usersCollection, restaurantsCollection,
				offersCollection, regionsCollection, fbAuth)
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
				mockPostsCollection       *mocks.OfferGroupPosts
				mockRestaurantsCollection *mocks.Restaurants
				mockUsersCollection       *mocks.Users
				mockOffersCollection      *mocks.Offers
				mockRegionsCollection     *mocks.Regions
				mockFBAuth                *mocks.Authenticator
				restaurantID              bson.ObjectId
				id                        bson.ObjectId
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
				mockOffersCollection = new(mocks.Offers)
				offersCollection = mockOffersCollection
				mockRegionsCollection = new(mocks.Regions)
				regionsCollection = mockRegionsCollection
				mockFBAuth = new(mocks.Authenticator)
				fbAuth = mockFBAuth

				id = bson.NewObjectId()

				restaurantID = bson.NewObjectId()
				restaurant := &model.Restaurant{
					ID: restaurantID,
				}
				user := &model.User{
					RestaurantIDs: []bson.ObjectId{restaurant.ID},
				}

				mockSessionManager.On("Get", mock.Anything).Return("session", nil)
				mockUsersCollection.On("GetSessionID", "session").Return(user, nil).Once()
				mockRestaurantsCollection.On("GetID", restaurantID).Return(restaurant, nil).Once()

				requestMethod = "POST"
			})

			AfterEach(func() {
				mockSessionManager.AssertExpectations(GinkgoT())
				mockPostsCollection.AssertExpectations(GinkgoT())
				mockRestaurantsCollection.AssertExpectations(GinkgoT())
				mockUsersCollection.AssertExpectations(GinkgoT())
			})

			Context("with valid input", func() {
				BeforeEach(func() {
					requestData = map[string]interface{}{
						"date": "2115-04-18",
					}
				})

				Context("with DB update succeeding", func() {
					BeforeEach(func() {
						mockPostsCollection.On("Insert", mock.AnythingOfType("[]*model.OfferGroupPost")).Return([]*model.OfferGroupPost{
							&model.OfferGroupPost{
								ID:              id,
								Date:            model.DateWithoutTime("2115-04-18"),
								MessageTemplate: "messagetemplate",
							},
						}, nil)
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

					It("should include the post with the new ID", func() {
						handler(responseRecorder, request)
						var post model.OfferGroupPost
						json.Unmarshal(responseRecorder.Body.Bytes(), &post)
						Expect(post.ID).To(Equal(id))
					})

					Context("with a restaurant with an associated FB page", func() {
						var fbPostID string

						BeforeEach(func() {
							fbPageID := "fbpageid"
							regionName := "regionname"
							restaurant := &model.Restaurant{
								ID:             restaurantID,
								FacebookPageID: fbPageID,
								Region:         regionName,
							}
							fbUserToken := &oauth2.Token{}
							fbPageToken := "afbpagetoken"
							user := &model.User{
								RestaurantIDs: []bson.ObjectId{restaurant.ID},
								Session: &model.UserSession{
									FacebookUserToken: *fbUserToken,
									FacebookPageToken: fbPageToken,
								},
							}
							mockRestaurantsCollection.GetID(restaurantID) // Best way I could think of getting rid of the previous mock
							mockRestaurantsCollection.On("GetID", restaurantID).Return(restaurant, nil)
							mockUsersCollection.GetSessionID("session")
							mockUsersCollection.On("GetSessionID", "session").Return(user, nil)
							fbAPI := new(mocks.API)
							mockFBAuth.On("APIConnection", fbUserToken).Return(fbAPI)
							region := &model.Region{
								Location: "UTC",
							}
							mockRegionsCollection.On("GetName", regionName).Return(region, nil)
							startTime := time.Date(2115, 04, 18, 0, 0, 0, 0, time.UTC)
							endTime := time.Date(2115, 04, 19, 0, 0, 0, 0, time.UTC)
							offers := []*model.Offer{
								&model.Offer{
									CommonOfferFields: model.CommonOfferFields{
										Title: "atitle",
										Price: 5.670000000000,
									},
								},
								&model.Offer{
									CommonOfferFields: model.CommonOfferFields{
										Title: "btitle",
										Price: 4.670000000000,
									},
								},
							}
							mockOffersCollection.On("GetForRestaurantWithinTimeBounds", restaurantID, startTime, endTime).Return(offers, nil)
							fbPostID = "fbpostid"
							fbAPI.On("PagePublish", fbPageToken, fbPageID, "messagetemplate\n\natitle - 5.67€\nbtitle - 4.67€").Return(&fbmodel.Post{
								ID: fbPostID,
							}, nil)
							mockPostsCollection.On("UpdateByID", id, &model.OfferGroupPost{
								ID:              id,
								Date:            model.DateWithoutTime("2115-04-18"),
								MessageTemplate: "messagetemplate",
								FBPostID:        fbPostID,
							}).Return(nil)
						})

						It("should succeed", func() {
							err := handler(responseRecorder, request)
							Expect(err).To(BeNil())
						})
					})
				})

				Context("with DB update failing", func() {
					BeforeEach(func() {
						mockPostsCollection.On("Insert", mock.AnythingOfType("[]*model.OfferGroupPost")).Return(nil, errors.New("things"))
					})

					It("should fail", func() {
						err := handler(responseRecorder, request)
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
					err := handler(responseRecorder, request)
					Expect(err).NotTo(BeNil())
				})
			})
		})
	})

	Describe("PUT /restaurant/posts/:date", func() {
		var (
			sessionManager        session.Manager
			postsCollection       db.OfferGroupPosts
			restaurantsCollection db.Restaurants
			usersCollection       db.Users
			offersCollection      db.Offers
			regionsCollection     db.Regions
			fbAuth                facebook.Authenticator
			handler               router.HandlerWithParams
		)

		JustBeforeEach(func() {
			handler = PutOfferGroupPost(postsCollection, sessionManager, usersCollection, restaurantsCollection,
				offersCollection, regionsCollection, fbAuth)
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
				id = bson.NewObjectId()

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

				requestMethod = "PUT"
				params = httprouter.Params{httprouter.Param{
					Key:   "date",
					Value: "2015-04-10",
				}}
			})

			AfterEach(func() {
				mockSessionManager.AssertExpectations(GinkgoT())
				mockPostsCollection.AssertExpectations(GinkgoT())
				mockRestaurantsCollection.AssertExpectations(GinkgoT())
				mockUsersCollection.AssertExpectations(GinkgoT())
			})

			Context("with valid input", func() {
				BeforeEach(func() {
					requestData = map[string]interface{}{
						"_id":              id,
						"date":             "2015-04-10",
						"message_template": "a message template",
					}
				})

				Context("with DB update succeeding", func() {
					BeforeEach(func() {
						mockPostsCollection.On("UpdateByID", id, &model.OfferGroupPost{
							ID:              id,
							Date:            "2015-04-10",
							RestaurantID:    restaurantID,
							MessageTemplate: "a message template",
						}).Return(nil)
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

			Context("with a non-matching date", func() {
				BeforeEach(func() {
					requestData = map[string]interface{}{
						"_id":  id,
						"date": "2015-04-11",
					}
				})

				It("should fail", func() {
					err := handler(responseRecorder, request, params)
					Expect(err).NotTo(BeNil())
				})
			})

			Context("without an id", func() {
				BeforeEach(func() {
					requestData = map[string]interface{}{
						"date": "2015-04-10",
					}
					mockPostsCollection.On("UpdateByID", bson.ObjectId(""), mock.AnythingOfType("*model.OfferGroupPost")).Return(errors.New("missing stuff"))
				})

				It("should fail", func() {
					err := handler(responseRecorder, request, params)
					Expect(err).NotTo(BeNil())
				})
			})
		})
	})
})
