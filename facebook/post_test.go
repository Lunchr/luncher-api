package facebook_test

import (
	"bytes"
	"errors"
	"hash/crc32"
	"image"
	"image/color"
	"image/jpeg"
	"io"
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
		images           *mocks.Images
		collageLayout    *mocks.Layout

		user       *model.User
		restaurant *model.Restaurant
	)

	BeforeEach(func() {
		groupPosts = new(mocks.OfferGroupPosts)
		offersCollection = new(mocks.Offers)
		regions = new(mocks.Regions)
		fbAuth = new(mocks.Authenticator)
		images = new(mocks.Images)
		collageLayout = new(mocks.Layout)

		facebookPost = facebook.NewPost(groupPosts, offersCollection, regions, fbAuth, images, collageLayout)
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
					Session: model.UserSession{
						FacebookUserToken: *facebookUserToken,
						FacebookPageTokens: []model.FacebookPageToken{model.FacebookPageToken{
							PageID: facebookPageID,
							Token:  facebookPageToken,
						}},
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
					photoID         = "photo ID"

					bluePixel         image.Image
					bluePixelJPEGData io.Reader
					bluePixelChecksum uint32

					redPixel         image.Image
					redPixelJPEGData io.Reader
					redPixelChecksum uint32

					greenPixel         image.Image
					greenPixelJPEGData io.Reader
					greenPixelChecksum uint32
				)

				var getSinglePixelImage = func(c color.Color) image.Image {
					i := image.NewRGBA(image.Rect(0, 0, 1, 1))
					i.Set(0, 0, c)
					return i
				}

				var getJPEGImageData = func(i image.Image) (io.Reader, uint32) {
					var imageData bytes.Buffer
					err := jpeg.Encode(&imageData, i, nil)
					Expect(err).NotTo(HaveOccurred())
					crc := crc32.ChecksumIEEE(imageData.Bytes())
					return &imageData, crc
				}

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

					bluePixel = getSinglePixelImage(color.RGBA{0x00, 0x00, 0xff, 0xff})
					bluePixelJPEGData, bluePixelChecksum = getJPEGImageData(bluePixel)
					redPixel = getSinglePixelImage(color.RGBA{0xff, 0x00, 0x00, 0xff})
					redPixelJPEGData, redPixelChecksum = getJPEGImageData(redPixel)
					greenPixel = getSinglePixelImage(color.RGBA{0xff, 0x00, 0x00, 0xff})
					greenPixelJPEGData, greenPixelChecksum = getJPEGImageData(greenPixel)
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
									ScheduledPublishTime: time.Date(2115, 01, 02, 8, 45, 0, 0, time.UTC),
								}).Return(&fbmodel.PostResponse{
									ID: facebookPostID,
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
									ID: facebookPostID,
								}, nil)

								err := facebookPost.Update(date, user, restaurant)
								Expect(err).To(BeNil())
								post := fbAPI.Calls[0].Arguments.Get(2).(*fbmodel.Post)
								Expect(post.Message).To(Equal(messageTemplate + "\n\natitle - 5.67€\nbtitle - 4.67€"))
								Expect(post.Published).To(BeFalse())
								Expect(post.ScheduledPublishTime.Sub(time.Now())).To(BeNumerically("~", 11*time.Minute, time.Second))
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
									ID: facebookPostID,
								}, nil)

								err := facebookPost.Update(date, user, restaurant)
								Expect(err).To(BeNil())
								post := fbAPI.Calls[0].Arguments.Get(2).(*fbmodel.Post)
								Expect(post.Message).To(Equal(messageTemplate + "\n\natitle - 5.67€\nbtitle - 4.67€"))
								Expect(post.Published).To(BeFalse())
								Expect(post.ScheduledPublishTime.Sub(time.Now())).To(BeNumerically("~", 11*time.Minute, time.Second))
							})
						})
					})

					Context("with there being offers with images for that date", func() {
						BeforeEach(func() {
							fbAuth.On("APIConnection", facebookUserToken).Return(fbAPI)
						})

						Describe("with one of the offers having an image", func() {
							BeforeEach(func() {
								groupPosts.On("UpdateByID", id, &model.OfferGroupPost{
									ID:                  id,
									Date:                date,
									MessageTemplate:     messageTemplate,
									FBPostID:            facebookPageID + "_" + photoID,
									PostedImageChecksum: bluePixelChecksum,
								}).Return(nil)
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
											ImageChecksum: "checksum1",
										},
										&model.Offer{
											CommonOfferFields: model.CommonOfferFields{
												Title:    "btitle",
												Price:    4.670000000000,
												FromTime: time.Date(2005, 01, 02, 9, 0, 0, 0, time.UTC),
											},
										},
									}, nil)
									images.On("GetOriginal", "checksum1").Return(bluePixel, nil)
								})

								Context("with photo response including a post ID", func() {
									BeforeEach(func() {
										fbAPI.On("PagePhotoCreate", facebookPageToken, facebookPageID, mock.AnythingOfType("*model.Photo")).Return(&fbmodel.PhotoResponse{
											ID:     "whatever",
											PostID: facebookPageID + "_" + photoID,
										}, nil)
									})

									It("should post the photot with the single image", func() {
										err := facebookPost.Update(date, user, restaurant)
										Expect(err).To(BeNil())
										post := fbAPI.Calls[0].Arguments.Get(2).(*fbmodel.Photo)
										Expect(post.Message).To(Equal(messageTemplate + "\n\natitle - 5.67€\nbtitle - 4.67€"))
										Expect(post.Published).To(BeFalse())
										Expect(post.ScheduledPublishTime.Sub(time.Now())).To(BeNumerically("~", 11*time.Minute, time.Second))
										Expect(post.Photo).To(Equal(bluePixelJPEGData))
									})
								})

								Context("without photo response including a post ID", func() {
									BeforeEach(func() {
										fbAPI.On("PagePhotoCreate", facebookPageToken, facebookPageID, mock.AnythingOfType("*model.Photo")).Return(&fbmodel.PhotoResponse{
											ID: photoID,
										}, nil)
									})

									It("should leave some time to still modify the offer before going live", func() {
										err := facebookPost.Update(date, user, restaurant)
										Expect(err).To(BeNil())
										post := fbAPI.Calls[0].Arguments.Get(2).(*fbmodel.Photo)
										Expect(post.Message).To(Equal(messageTemplate + "\n\natitle - 5.67€\nbtitle - 4.67€"))
										Expect(post.Published).To(BeFalse())
										Expect(post.ScheduledPublishTime.Sub(time.Now())).To(BeNumerically("~", 11*time.Minute, time.Second))
										Expect(post.Photo).To(Equal(bluePixelJPEGData))
									})
								})
							})
						})

						Describe("with two of the offers having an image", func() {
							BeforeEach(func() {
								groupPosts.On("UpdateByID", id, &model.OfferGroupPost{
									ID:                  id,
									Date:                date,
									MessageTemplate:     messageTemplate,
									FBPostID:            facebookPageID + "_" + photoID,
									PostedImageChecksum: greenPixelChecksum,
								}).Return(nil)
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
											ImageChecksum: "checksum1",
										},
										&model.Offer{
											CommonOfferFields: model.CommonOfferFields{
												Title:    "btitle",
												Price:    4.670000000000,
												FromTime: time.Date(2005, 01, 02, 9, 0, 0, 0, time.UTC),
											},
											ImageChecksum: "checksum2",
										},
									}, nil)
									images.On("GetOriginal", "checksum1").Return(bluePixel, nil)
									images.On("GetOriginal", "checksum2").Return(redPixel, nil)
									picassoNode := new(mocks.Node)
									collageLayout.On("Compose", []image.Image{bluePixel, redPixel}).Return(picassoNode)
									white := color.RGBA{0xff, 0xff, 0xff, 0xff}
									picassoNode.On("DrawWithBorder", 800, 800, white, 2).Return(greenPixel)
								})

								Context("with photo response including a post ID", func() {
									BeforeEach(func() {
										fbAPI.On("PagePhotoCreate", facebookPageToken, facebookPageID, mock.AnythingOfType("*model.Photo")).Return(&fbmodel.PhotoResponse{
											ID:     "whatever",
											PostID: facebookPageID + "_" + photoID,
										}, nil)
									})

									It("should post the photot with the single image", func() {
										err := facebookPost.Update(date, user, restaurant)
										Expect(err).To(BeNil())
										post := fbAPI.Calls[0].Arguments.Get(2).(*fbmodel.Photo)
										Expect(post.Message).To(Equal(messageTemplate + "\n\natitle - 5.67€\nbtitle - 4.67€"))
										Expect(post.Published).To(BeFalse())
										Expect(post.ScheduledPublishTime.Sub(time.Now())).To(BeNumerically("~", 11*time.Minute, time.Second))
										Expect(post.Photo).To(Equal(greenPixelJPEGData))
									})
								})
							})
						})

						Describe("with five (> 4) of the offers having an image", func() {
							BeforeEach(func() {
								groupPosts.On("UpdateByID", id, &model.OfferGroupPost{
									ID:                  id,
									Date:                date,
									MessageTemplate:     messageTemplate,
									FBPostID:            facebookPageID + "_" + photoID,
									PostedImageChecksum: greenPixelChecksum,
								}).Return(nil)
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
											ImageChecksum: "checksum1",
										},
										&model.Offer{
											CommonOfferFields: model.CommonOfferFields{
												Title:    "btitle",
												Price:    4.670000000000,
												FromTime: time.Date(2005, 01, 02, 9, 0, 0, 0, time.UTC),
											},
											ImageChecksum: "checksum2",
										},
										&model.Offer{
											CommonOfferFields: model.CommonOfferFields{
												Title:    "btitle",
												Price:    4.670000000000,
												FromTime: time.Date(2005, 01, 02, 9, 0, 0, 0, time.UTC),
											},
											ImageChecksum: "checksum1",
										},
										&model.Offer{
											CommonOfferFields: model.CommonOfferFields{
												Title:    "btitle",
												Price:    4.670000000000,
												FromTime: time.Date(2005, 01, 02, 9, 0, 0, 0, time.UTC),
											},
											ImageChecksum: "checksum2",
										},
										&model.Offer{
											CommonOfferFields: model.CommonOfferFields{
												Title:    "btitle",
												Price:    4.670000000000,
												FromTime: time.Date(2005, 01, 02, 9, 0, 0, 0, time.UTC),
											},
											ImageChecksum: "checksum1",
										},
									}, nil)
									images.On("GetOriginal", "checksum1").Return(bluePixel, nil)
									images.On("GetOriginal", "checksum2").Return(redPixel, nil)
									picassoNode := new(mocks.Node)
									collageLayout.On("Compose", []image.Image{bluePixel, redPixel, bluePixel, redPixel}).Return(picassoNode)
									white := color.RGBA{0xff, 0xff, 0xff, 0xff}
									picassoNode.On("DrawWithBorder", 800, 800, white, 2).Return(greenPixel)
								})

								Context("with photo response including a post ID", func() {
									BeforeEach(func() {
										fbAPI.On("PagePhotoCreate", facebookPageToken, facebookPageID, mock.AnythingOfType("*model.Photo")).Return(&fbmodel.PhotoResponse{
											ID:     "whatever",
											PostID: facebookPageID + "_" + photoID,
										}, nil)
									})

									It("should post the photot with the single image", func() {
										err := facebookPost.Update(date, user, restaurant)
										Expect(err).To(BeNil())
										post := fbAPI.Calls[0].Arguments.Get(2).(*fbmodel.Photo)
										Expect(post.Message).To(Equal(messageTemplate + "\n\natitle - 5.67€\nbtitle - 4.67€\nbtitle - 4.67€\nbtitle - 4.67€\nbtitle - 4.67€"))
										Expect(post.Published).To(BeFalse())
										Expect(post.ScheduledPublishTime.Sub(time.Now())).To(BeNumerically("~", 11*time.Minute, time.Second))
										Expect(post.Photo).To(Equal(greenPixelJPEGData))
									})
								})
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
							var originalPublishTime = time.Date(2015, 01, 02, 10, 0, 0, 0, time.UTC)
							BeforeEach(func() {
								fbAPI.On("Post", facebookPageToken, facebookPostID).Return(&fbmodel.PostResponse{
									IsPublished: true,
									CreatedTime: originalPublishTime,
								}, nil)
							})

							Context("without the offers for the date having images", func() {
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

								It("should update the post as published", func() {
									fbAPI.On("PostUpdate", facebookPageToken, facebookPostID, &fbmodel.Post{
										Message:   messageTemplate + "\n\natitle - 5.67€\nbtitle - 4.67€",
										Published: true,
									}).Return(nil)

									err := facebookPost.Update(date, user, restaurant)
									Expect(err).To(BeNil())
								})

								Context("with previous post having a collage", func() {
									BeforeEach(func() {
										offerGroupPost.PostedImageChecksum = bluePixelChecksum
									})

									It("deletes the current post and creates a new backdated one without a collage", func() {
										fbAPI.On("PostDelete", facebookPageToken, facebookPostID).Return(nil)
										fbAPI.On("PagePublish", facebookPageToken, facebookPageID, mock.AnythingOfType("*model.Post")).Return(&fbmodel.PostResponse{
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
										post := fbAPI.Calls[2].Arguments.Get(2).(*fbmodel.Post)
										Expect(post.Message).To(Equal(messageTemplate + "\n\natitle - 5.67€\nbtitle - 4.67€"))
										Expect(post.Published).To(BeTrue())
										Expect(post.BackdatedTime).To(Equal(originalPublishTime))
									})
								})
							})

							Context("with the offers for the date having images", func() {
								BeforeEach(func() {
									offersCollection.On("GetForRestaurantWithinTimeBounds", restaurantID, startTime, endTime).Return([]*model.Offer{
										&model.Offer{
											CommonOfferFields: model.CommonOfferFields{
												Title:    "atitle",
												Price:    5.670000000000,
												FromTime: time.Date(2115, 01, 02, 10, 0, 0, 0, time.UTC),
											},
											ImageChecksum: "checksum1",
										},
										&model.Offer{
											CommonOfferFields: model.CommonOfferFields{
												Title:    "btitle",
												Price:    4.670000000000,
												FromTime: time.Date(2115, 01, 02, 9, 0, 0, 0, time.UTC),
											},
										},
									}, nil)
									images.On("GetOriginal", "checksum1").Return(bluePixel, nil)
								})

								Context("with previous post having a collage", func() {
									Context("with the same collage", func() {
										BeforeEach(func() {
											offerGroupPost.PostedImageChecksum = bluePixelChecksum
										})

										It("updates the current post's message", func() {
											fbAPI.On("PostUpdate", facebookPageToken, facebookPostID, &fbmodel.Post{
												Message:   messageTemplate + "\n\natitle - 5.67€\nbtitle - 4.67€",
												Published: true,
											}).Return(nil)

											err := facebookPost.Update(date, user, restaurant)
											Expect(err).To(BeNil())
										})
									})

									Context("with a different collage", func() {
										BeforeEach(func() {
											offerGroupPost.PostedImageChecksum = greenPixelChecksum
										})

										It("deletes the current post and creates a new backdated one", func() {
											fbAPI.On("PostDelete", facebookPageToken, facebookPostID).Return(nil)
											fbAPI.On("PagePhotoCreate", facebookPageToken, facebookPageID, mock.AnythingOfType("*model.Photo")).Return(&fbmodel.PhotoResponse{
												ID:     "whatever",
												PostID: facebookPageID + "_" + photoID,
											}, nil)
											groupPosts.On("UpdateByID", id, &model.OfferGroupPost{
												ID:                  id,
												Date:                date,
												MessageTemplate:     messageTemplate,
												FBPostID:            facebookPageID + "_" + photoID,
												PostedImageChecksum: bluePixelChecksum,
											}).Return(nil)

											err := facebookPost.Update(date, user, restaurant)
											Expect(err).To(BeNil())
											post := fbAPI.Calls[2].Arguments.Get(2).(*fbmodel.Photo)
											Expect(post.Message).To(Equal(messageTemplate + "\n\natitle - 5.67€\nbtitle - 4.67€"))
											Expect(post.Published).To(BeTrue())
											Expect(post.BackdatedTime).To(Equal(originalPublishTime))
											Expect(post.Photo).To(Equal(bluePixelJPEGData))
										})
									})
								})

								Context("without previous post having a collage", func() {
									It("deletes the current post and creates a new backdated one", func() {
										fbAPI.On("PostDelete", facebookPageToken, facebookPostID).Return(nil)
										fbAPI.On("PagePhotoCreate", facebookPageToken, facebookPageID, mock.AnythingOfType("*model.Photo")).Return(&fbmodel.PhotoResponse{
											ID:     "whatever",
											PostID: facebookPageID + "_" + photoID,
										}, nil)
										groupPosts.On("UpdateByID", id, &model.OfferGroupPost{
											ID:                  id,
											Date:                date,
											MessageTemplate:     messageTemplate,
											FBPostID:            facebookPageID + "_" + photoID,
											PostedImageChecksum: bluePixelChecksum,
										}).Return(nil)

										err := facebookPost.Update(date, user, restaurant)
										Expect(err).To(BeNil())
										post := fbAPI.Calls[2].Arguments.Get(2).(*fbmodel.Photo)
										Expect(post.Message).To(Equal(messageTemplate + "\n\natitle - 5.67€\nbtitle - 4.67€"))
										Expect(post.Published).To(BeTrue())
										Expect(post.BackdatedTime).To(Equal(originalPublishTime))
										Expect(post.Photo).To(Equal(bluePixelJPEGData))
									})
								})
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
										ScheduledPublishTime: time.Date(2115, 01, 02, 8, 45, 0, 0, time.UTC),
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
									Expect(post.ScheduledPublishTime.Sub(time.Now())).To(BeNumerically("~", 11*time.Minute, time.Second))
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
									Expect(post.ScheduledPublishTime.Sub(time.Now())).To(BeNumerically("~", 11*time.Minute, time.Second))
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
