package facebook_test

import (
	"time"

	"github.com/Lunchr/luncher-api/db/model"
	"github.com/Lunchr/luncher-api/facebook"
	"github.com/Lunchr/luncher-api/facebook/mocks"
	fbmodel "github.com/deiwin/facebook/model"
	"golang.org/x/oauth2"
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

				startTime := time.Date(2011, 04, 24, 0, 0, 0, 0, time.UTC)
				endTime := time.Date(2011, 04, 25, 0, 0, 0, 0, time.UTC)
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
				fbAuth.On("APIConnection", facebookUserToken).Return(fbAPI)
			})

			Context("with an existing OfferGroupPost", func() {
				var (
					messageTemplate string
					offerGroupPost  *model.OfferGroupPost
					facebookPostID  string
				)

				BeforeEach(func() {
					id := bson.NewObjectId()
					messageTemplate = "a message template"
					offerGroupPost = &model.OfferGroupPost{
						ID:              id,
						Date:            date,
						MessageTemplate: messageTemplate,
					}
					groupPosts.On("GetByDate", date, restaurantID).Return(offerGroupPost, nil)

					facebookPostID = "fb post id"
					fbAPI.On("PagePublish", facebookPageToken, facebookPageID,
						messageTemplate+"\n\natitle - 5.67€\nbtitle - 4.67€").Return(&fbmodel.Post{
						ID: facebookPostID,
					}, nil)
					groupPosts.On("UpdateByID", id, &model.OfferGroupPost{
						ID:              id,
						Date:            date,
						MessageTemplate: messageTemplate,
						FBPostID:        facebookPostID,
					}).Return(nil)
				})

				Context("without a previous associated FB post", func() {
					AfterEach(func() {
						groupPosts.AssertExpectations(GinkgoT())
						offersCollection.AssertExpectations(GinkgoT())
						regions.AssertExpectations(GinkgoT())
						fbAuth.AssertExpectations(GinkgoT())
					})

					It("should succeed", func() {
						err := facebookPost.Update(date, user, restaurant)
						Expect(err).To(BeNil())
					})
				})
			})
		})
	})
})
