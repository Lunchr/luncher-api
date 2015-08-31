package facebook_test

import (
	"errors"
	"net/http"
	"time"

	"github.com/Lunchr/luncher-api/db/model"
	"github.com/Lunchr/luncher-api/facebook"
	"github.com/Lunchr/luncher-api/facebook/mocks"
	fbmodel "github.com/deiwin/facebook/model"
	"golang.org/x/oauth2"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Post", func() {
	var (
		facebookPost facebook.Post

		groupPosts       *mocks.OfferGroupPosts
		offersCollection *mocks.Offers
		regions          *mocks.Regions
		fbAuth           *mocks.Authenticator

		user       *model.User
		restaurant *model.Restaurant
	)

	BeforeEach(func() {
		groupPosts = new(mocks.OfferGroupPosts)
		offersCollection = new(mocks.Offers)
		regions = new(mocks.Regions)
		fbAuth = new(mocks.Authenticator)

		facebookPost = facebook.NewPost(groupPosts, offersCollection, regions, fbAuth)
	})

	JustBeforeEach(func() {
	})

	Describe("Update", func() {
		var date model.DateWithoutTime

		Context("for restaurants without an associated FB page", func() {
			BeforeEach(func() {
				restaurant = &model.Restaurant{
					FacebookPageID: "",
				}
				date = model.DateWithoutTime("2011-04-24")
			})

			It("does nothing", func() {
				err := facebookPost.Update(date, user, restaurant)
				Expect(err).NotTo(HaveOccurred())
			})
		})

		Context("for a restaurant with an associated FB page", func() {
			var (
				restaurantID      bson.ObjectId
				facebookPageID    string
				facebookUserToken *oauth2.Token
				facebookPageToken string
				fbAPI             *mocks.API
			)

			BeforeEach(func() {
				restaurantID = bson.NewObjectId()
				facebookPageID = "a page ID"
				facebookUserToken = &oauth2.Token{
					AccessToken: "a user token",
				}
				facebookPageToken = "a page token"

				regionName := "a region"
				region := &model.Region{
					Name:     regionName,
					Location: "UTC",
				}
				regions.On("GetName", regionName).Return(region, nil)

				restaurant = &model.Restaurant{
					ID:             restaurantID,
					FacebookPageID: facebookPageID,
					Region:         regionName,
				}
				user = &model.User{
					Session: &model.UserSession{
						FacebookUserToken: *facebookUserToken,
						FacebookPageToken: facebookPageToken,
					},
				}
				fbAPI = new(mocks.API)
			})

			Context("with an existing OfferGroupPost", func() {
				var (
					messageTemplate string
					offerGroupPost  *model.OfferGroupPost
					facebookPostID  string
					id              bson.ObjectId
					startTime       time.Time
					endTime         time.Time
				)

				BeforeEach(func() {
					id = bson.NewObjectId()
					messageTemplate = "a message template"
					offerGroupPost = &model.OfferGroupPost{
						ID:              id,
						Date:            date,
						MessageTemplate: messageTemplate,
					}
					groupPosts.On("GetByDate", date, restaurantID).Return(offerGroupPost, nil)

					facebookPostID = "fb post id"
					startTime = time.Date(2011, 04, 24, 0, 0, 0, 0, time.UTC)
					endTime = time.Date(2011, 04, 25, 0, 0, 0, 0, time.UTC)
				})

				AfterEach(func() {
					groupPosts.AssertExpectations(GinkgoT())
					offersCollection.AssertExpectations(GinkgoT())
					regions.AssertExpectations(GinkgoT())
					fbAuth.AssertExpectations(GinkgoT())
				})

				Context("without a previous associated FB post", func() {
					Context("with there being offers for that date", func() {
						BeforeEach(func() {
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
							offersCollection.On("GetForRestaurantWithinTimeBounds", restaurantID, startTime, endTime).Return(offers, nil)
							fbAuth.On("APIConnection", facebookUserToken).Return(fbAPI)
						})

						It("posts to FB and updates FB post ID in DB", func() {
							fbAPI.On("PagePublish", facebookPageToken, facebookPageID, &fbmodel.Post{
								Message: messageTemplate + "\n\natitle - 5.67€\nbtitle - 4.67€",
							}).Return(&fbmodel.Post{
								ID: facebookPostID,
							}, nil)
							groupPosts.On("UpdateByID", id, &model.OfferGroupPost{
								ID:              id,
								Date:            date,
								MessageTemplate: messageTemplate,
								FBPostID:        facebookPostID,
							}).Return(nil)

							err := facebookPost.Update(date, user, restaurant)
							Expect(err).To(BeNil())
						})
					})

					Context("without there being offers for that date", func() {
						BeforeEach(func() {
							offersCollection.On("GetForRestaurantWithinTimeBounds", restaurantID, startTime, endTime).Return([]*model.Offer{}, nil)
						})

						It("succeeds", func() {
							err := facebookPost.Update(date, user, restaurant)
							Expect(err).To(BeNil())
						})
					})
				})

				Context("with a previous associated FB post", func() {
					BeforeEach(func() {
						offerGroupPost.FBPostID = facebookPostID
					})

					Context("with there being offers for that date", func() {
						BeforeEach(func() {
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
							offersCollection.On("GetForRestaurantWithinTimeBounds", restaurantID, startTime, endTime).Return(offers, nil)
							fbAuth.On("APIConnection", facebookUserToken).Return(fbAPI)
						})

						It("updates the FB post", func() {
							fbAPI.On("PostUpdate", facebookPageToken, facebookPostID, &fbmodel.Post{
								Message: messageTemplate + "\n\natitle - 5.67€\nbtitle - 4.67€",
							}).Return(nil)

							err := facebookPost.Update(date, user, restaurant)
							Expect(err).To(BeNil())
						})
					})

					Context("without there being offers for that date", func() {
						BeforeEach(func() {
							offersCollection.On("GetForRestaurantWithinTimeBounds", restaurantID, startTime, endTime).Return([]*model.Offer{}, nil)
							fbAuth.On("APIConnection", facebookUserToken).Return(fbAPI)
						})

						Context("with post deletion failing", func() {
							var err error

							BeforeEach(func() {
								err = errors.New("something went wrong")
								fbAPI.On("PostDelete", facebookPageToken, facebookPostID).Return(err)
							})

							It("fails with the given error", func() {
								handlerErr := facebookPost.Update(date, user, restaurant)
								Expect(handlerErr).NotTo(BeNil())
								Expect(handlerErr.Err).To(Equal(err))
								Expect(handlerErr.Code).To(Equal(http.StatusBadGateway))
							})
						})

						Context("with post deletion succeeding", func() {
							BeforeEach(func() {
								fbAPI.On("PostDelete", facebookPageToken, facebookPostID).Return(nil)
							})

							It("removes the post ID from memory", func() {
								groupPosts.On("UpdateByID", id, &model.OfferGroupPost{
									ID:              id,
									Date:            date,
									MessageTemplate: messageTemplate,
									FBPostID:        "",
								}).Return(nil)

								err := facebookPost.Update(date, user, restaurant)
								Expect(err).To(BeNil())
							})
						})
					})
				})
			})

			Context("without an existing offer group post", func() {
				var defaultTemplate string

				BeforeEach(func() {
					defaultTemplate = "a default template message"
					restaurant.DefaultGroupPostMessageTemplate = defaultTemplate
					groupPosts.On("GetByDate", date, restaurantID).Return(nil, mgo.ErrNotFound)
				})

				Context("with insert failing", func() {
					var err error
					BeforeEach(func() {
						err = errors.New("something went wrong")
						groupPosts.On("Insert", []*model.OfferGroupPost{&model.OfferGroupPost{
							RestaurantID:    restaurantID,
							Date:            date,
							MessageTemplate: defaultTemplate,
						}}).Return(nil, err)
					})

					AfterEach(func() {
						groupPosts.AssertExpectations(GinkgoT())
					})

					It("fails with the given error", func() {
						handlerErr := facebookPost.Update(date, user, restaurant)
						Expect(handlerErr).NotTo(BeNil())
						Expect(handlerErr.Err).To(Equal(err))
						Expect(handlerErr.Code).To(Equal(http.StatusInternalServerError))
					})
				})
			})
		})
	})
})
