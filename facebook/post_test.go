package facebook_test

import (
	"errors"
	"net/http"
	"time"

	"github.com/Lunchr/luncher-api/db/model"
	"github.com/Lunchr/luncher-api/facebook"
	"github.com/Lunchr/luncher-api/facebook/mocks"
	fbmodel "github.com/deiwin/facebook/model"
	"github.com/stretchr/testify/mock"
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

		BeforeEach(func() {
			date = model.DateWithoutTime("2011-04-24")
		})

		Context("for restaurants without an associated FB page", func() {
			BeforeEach(func() {
				restaurant = &model.Restaurant{
					FacebookPageID: "",
				}
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
							fbAuth.On("APIConnection", facebookUserToken).Return(fbAPI)
							groupPosts.On("UpdateByID", id, &model.OfferGroupPost{
								ID:              id,
								Date:            date,
								MessageTemplate: messageTemplate,
								FBPostID:        facebookPostID,
							}).Return(nil)
						})

						Context("for far future offers", func() {
							BeforeEach(func() {
								offersCollection.On("GetForRestaurantWithinTimeBounds", restaurantID, startTime, endTime).Return([]*model.Offer{
									&model.Offer{
										CommonOfferFields: model.CommonOfferFields{
											Title:    "atitle",
											Price:    5.670000000000,
											FromTime: time.Date(2115, 01, 02, 10, 0, 0, 0, time.UTC),
										},
									},
									&model.Offer{
										CommonOfferFields: model.CommonOfferFields{
											Title:    "btitle",
											Price:    4.670000000000,
											FromTime: time.Date(2115, 01, 02, 9, 0, 0, 0, time.UTC),
										},
									},
								}, nil)
							})

							It("should set the post to be published right before the earliest offer", func() {
								fbAPI.On("PagePublish", facebookPageToken, facebookPageID, &fbmodel.Post{
									Message:              messageTemplate + "\n\natitle - 5.67€\nbtitle - 4.67€",
									Published:            false,
									ScheduledPublishTime: time.Date(2115, 01, 02, 8, 30, 0, 0, time.UTC),
								}).Return(&fbmodel.PostResponse{
									Post: fbmodel.Post{
										ID: facebookPostID,
									},
								}, nil)

								err := facebookPost.Update(date, user, restaurant)
								Expect(err).To(BeNil())
							})
						})

						Context("for near future offers", func() {
							BeforeEach(func() {
								offersCollection.On("GetForRestaurantWithinTimeBounds", restaurantID, startTime, endTime).Return([]*model.Offer{
									&model.Offer{
										CommonOfferFields: model.CommonOfferFields{
											Title:    "atitle",
											Price:    5.670000000000,
											FromTime: time.Now().Add(time.Minute),
										},
									},
									&model.Offer{
										CommonOfferFields: model.CommonOfferFields{
											Title:    "btitle",
											Price:    4.670000000000,
											FromTime: time.Now().Add(time.Hour),
										},
									},
								}, nil)
							})

							It("should leave some time to still modify the offer before going live", func() {
								fbAPI.On("PagePublish", facebookPageToken, facebookPageID, mock.AnythingOfType("*model.Post")).Return(&fbmodel.PostResponse{
									Post: fbmodel.Post{
										ID: facebookPostID,
									},
								}, nil)

								err := facebookPost.Update(date, user, restaurant)
								Expect(err).To(BeNil())
								post := fbAPI.Calls[0].Arguments.Get(2).(*fbmodel.Post)
								Expect(post.Message).To(Equal(messageTemplate + "\n\natitle - 5.67€\nbtitle - 4.67€"))
								Expect(post.Published).To(BeFalse())
								Expect(post.ScheduledPublishTime.Sub(time.Now())).To(BeNumerically("~", 5*time.Minute, time.Second))
							})
						})

						Context("for past offers", func() {
							BeforeEach(func() {
								offersCollection.On("GetForRestaurantWithinTimeBounds", restaurantID, startTime, endTime).Return([]*model.Offer{
									&model.Offer{
										CommonOfferFields: model.CommonOfferFields{
											Title:    "atitle",
											Price:    5.670000000000,
											FromTime: time.Date(2005, 01, 02, 9, 0, 0, 0, time.UTC),
										},
									},
									&model.Offer{
										CommonOfferFields: model.CommonOfferFields{
											Title:    "btitle",
											Price:    4.670000000000,
											FromTime: time.Date(2005, 01, 02, 9, 0, 0, 0, time.UTC),
										},
									},
								}, nil)
							})

							It("should leave some time to still modify the offer before going live", func() {
								fbAPI.On("PagePublish", facebookPageToken, facebookPageID, mock.AnythingOfType("*model.Post")).Return(&fbmodel.PostResponse{
									Post: fbmodel.Post{
										ID: facebookPostID,
									},
								}, nil)

								err := facebookPost.Update(date, user, restaurant)
								Expect(err).To(BeNil())
								post := fbAPI.Calls[0].Arguments.Get(2).(*fbmodel.Post)
								Expect(post.Message).To(Equal(messageTemplate + "\n\natitle - 5.67€\nbtitle - 4.67€"))
								Expect(post.Published).To(BeFalse())
								Expect(post.ScheduledPublishTime.Sub(time.Now())).To(BeNumerically("~", 5*time.Minute, time.Second))
							})
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
							fbAuth.On("APIConnection", facebookUserToken).Return(fbAPI)
						})

						Context("with existing post being published", func() {
							BeforeEach(func() {
								fbAPI.On("Post", facebookPageToken, facebookPostID).Return(&fbmodel.PostResponse{
									IsPublished: true,
								}, nil)
								offersCollection.On("GetForRestaurantWithinTimeBounds", restaurantID, startTime, endTime).Return([]*model.Offer{
									&model.Offer{
										CommonOfferFields: model.CommonOfferFields{
											Title:    "atitle",
											Price:    5.670000000000,
											FromTime: time.Date(2115, 01, 02, 10, 0, 0, 0, time.UTC),
										},
									},
									&model.Offer{
										CommonOfferFields: model.CommonOfferFields{
											Title:    "btitle",
											Price:    4.670000000000,
											FromTime: time.Date(2115, 01, 02, 9, 0, 0, 0, time.UTC),
										},
									},
								}, nil)
							})

							It("should update the post as published", func() {
								fbAPI.On("PostUpdate", facebookPageToken, facebookPostID, &fbmodel.Post{
									Message:   messageTemplate + "\n\natitle - 5.67€\nbtitle - 4.67€",
									Published: true,
								}).Return(nil)

								err := facebookPost.Update(date, user, restaurant)
								Expect(err).To(BeNil())
							})
						})

						Context("with existing post not being published", func() {
							BeforeEach(func() {
								fbAPI.On("Post", facebookPageToken, facebookPostID).Return(&fbmodel.PostResponse{
									IsPublished: false,
								}, nil)
							})

							Context("for far future offers", func() {
								BeforeEach(func() {
									offersCollection.On("GetForRestaurantWithinTimeBounds", restaurantID, startTime, endTime).Return([]*model.Offer{
										&model.Offer{
											CommonOfferFields: model.CommonOfferFields{
												Title:    "atitle",
												Price:    5.670000000000,
												FromTime: time.Date(2115, 01, 02, 10, 0, 0, 0, time.UTC),
											},
										},
										&model.Offer{
											CommonOfferFields: model.CommonOfferFields{
												Title:    "btitle",
												Price:    4.670000000000,
												FromTime: time.Date(2115, 01, 02, 9, 0, 0, 0, time.UTC),
											},
										},
									}, nil)
								})

								It("should set the post to be published right before the earliest offer", func() {
									fbAPI.On("PostUpdate", facebookPageToken, facebookPostID, &fbmodel.Post{
										Message:              messageTemplate + "\n\natitle - 5.67€\nbtitle - 4.67€",
										Published:            false,
										ScheduledPublishTime: time.Date(2115, 01, 02, 8, 30, 0, 0, time.UTC),
									}).Return(nil)

									err := facebookPost.Update(date, user, restaurant)
									Expect(err).To(BeNil())
								})
							})

							Context("for near future offers", func() {
								BeforeEach(func() {
									offersCollection.On("GetForRestaurantWithinTimeBounds", restaurantID, startTime, endTime).Return([]*model.Offer{
										&model.Offer{
											CommonOfferFields: model.CommonOfferFields{
												Title:    "atitle",
												Price:    5.670000000000,
												FromTime: time.Now().Add(time.Minute),
											},
										},
										&model.Offer{
											CommonOfferFields: model.CommonOfferFields{
												Title:    "btitle",
												Price:    4.670000000000,
												FromTime: time.Now().Add(time.Hour),
											},
										},
									}, nil)
								})

								It("should leave some time to still modify the offer before going live", func() {
									fbAPI.On("PostUpdate", facebookPageToken, facebookPostID, mock.AnythingOfType("*model.Post")).Return(nil)

									err := facebookPost.Update(date, user, restaurant)
									Expect(err).To(BeNil())
									post := fbAPI.Calls[1].Arguments.Get(2).(*fbmodel.Post)
									Expect(post.Message).To(Equal(messageTemplate + "\n\natitle - 5.67€\nbtitle - 4.67€"))
									Expect(post.Published).To(BeFalse())
									Expect(post.ScheduledPublishTime.Sub(time.Now())).To(BeNumerically("~", 5*time.Minute, time.Second))
								})
							})

							Context("for past offers", func() {
								BeforeEach(func() {
									offersCollection.On("GetForRestaurantWithinTimeBounds", restaurantID, startTime, endTime).Return([]*model.Offer{
										&model.Offer{
											CommonOfferFields: model.CommonOfferFields{
												Title:    "atitle",
												Price:    5.670000000000,
												FromTime: time.Date(2005, 01, 02, 9, 0, 0, 0, time.UTC),
											},
										},
										&model.Offer{
											CommonOfferFields: model.CommonOfferFields{
												Title:    "btitle",
												Price:    4.670000000000,
												FromTime: time.Date(2005, 01, 02, 9, 0, 0, 0, time.UTC),
											},
										},
									}, nil)
								})

								It("should leave some time to still modify the offer before going live", func() {
									fbAPI.On("PostUpdate", facebookPageToken, facebookPostID, mock.AnythingOfType("*model.Post")).Return(nil)

									err := facebookPost.Update(date, user, restaurant)
									Expect(err).To(BeNil())
									post := fbAPI.Calls[1].Arguments.Get(2).(*fbmodel.Post)
									Expect(post.Message).To(Equal(messageTemplate + "\n\natitle - 5.67€\nbtitle - 4.67€"))
									Expect(post.Published).To(BeFalse())
									Expect(post.ScheduledPublishTime.Sub(time.Now())).To(BeNumerically("~", 5*time.Minute, time.Second))
								})
							})
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
