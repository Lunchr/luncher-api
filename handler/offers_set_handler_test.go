package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/deiwin/facebook"
	fbmodel "github.com/deiwin/facebook/model"
	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	. "github.com/deiwin/luncher-api/handler"
	"github.com/deiwin/luncher-api/router"
	"github.com/deiwin/luncher-api/session"
	"github.com/deiwin/luncher-api/storage"
	"github.com/julienschmidt/httprouter"
	"golang.org/x/oauth2"
	"gopkg.in/mgo.v2/bson"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OffersHandler", func() {

	var (
		offersCollection db.Offers
	)

	BeforeEach(func() {
		offersCollection = &mockOffers{}
	})

	Describe("PostOffers", func() {
		var (
			usersCollection       db.Users
			restaurantsCollection db.Restaurants
			handler               router.Handler
			authenticator         facebook.Authenticator
			sessionManager        session.Manager
			imageStorage          storage.Images
		)

		BeforeEach(func() {
			usersCollection = &mockUsers{}
			restaurantsCollection = &mockRestaurants{}
			authenticator = &mockAuthenticator{}
		})

		JustBeforeEach(func() {
			handler = PostOffers(offersCollection, usersCollection, restaurantsCollection, sessionManager, authenticator, imageStorage)
		})

		ExpectUserToBeLoggedIn(func() *router.HandlerError {
			return handler(responseRecorder, request)
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
					"image":       "image data url",
				}
				authenticator = &mockAuthenticator{
					api: &mockAPI{
						message: "thetitle - Ingredient1, ingredient2, ingredient3",
					},
				}
				imageStorage = mockImageStorage{}
			})

			Context("with post to FB failing", func() {
				BeforeEach(func() {
					authenticator = &mockAuthenticator{
						api: &mockAPI{
							message:    "postmessage",
							shouldFail: true,
						},
					}
				})

				It("should fail", func(done Done) {
					defer close(done)
					err := handler(responseRecorder, request)
					Expect(err.Code).To(Equal(http.StatusBadGateway))
				})
			})

			It("should succeed", func(done Done) {
				defer close(done)
				err := handler(responseRecorder, request)
				Expect(err).To(BeNil())
			})

			It("should return json", func(done Done) {
				defer close(done)
				handler(responseRecorder, request)
				contentTypes := responseRecorder.HeaderMap["Content-Type"]
				Expect(contentTypes).To(HaveLen(1))
				Expect(contentTypes[0]).To(Equal("application/json"))
			})

			It("should include the offer with the new ID", func(done Done) {
				defer close(done)
				handler(responseRecorder, request)
				var offer *model.Offer
				json.Unmarshal(responseRecorder.Body.Bytes(), &offer)
				Expect(offer.ID).To(Equal(objectID))
				Expect(offer.FBPostID).To(Equal("postid"))
			})
		})
	})

	Describe("PutOffers", func() {
		var (
			usersCollection       db.Users
			restaurantsCollection db.Restaurants
			handler               router.HandlerWithParams
			authenticator         facebook.Authenticator
			sessionManager        session.Manager
			imageStorage          storage.Images
			params                httprouter.Params
		)

		BeforeEach(func() {
			usersCollection = &mockUsers{}
			restaurantsCollection = &mockRestaurants{}
			authenticator = &mockAuthenticator{}
			params = httprouter.Params{httprouter.Param{
				Key:   "id",
				Value: objectID.Hex(),
			}}
		})

		JustBeforeEach(func() {
			handler = PutOffers(offersCollection, usersCollection, restaurantsCollection, sessionManager, authenticator, imageStorage)
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

			It("should fail", func(done Done) {
				defer close(done)
				err := handler(responseRecorder, request, params)
				Expect(err.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("with an ID that's not an object ID", func() {
			BeforeEach(func() {
				params = httprouter.Params{httprouter.Param{
					Key:   "id",
					Value: "not a proper bson.ObjectId",
				}}
				sessionManager = &mockSessionManager{isSet: true, id: "correctSession"}
			})

			It("should fail", func(done Done) {
				defer close(done)
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
					"image":       "images/a large image path",
				}
				currentOffer := &model.Offer{
					Title:    "an offer title",
					FBPostID: "fb post id",
					Image:    "image checksum",
				}
				offersCollection = &mockOffers{
					mockOffer:        currentOffer,
					imageIsUnchanged: true,
				}
				authenticator = &mockAuthenticator{
					api: &mockAPI{
						message: "thetitle - Ingredient1, ingredient2, ingredient3",
					},
				}
				imageStorage = mockImageStorage{}
			})

			It("should succeed", func(done Done) {
				defer close(done)
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
					"image":       "image data url",
				}
				currentOffer := &model.Offer{
					Title:    "an offer title",
					FBPostID: "fb post id",
					Image:    "image checksum",
				}
				offersCollection = &mockOffers{
					mockOffer: currentOffer,
				}
				authenticator = &mockAuthenticator{
					api: &mockAPI{
						message: "thetitle - Ingredient1, ingredient2, ingredient3",
					},
				}
				imageStorage = mockImageStorage{}
			})

			Context("with post to FB failing", func() {
				BeforeEach(func() {
					authenticator = &mockAuthenticator{
						api: &mockAPI{
							message:    "postmessage",
							shouldFail: true,
						},
					}
				})

				It("should fail", func(done Done) {
					defer close(done)
					err := handler(responseRecorder, request, params)
					Expect(err.Code).To(Equal(http.StatusBadGateway))
				})
			})

			It("should succeed", func(done Done) {
				defer close(done)
				err := handler(responseRecorder, request, params)
				Expect(err).To(BeNil())
			})

			It("should return json", func(done Done) {
				defer close(done)
				handler(responseRecorder, request, params)
				contentTypes := responseRecorder.HeaderMap["Content-Type"]
				Expect(contentTypes).To(HaveLen(1))
				Expect(contentTypes[0]).To(Equal("application/json"))
			})

			It("should include the updated offer in the response", func(done Done) {
				defer close(done)
				handler(responseRecorder, request, params)
				var offer *model.Offer
				json.Unmarshal(responseRecorder.Body.Bytes(), &offer)
				Expect(offer.ID).To(Equal(objectID))
				Expect(offer.FBPostID).To(Equal("postid"))
			})
		})
	})
})

var objectID = bson.NewObjectId()

type mockUsers struct {
	db.Users
}

func (m mockUsers) GetSessionID(session string) (*model.User, error) {
	if session != "correctSession" {
		return nil, errors.New("wrong session")
	}
	user := &model.User{
		ID:             objectID,
		FacebookPageID: "pageid",
		RestaurantID:   "restid",
		Session: &model.UserSession{
			FacebookUserToken: oauth2.Token{
				AccessToken: "usertoken",
			},
			FacebookPageToken: "pagetoken",
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
	Expect(offer.FBPostID).To(Equal("postid"))
	Expect(offer.Title).To(Equal("thetitle"))
	Expect(offer.Ingredients).To(HaveLen(3))
	Expect(offer.Ingredients).To(ContainElement("ingredient1"))
	Expect(offer.Ingredients).To(ContainElement("ingredient2"))
	Expect(offer.Ingredients).To(ContainElement("ingredient3"))
	Expect(offer.Tags).To(HaveLen(2))
	Expect(offer.Tags).To(ContainElement("tag1"))
	Expect(offer.Tags).To(ContainElement("tag2"))
	Expect(offer.Price).To(BeNumerically("~", 123.58))
	Expect(offer.Restaurant.Name).To(Equal("Asian Chef"))
	Expect(offer.Restaurant.Region).To(Equal("Tartu"))
	Expect(offer.Restaurant.Address).To(Equal("an-address"))
	Expect(offer.Restaurant.Location.Coordinates[0]).To(BeNumerically("~", 26.7))
	Expect(offer.Restaurant.Location.Coordinates[1]).To(BeNumerically("~", 58.4))
	Expect(offer.FromTime).To(Equal(time.Date(2014, 11, 11, 9, 0, 0, 0, time.UTC)))
	Expect(offer.ToTime).To(Equal(time.Date(2014, 11, 11, 11, 0, 0, 0, time.UTC)))
	Expect(offer.Image).To(Equal("image checksum"))

	offers[0].ID = objectID
	return offers, nil
}

func (m mockOffers) UpdateID(id bson.ObjectId, offer *model.Offer) error {
	Expect(id).To(Equal(objectID))
	Expect(offer.FBPostID).To(Equal("postid"))
	Expect(offer.Title).To(Equal("thetitle"))
	Expect(offer.Ingredients).To(HaveLen(3))
	Expect(offer.Ingredients).To(ContainElement("ingredient1"))
	Expect(offer.Ingredients).To(ContainElement("ingredient2"))
	Expect(offer.Ingredients).To(ContainElement("ingredient3"))
	Expect(offer.Tags).To(HaveLen(2))
	Expect(offer.Tags).To(ContainElement("tag1"))
	Expect(offer.Tags).To(ContainElement("tag2"))
	Expect(offer.Price).To(BeNumerically("~", 123.58))
	Expect(offer.Restaurant.Name).To(Equal("Asian Chef"))
	Expect(offer.Restaurant.Region).To(Equal("Tartu"))
	Expect(offer.Restaurant.Address).To(Equal("an-address"))
	Expect(offer.Restaurant.Location.Coordinates[0]).To(BeNumerically("~", 26.7))
	Expect(offer.Restaurant.Location.Coordinates[1]).To(BeNumerically("~", 58.4))
	Expect(offer.FromTime).To(Equal(time.Date(2014, 11, 11, 9, 0, 0, 0, time.UTC)))
	Expect(offer.ToTime).To(Equal(time.Date(2014, 11, 11, 11, 0, 0, 0, time.UTC)))
	Expect(offer.Image).To(Equal("image checksum"))
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
	Expect(id).To(Equal(bson.ObjectId("restid")))
	restaurant := &model.Restaurant{
		Name:    "Asian Chef",
		Region:  "Tartu",
		Address: "an-address",
		Location: model.Location{
			Type:        "Point",
			Coordinates: []float64{26.7, 58.4},
		},
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

func (m mockAPI) PagePublish(pageAccessToken, pageID, message string) (*fbmodel.Post, error) {
	if m.shouldFail {
		return nil, errors.New("post to FB failed")
	}

	Expect(pageAccessToken).To(Equal("pagetoken"))
	Expect(pageID).To(Equal("pageid"))
	Expect(message).To(Equal(m.message))

	post := &fbmodel.Post{
		ID: "postid",
	}
	return post, nil
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

type mockImageStorage struct {
	storage.Images
}

func (m mockImageStorage) ChecksumDataURL(dataURL string) (string, error) {
	Expect(dataURL).To(Equal("image data url"))
	return "image checksum", nil
}

func (m mockImageStorage) StoreDataURL(dataURL string) error {
	Expect(dataURL).To(Equal("image data url"))
	return nil
}

func (m mockImageStorage) PathForLarge(checksum string) (string, error) {
	Expect(checksum).To(Equal("image checksum"))
	return "images/a large image path", nil
}
