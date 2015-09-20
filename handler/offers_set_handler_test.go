package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/db/model"
	. "github.com/Lunchr/luncher-api/handler"
	"github.com/Lunchr/luncher-api/handler/mocks"
	"github.com/Lunchr/luncher-api/router"
	"github.com/Lunchr/luncher-api/session"
	"github.com/deiwin/facebook"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/mock"
	"golang.org/x/oauth2"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OffersHandler", func() {

	var (
		offersCollection  db.Offers
		imageStorage      *mocks.Images
		regionsCollection *mocks.Regions
	)

	BeforeEach(func() {
		offersCollection = &mockOffers{}
		imageStorage = new(mocks.Images)
		imageStorage.On("ChecksumDataURL", "image data url").Return("image checksum", nil)
		imageStorage.On("HasChecksum", "image checksum").Return(false, nil)
		imageStorage.On("StoreDataURL", "image data url").Return(nil)
		imageStorage.On("PathsFor", "image checksum").Return(&model.OfferImagePaths{
			Large:     "images/a large image path",
			Thumbnail: "images/thumbnail",
		}, nil)
		regionsCollection = new(mocks.Regions)
		regionsCollection.On("GetName", "Tartu").Return(&model.Region{
			Name:     "Tartu",
			Location: "Europe/Tallinn",
		}, nil)
	})

	Describe("PostOffers", func() {
		var (
			usersCollection       db.Users
			restaurantsCollection *mocks.Restaurants
			handler               router.HandlerWithParams
			params                httprouter.Params
			sessionManager        session.Manager
			facebookPost          *mocks.Post
		)

		BeforeEach(func() {
			usersCollection = &mockUsers{}
			restaurantsCollection = new(mocks.Restaurants)
			facebookPost = new(mocks.Post)

			restaurantID := bson.ObjectId("12letrrestid")
			restaurant := &model.Restaurant{
				ID:      restaurantID,
				Name:    "Asian Chef",
				Region:  "Tartu",
				Address: "an-address",
				Location: model.Location{
					Type:        "Point",
					Coordinates: []float64{26.7, 58.4},
				},
				Phone: "+372 5678 910",
			}
			params = httprouter.Params{httprouter.Param{
				Key:   "restaurantID",
				Value: restaurantID.Hex(),
			}}
			restaurantsCollection.On("GetID", restaurantID).Return(restaurant, nil).Once()
			facebookPost = new(mocks.Post)
			facebookPost.On("Update", model.DateWithoutTime("2014-11-11"), mock.AnythingOfType("*model.User"), restaurant).Return(nil)
		})

		JustBeforeEach(func() {
			handler = PostOffers(offersCollection, usersCollection, restaurantsCollection, sessionManager,
				imageStorage, facebookPost, regionsCollection)
		})

		ExpectUserToBeLoggedIn(func() *router.HandlerError {
			return handler(responseRecorder, request, params)
		}, func(mgr session.Manager, users db.Users) {
			sessionManager = mgr
			usersCollection = users
		})

		Context("with session set and a matching user in DB", func() {
			BeforeEach(func() {
				sessionManager = &mockSessionManager{isSet: true, id: "correctSession"}
				requestMethod = "POST"
				requestData = map[string]interface{}{
					"title":       "thetitle",
					"ingredients": []string{"ingredient1", "ingredient2", "ingredient3"},
					"tags":        []string{"tag1", "tag2"},
					"price":       123.58,
					"from_time":   "2014-11-11T09:00:00.000Z",
					"to_time":     "2014-11-11T11:00:00.000Z",
					"image_data":  "image data url",
				}
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

			It("should include the offer with the new ID", func() {
				handler(responseRecorder, request, params)
				var offer model.OfferJSON
				json.Unmarshal(responseRecorder.Body.Bytes(), &offer)
				Expect(offer.ID).To(Equal(objectID))
				Expect(offer.Image.Large).To(Equal("images/a large image path"))
			})

			Context("with the image having already been stored once", func() {
				BeforeEach(func() {
					imageStorage.ExpectedCalls = make([]*mock.Call, 0)
					imageStorage.On("ChecksumDataURL", "image data url").Return("image checksum", nil)
					imageStorage.On("HasChecksum", "image checksum").Return(true, nil)
					imageStorage.On("StoreDataURL", "image data url").Return(errors.New("already stored"))
					imageStorage.On("PathsFor", "image checksum").Return(&model.OfferImagePaths{
						Large:     "images/a large image path",
						Thumbnail: "images/thumbnail",
					}, nil)
				})

				It("succeeds", func() {
					err := handler(responseRecorder, request, params)
					Expect(err).To(BeNil())
				})

				It("includes the image paths in the response", func() {
					handler(responseRecorder, request, params)
					var offer model.OfferJSON
					json.Unmarshal(responseRecorder.Body.Bytes(), &offer)
					Expect(offer.ID).To(Equal(objectID))
					Expect(offer.Image.Large).To(Equal("images/a large image path"))
					Expect(offer.Image.Thumbnail).To(Equal("images/thumbnail"))
				})
			})
		})
	})

	Describe("PutOffers", func() {
		var (
			usersCollection       db.Users
			restaurantsCollection *mocks.Restaurants
			handler               router.HandlerWithParams
			sessionManager        session.Manager
			facebookPost          *mocks.Post
			params                httprouter.Params
			restaurantID          bson.ObjectId
		)

		BeforeEach(func() {
			usersCollection = &mockUsers{}
			restaurantsCollection = new(mocks.Restaurants)

			restaurantID = bson.ObjectId("12letrrestid")
			params = httprouter.Params{httprouter.Param{
				Key:   "id",
				Value: objectID.Hex(),
			}, httprouter.Param{
				Key:   "restaurantID",
				Value: restaurantID.Hex(),
			}}

			restaurant := &model.Restaurant{
				ID:      restaurantID,
				Name:    "Asian Chef",
				Region:  "Tartu",
				Address: "an-address",
				Location: model.Location{
					Type:        "Point",
					Coordinates: []float64{26.7, 58.4},
				},
				Phone: "+372 5678 910",
			}
			restaurantsCollection.On("GetID", restaurantID).Return(restaurant, nil).Once()
			facebookPost = new(mocks.Post)
			facebookPost.On("Update", model.DateWithoutTime("2014-11-11"), mock.AnythingOfType("*model.User"), restaurant).Return(nil)
		})

		JustBeforeEach(func() {
			handler = PutOffers(offersCollection, usersCollection, restaurantsCollection, sessionManager,
				imageStorage, facebookPost, regionsCollection)
		})

		ExpectUserToBeLoggedIn(func() *router.HandlerError {
			return handler(responseRecorder, request, nil)
		}, func(mgr session.Manager, users db.Users) {
			sessionManager = mgr
			usersCollection = users
		})

		Context("with no matching offer in the DB", func() {
			BeforeEach(func() {
				sessionManager = &mockSessionManager{isSet: true, id: "correctSession"}
			})

			It("should fail", func() {
				err := handler(responseRecorder, request, params)
				Expect(err.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("with an ID that's not an object ID", func() {
			BeforeEach(func() {
				params = httprouter.Params{httprouter.Param{
					Key:   "id",
					Value: "not a proper bson.ObjectId",
				}, httprouter.Param{
					Key:   "restaurantID",
					Value: restaurantID.Hex(),
				}}
				sessionManager = &mockSessionManager{isSet: true, id: "correctSession"}
			})

			It("should fail", func() {
				err := handler(responseRecorder, request, params)
				Expect(err.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("with image not changed", func() {
			BeforeEach(func() {
				sessionManager = &mockSessionManager{isSet: true, id: "correctSession"}
				requestMethod = "PUT"
				requestData = map[string]interface{}{
					"title":       "thetitle",
					"ingredients": []string{"ingredient1", "ingredient2", "ingredient3"},
					"tags":        []string{"tag1", "tag2"},
					"price":       123.58,
					"from_time":   "2014-11-11T09:00:00.000Z",
					"to_time":     "2014-11-11T11:00:00.000Z",
					"image": map[string]interface{}{
						"large":     "images/a large image path",
						"thumbnail": "images/a thumbnail path",
					},
				}
				currentOffer := &model.Offer{
					CommonOfferFields: model.CommonOfferFields{
						ID:       objectID2,
						Title:    "an offer title",
						FromTime: time.Date(2014, 11, 11, 9, 0, 0, 0, time.UTC),
					},
					ImageChecksum: "image checksum",
				}
				offersCollection = &mockOffers{
					mockOffer:        currentOffer,
					imageIsUnchanged: true,
				}
				imageStorage.On("PathsFor", "").Return(&model.OfferImagePaths{}, nil)
			})

			It("should succeed", func() {
				err := handler(responseRecorder, request, params)
				Expect(err).To(BeNil())
			})
		})

		Context("with session set, a matching user in DB and an offer in DB", func() {
			BeforeEach(func() {
				sessionManager = &mockSessionManager{isSet: true, id: "correctSession"}
				requestMethod = "PUT"
				requestData = map[string]interface{}{
					"title":       "thetitle",
					"ingredients": []string{"ingredient1", "ingredient2", "ingredient3"},
					"tags":        []string{"tag1", "tag2"},
					"price":       123.58,
					"from_time":   "2014-11-11T09:00:00.000Z",
					"to_time":     "2014-11-11T11:00:00.000Z",
					"image_data":  "image data url",
				}
				currentOffer := &model.Offer{
					CommonOfferFields: model.CommonOfferFields{
						ID:       objectID,
						Title:    "an offer title",
						FromTime: time.Date(2014, 11, 11, 9, 0, 0, 0, time.UTC),
					},
					ImageChecksum: "image checksum",
				}
				offersCollection = &mockOffers{
					mockOffer: currentOffer,
				}
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

			It("should include the updated offer in the response", func() {
				handler(responseRecorder, request, params)
				var offer *model.OfferJSON
				json.Unmarshal(responseRecorder.Body.Bytes(), &offer)
				Expect(offer.ID).To(Equal(objectID))
				Expect(offer.Image.Large).To(Equal("images/a large image path"))
			})

			Context("with the image having already been stored once", func() {
				BeforeEach(func() {
					imageStorage.ExpectedCalls = make([]*mock.Call, 0)
					imageStorage.On("ChecksumDataURL", "image data url").Return("image checksum", nil)
					imageStorage.On("HasChecksum", "image checksum").Return(true, nil)
					imageStorage.On("StoreDataURL", "image data url").Return(errors.New("already stored"))
					imageStorage.On("PathsFor", "image checksum").Return(&model.OfferImagePaths{
						Large:     "images/a large image path",
						Thumbnail: "images/thumbnail",
					}, nil)
				})

				It("succeeds", func() {
					err := handler(responseRecorder, request, params)
					Expect(err).To(BeNil())
				})

				It("includes the image paths in the response", func() {
					handler(responseRecorder, request, params)
					var offer model.OfferJSON
					json.Unmarshal(responseRecorder.Body.Bytes(), &offer)
					Expect(offer.ID).To(Equal(objectID))
					Expect(offer.Image.Large).To(Equal("images/a large image path"))
					Expect(offer.Image.Thumbnail).To(Equal("images/thumbnail"))
				})
			})

			Describe("updating group post", func() {
				AfterEach(func() {
					facebookPost.AssertExpectations(GinkgoT())
				})

				BeforeEach(func() {
					facebookPost = new(mocks.Post)
					facebookPost.On("Update", model.DateWithoutTime("2014-11-11"), mock.AnythingOfType("*model.User"), mock.AnythingOfType("*model.Restaurant")).Return(nil)
				})

				It("succeeds", func() {
					err := handler(responseRecorder, request, params)
					Expect(err).NotTo(HaveOccurred())
				})

				Context("with offer date changed", func() {
					BeforeEach(func() {
						requestData = map[string]interface{}{
							"title":       "thetitle",
							"ingredients": []string{"ingredient1", "ingredient2", "ingredient3"},
							"tags":        []string{"tag1", "tag2"},
							"price":       123.58,
							"from_time":   "2014-11-15T09:00:00.000Z",
							"to_time":     "2014-11-15T11:00:00.000Z",
							"image_data":  "image data url",
						}
						date := model.DateWithoutTime("2014-11-15")
						facebookPost.On("Update", date, mock.AnythingOfType("*model.User"), mock.AnythingOfType("*model.Restaurant")).Return(nil)
					})

					It("succeeds updating the group posts for both the previous and new day", func() {
						err := handler(responseRecorder, request, params)
						Expect(err).NotTo(HaveOccurred())
					})
				})
			})
		})
	})

	Describe("DeleteOffers", func() {
		var (
			usersCollection       db.Users
			handler               router.HandlerWithParams
			sessionManager        session.Manager
			restaurantsCollection *mocks.Restaurants
			facebookPost          *mocks.Post
			params                httprouter.Params
			restaurantID          bson.ObjectId
		)

		BeforeEach(func() {
			usersCollection = &mockUsers{}
			restaurantsCollection = new(mocks.Restaurants)

			restaurantID = bson.ObjectId("12letrrestid")
			params = httprouter.Params{httprouter.Param{
				Key:   "id",
				Value: objectID.Hex(),
			}, httprouter.Param{
				Key:   "restaurantID",
				Value: restaurantID.Hex(),
			}}

			restaurant := &model.Restaurant{
				ID:      restaurantID,
				Name:    "Asian Chef",
				Region:  "Tartu",
				Address: "an-address",
				Location: model.Location{
					Type:        "Point",
					Coordinates: []float64{26.7, 58.4},
				},
				Phone: "+372 5678 910",
			}
			restaurantsCollection.On("GetID", restaurantID).Return(restaurant, nil).Once()
			facebookPost = new(mocks.Post)
			facebookPost.On("Update", model.DateWithoutTime("2014-11-11"), mock.AnythingOfType("*model.User"), restaurant).Return(nil)
		})

		JustBeforeEach(func() {
			handler = DeleteOffers(offersCollection, usersCollection, sessionManager, restaurantsCollection, facebookPost, regionsCollection)
		})

		ExpectUserToBeLoggedIn(func() *router.HandlerError {
			return handler(responseRecorder, request, nil)
		}, func(mgr session.Manager, users db.Users) {
			sessionManager = mgr
			usersCollection = users
		})

		Context("with no matching offer in the DB", func() {
			BeforeEach(func() {
				sessionManager = &mockSessionManager{isSet: true, id: "correctSession"}
			})

			It("should fail", func() {
				err := handler(responseRecorder, request, params)
				Expect(err.Code).To(Equal(http.StatusNotFound))
			})
		})

		Context("with an ID that's not an object ID", func() {
			BeforeEach(func() {
				params = httprouter.Params{httprouter.Param{
					Key:   "id",
					Value: "not a proper bson.ObjectId",
				}, httprouter.Param{
					Key:   "restaurantID",
					Value: restaurantID.Hex(),
				}}
				sessionManager = &mockSessionManager{isSet: true, id: "correctSession"}
			})

			It("should fail", func() {
				err := handler(responseRecorder, request, params)
				Expect(err.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("with session set, a matching user in DB and an offer in DB", func() {
			BeforeEach(func() {
				sessionManager = &mockSessionManager{isSet: true, id: "correctSession"}
				requestMethod = "DELETE"
				currentOffer := &model.Offer{
					CommonOfferFields: model.CommonOfferFields{
						ID:       objectID,
						Title:    "an offer title",
						FromTime: time.Date(2014, 11, 11, 0, 0, 0, 0, time.UTC),
					},
					ImageChecksum: "image checksum",
				}
				offersCollection = &mockOffers{
					mockOffer: currentOffer,
				}
			})

			It("should succeed", func() {
				err := handler(responseRecorder, request, params)
				Expect(err).To(BeNil())
			})
		})
	})
})

var objectID = bson.NewObjectId()
var objectID2 = bson.NewObjectId()

type mockUsers struct {
	db.Users
}

func (m mockUsers) GetSessionID(session string) (*model.User, error) {
	if session != "correctSession" {
		return nil, mgo.ErrNotFound
	}
	user := &model.User{
		ID:            objectID,
		RestaurantIDs: []bson.ObjectId{"12letrrestid"},
		Session: model.UserSession{
			FacebookUserToken: oauth2.Token{
				AccessToken: "usertoken",
			},
		},
	}
	return user, nil
}

type mockOffers struct {
	getForTimeRangeFunc func(time.Time, time.Time) ([]*model.Offer, error)
	mockOffer           *model.Offer
	imageIsUnchanged    bool
	db.Offers
}

func (m mockOffers) Insert(offers ...*model.Offer) ([]*model.Offer, error) {
	Expect(offers).To(HaveLen(1))
	offer := offers[0]
	Expect(offer.Title).To(Equal("thetitle"))
	Expect(offer.Ingredients).To(HaveLen(3))
	Expect(offer.Ingredients).To(ContainElement("ingredient1"))
	Expect(offer.Ingredients).To(ContainElement("ingredient2"))
	Expect(offer.Ingredients).To(ContainElement("ingredient3"))
	Expect(offer.Tags).To(HaveLen(2))
	Expect(offer.Tags).To(ContainElement("tag1"))
	Expect(offer.Tags).To(ContainElement("tag2"))
	Expect(offer.Price).To(BeNumerically("~", 123.58))
	Expect(offer.Restaurant.ID).To(Equal(bson.ObjectId("12letrrestid")))
	Expect(offer.Restaurant.Name).To(Equal("Asian Chef"))
	Expect(offer.Restaurant.Region).To(Equal("Tartu"))
	Expect(offer.Restaurant.Address).To(Equal("an-address"))
	Expect(offer.Restaurant.Location.Coordinates[0]).To(BeNumerically("~", 26.7))
	Expect(offer.Restaurant.Location.Coordinates[1]).To(BeNumerically("~", 58.4))
	Expect(offer.Restaurant.Phone).To(Equal("+372 5678 910"))
	Expect(offer.FromTime).To(Equal(time.Date(2014, 11, 11, 9, 0, 0, 0, time.UTC)))
	Expect(offer.ToTime).To(Equal(time.Date(2014, 11, 11, 11, 0, 0, 0, time.UTC)))
	Expect(offer.ImageChecksum).To(Equal("image checksum"))

	offers[0].ID = objectID
	return offers, nil
}

func (m mockOffers) UpdateID(id bson.ObjectId, offer *model.Offer) error {
	Expect(offer.Title).To(Equal("thetitle"))
	Expect(offer.Ingredients).To(HaveLen(3))
	Expect(offer.Ingredients).To(ContainElement("ingredient1"))
	Expect(offer.Ingredients).To(ContainElement("ingredient2"))
	Expect(offer.Ingredients).To(ContainElement("ingredient3"))
	Expect(offer.Tags).To(HaveLen(2))
	Expect(offer.Tags).To(ContainElement("tag1"))
	Expect(offer.Tags).To(ContainElement("tag2"))
	Expect(offer.Price).To(BeNumerically("~", 123.58))
	Expect(offer.Restaurant.ID).To(Equal(bson.ObjectId("12letrrestid")))
	Expect(offer.Restaurant.Name).To(Equal("Asian Chef"))
	Expect(offer.Restaurant.Region).To(Equal("Tartu"))
	Expect(offer.Restaurant.Address).To(Equal("an-address"))
	Expect(offer.Restaurant.Location.Coordinates[0]).To(BeNumerically("~", 26.7))
	Expect(offer.Restaurant.Location.Coordinates[1]).To(BeNumerically("~", 58.4))
	Expect(offer.Restaurant.Phone).To(Equal("+372 5678 910"))
	if id == objectID {
		Expect(offer.ImageChecksum).To(Equal("image checksum"))
	} else if id == objectID2 {
		Expect(offer.ImageChecksum).To(Equal(""))
	} else {
		Fail("Unexpected id")
	}
	return nil
}

func (m mockOffers) RemoveID(id bson.ObjectId) error {
	Expect(id).To(Equal(objectID))
	return nil
}

func (m mockOffers) GetID(id bson.ObjectId) (*model.Offer, error) {
	Expect(id).To(Equal(objectID))
	if m.mockOffer == nil {
		return nil, errors.New("offer not found")
	}
	return m.mockOffer, nil
}

func (m mockRestaurants) GetID(id bson.ObjectId) (*model.Restaurant, error) {
	Expect(id).To(Equal(bson.ObjectId("12letrrestid")))
	restaurant := &model.Restaurant{
		ID:      id,
		Name:    "Asian Chef",
		Region:  "Tartu",
		Address: "an-address",
		Location: model.Location{
			Type:        "Point",
			Coordinates: []float64{26.7, 58.4},
		},
		Phone: "+372 5678 910",
	}
	return restaurant, nil
}

func (m mockAuthenticator) APIConnection(tok *oauth2.Token) facebook.API {
	Expect(tok.AccessToken).To(Equal("usertoken"))
	return m.api
}

type mockAPI struct {
	shouldFail bool
	message    string
	facebook.API
}

func (m mockAPI) PostDelete(pageAccessToken, postID string) error {
	if m.shouldFail {
		return errors.New("post to FB failed")
	}
	Expect(pageAccessToken).To(Equal("pagetoken"))
	Expect(postID).To(Equal("fb post id"))
	return nil
}

type mockRegions struct {
	getAllFunc func() db.RegionIter
	db.Regions
}

func (m mockRegions) GetName(name string) (*model.Region, error) {
	Expect(name).To(Equal("Tartu"))
	region := &model.Region{
		Name:     "Tartu",
		Location: "Europe/Tallinn",
	}
	return region, nil
}
