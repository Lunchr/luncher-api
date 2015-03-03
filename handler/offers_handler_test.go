package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/deiwin/facebook"
	fbmodel "github.com/deiwin/facebook/model"
	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	. "github.com/deiwin/luncher-api/handler"
	"github.com/deiwin/luncher-api/session"
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

	Describe("Offers", func() {
		var (
			handler           Handler
			regionsCollection db.Regions
		)

		BeforeEach(func() {
			regionsCollection = &mockRegions{}
		})

		JustBeforeEach(func() {
			handler = Offers(offersCollection, regionsCollection)
		})

		Context("with no region specified", func() {
			It("should fail", func(done Done) {
				defer close(done)
				handler.ServeHTTP(responseRecorder, request)
				Expect(responseRecorder.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("with region specified", func() {
			BeforeEach(func() {
				requestQuery = url.Values{
					"region": {"Tartu"},
				}
			})

			It("should succeed", func(done Done) {
				defer close(done)
				handler.ServeHTTP(responseRecorder, request)
				Expect(responseRecorder.Code).To(Equal(http.StatusOK))
			})

			It("should return json", func(done Done) {
				defer close(done)
				handler.ServeHTTP(responseRecorder, request)
				contentTypes := responseRecorder.HeaderMap["Content-Type"]
				Expect(contentTypes).To(HaveLen(1))
				Expect(contentTypes[0]).To(Equal("application/json"))
			})

			Context("with simple mocked result from DB", func() {
				var (
					mockResult []*model.Offer
				)
				BeforeEach(func() {
					mockResult = []*model.Offer{&model.Offer{Title: "sometitle"}}
					offersCollection = &mockOffers{
						func(startTime time.Time, endTime time.Time) (offers []*model.Offer, err error) {
							duration := endTime.Sub(startTime)
							// Due to daylight saving etc it's not always exactly 24h, but
							// I think with +- 1h it should always pass.
							Expect(duration).To(BeNumerically("~", 24*time.Hour, time.Hour))

							loc, err := time.LoadLocation("Europe/Tallinn")
							Expect(err).NotTo(HaveOccurred())
							Expect(startTime.Location()).To(Equal(loc))
							Expect(endTime.Location()).To(Equal(loc))

							offers = mockResult
							return
						},
						nil,
					}
				})

				It("should write the returned data to responsewriter", func(done Done) {
					defer close(done)
					handler.ServeHTTP(responseRecorder, request)
					// Expect(responseRecorder.Flushed).To(BeTrue()) // TODO check if this should be true
					var result []*model.Offer
					json.Unmarshal(responseRecorder.Body.Bytes(), &result)
					Expect(result).To(HaveLen(1))
					Expect(result[0].Title).To(Equal(mockResult[0].Title))
				})
			})

			Context("with an error returned from the DB", func() {
				var dbErr = errors.New("DB stuff failed")

				BeforeEach(func() {
					offersCollection = &mockOffers{
						getForTimeRangeFunc: func(startTime time.Time, endTime time.Time) (offers []*model.Offer, err error) {
							err = dbErr
							return
						},
					}
				})

				It("should return error 500", func(done Done) {
					defer close(done)
					handler.ServeHTTP(responseRecorder, request)
					Expect(responseRecorder.Code).To(Equal(http.StatusInternalServerError))
				})
			})
		})
	})

	Describe("PostOffers", func() {
		var (
			usersCollection       db.Users
			restaurantsCollection db.Restaurants
			handler               Handler
			authenticator         facebook.Authenticator
			sessionManager        session.Manager
		)

		BeforeEach(func() {
			usersCollection = &mockUsers{}
			restaurantsCollection = &mockRestaurants{}
			authenticator = &mockAuthenticator{}
		})

		JustBeforeEach(func() {
			handler = PostOffers(offersCollection, usersCollection, restaurantsCollection, sessionManager, authenticator)
		})
		Context("with no session set", func() {
			BeforeEach(func() {
				sessionManager = &mockSessionManager{}
			})

			It("should be forbidden", func(done Done) {
				defer close(done)
				handler.ServeHTTP(responseRecorder, request)
				Expect(responseRecorder.Code).To(Equal(http.StatusForbidden))
			})
		})

		Context("with session set, but no matching user in DB", func() {
			BeforeEach(func() {
				sessionManager = &mockSessionManager{isSet: true}
			})

			It("should be forbidden", func(done Done) {
				defer close(done)
				handler.ServeHTTP(responseRecorder, request)
				Expect(responseRecorder.Code).To(Equal(http.StatusForbidden))
			})
		})

		Context("with session set and a matching user in DB", func() {
			BeforeEach(func() {
				sessionManager = &mockSessionManager{isSet: true, id: "correctSession"}
				requestMethod = "POST"
				requestData = url.Values{
					"title":       {"thetitle"},
					"description": {"thedescription"},
					"tags":        {"tag1", "tag2"},
					"price":       {"123.58"},
					"from_time":   {"2014-11-11T09:00:00.000Z"},
					"to_time":     {"2014-11-11T11:00:00.000Z"},
				}
				authenticator = &mockAuthenticator{
					api: &mockAPI{
						message: "thetitle - thedescription",
					},
				}
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
					handler.ServeHTTP(responseRecorder, request)
					Expect(responseRecorder.Code).To(Equal(http.StatusBadGateway))
				})
			})

			It("should succeed", func(done Done) {
				defer close(done)
				handler.ServeHTTP(responseRecorder, request)
				Expect(responseRecorder.Code).To(Equal(http.StatusOK))
			})

			It("should return json", func(done Done) {
				defer close(done)
				handler.ServeHTTP(responseRecorder, request)
				contentTypes := responseRecorder.HeaderMap["Content-Type"]
				Expect(contentTypes).To(HaveLen(1))
				Expect(contentTypes[0]).To(Equal("application/json"))
			})

			It("should include the offer with the new ID", func(done Done) {
				defer close(done)
				handler.ServeHTTP(responseRecorder, request)
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

func (m mockUsers) GetBySessionID(session string) (*model.User, error) {
	if session != "correctSession" {
		return nil, errors.New("wrong session")
	}
	user := &model.User{
		FacebookPageID: "pageid",
		RestaurantID:   "restid",
		Session: model.UserSession{
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
	db.Offers
}

func (m mockOffers) Get(region string, startTime, endTime time.Time) (offers []*model.Offer, err error) {
	Expect(region).To(Equal("Tartu"))
	if m.getForTimeRangeFunc != nil {
		offers, err = m.getForTimeRangeFunc(startTime, endTime)
	}
	return
}

func (m mockOffers) Insert(offers ...*model.Offer) ([]*model.Offer, error) {
	Expect(offers).To(HaveLen(1))
	offer := offers[0]
	Expect(offer.FBPostID).To(Equal("postid"))
	Expect(offer.Title).To(Equal("thetitle"))
	Expect(offer.Description).To(Equal("thedescription"))
	Expect(offer.Tags).To(HaveLen(2))
	Expect(offer.Tags).To(ContainElement("tag1"))
	Expect(offer.Tags).To(ContainElement("tag2"))
	Expect(offer.Price).To(BeNumerically("~", 123.58))
	Expect(offer.Restaurant.Name).To(Equal("Asian Chef"))
	Expect(offer.Restaurant.Region).To(Equal("Tartu"))
	Expect(offer.FromTime).To(Equal(time.Date(2014, 11, 11, 9, 0, 0, 0, time.UTC)))
	Expect(offer.ToTime).To(Equal(time.Date(2014, 11, 11, 11, 0, 0, 0, time.UTC)))

	offers[0].ID = objectID
	return offers, nil
}

func (m mockRestaurants) GetByID(id bson.ObjectId) (*model.Restaurant, error) {
	Expect(id).To(Equal(bson.ObjectId("restid")))
	restaurant := &model.Restaurant{
		Name:   "Asian Chef",
		Region: "Tartu",
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
	db.Regions
}

func (m mockRegions) Get(name string) (*model.Region, error) {
	Expect(name).To(Equal("Tartu"))
	region := &model.Region{
		Name:     "Tartu",
		Location: "Europe/Tallinn",
	}
	return region, nil
}
