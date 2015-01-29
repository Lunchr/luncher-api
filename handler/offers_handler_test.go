package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"time"

	"github.com/deiwin/facebook"
	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	. "github.com/deiwin/luncher-api/handler"
	"github.com/deiwin/luncher-api/session"

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
			// TODO the header assertion could be made a custom matcher
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
			sessionManager = &mockSessionManager{}
		})

		JustBeforeEach(func() {
			handler = PostOffers(offersCollection, usersCollection, sessionManager, authenticator)
		})

		It("should return HTTP 201: Created", func(done Done) {
			defer close(done)
			handler.ServeHTTP(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusCreated))
		})
	})
})

type mockUsers struct {
	db.Users
}

type mockOffers struct {
	getForTimeRangeFunc func(time.Time, time.Time) ([]*model.Offer, error)
	db.Offers
}

func (mock mockOffers) GetForTimeRange(startTime time.Time, endTime time.Time) (offers []*model.Offer, err error) {
	if mock.getForTimeRangeFunc != nil {
		offers, err = mock.getForTimeRangeFunc(startTime, endTime)
	}
	return
}
