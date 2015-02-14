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
			handler Handler
		)

		JustBeforeEach(func() {
			handler = Offers(offersCollection)
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
					func(startTime time.Time, endTime time.Time) (offers []*model.Offer, err error) {
						err = dbErr
						return
					},
					nil,
				}
			})

			It("should return error 500", func(done Done) {
				defer close(done)
				handler.ServeHTTP(responseRecorder, request)
				Expect(responseRecorder.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})

	Describe("PostOffers", func() {
		var (
			usersCollection db.Users
			handler         Handler
			authenticator   facebook.Authenticator
			sessionManager  session.Manager
		)

		BeforeEach(func() {
			usersCollection = &mockUsers{}
			authenticator = &mockAuthenticator{}
		})

		JustBeforeEach(func() {
			handler = PostOffers(offersCollection, usersCollection, sessionManager, authenticator)
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
					"message": {"postmessage"},
				}
			})

			Context("with post to FB failing", func() {
				BeforeEach(func() {
					authenticator = &mockAuthenticator{
						apiCallsShouldFail: true,
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
		})
	})
})

type mockUsers struct {
	db.Users
}

func (m mockUsers) GetBySessionID(session string) (*model.User, error) {
	if session != "correctSession" {
		return nil, errors.New("wrong session")
	}
	user := &model.User{
		FacebookPageID: "pageid",
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

func (m mockOffers) GetForTimeRange(startTime time.Time, endTime time.Time) (offers []*model.Offer, err error) {
	if m.getForTimeRangeFunc != nil {
		offers, err = m.getForTimeRangeFunc(startTime, endTime)
	}
	return
}

func (m mockAuthenticator) APIConnection(tok *oauth2.Token) facebook.API {
	Expect(tok.AccessToken).To(Equal("usertoken"))
	return &mockAPI{
		shouldFail: m.apiCallsShouldFail,
	}
}

type mockAPI struct {
	shouldFail bool
	facebook.API
}

func (m mockAPI) PagePublish(pageAccessToken, pageID, message string) (*fbmodel.Post, error) {
	Expect(pageAccessToken).To(Equal("pagetoken"))
	Expect(pageID).To(Equal("pageid"))
	Expect(message).To(Equal("postmessage"))

	if m.shouldFail {
		return nil, errors.New("post to FB failed")
	}

	post := &fbmodel.Post{
		ID: "postid",
	}
	return post, nil
}
