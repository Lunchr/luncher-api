package handler_test

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/deiwin/praad-api/db"
	"github.com/deiwin/praad-api/db/model"
	. "github.com/deiwin/praad-api/handler"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

type handlerFunc func(http.ResponseWriter, *http.Request)

var _ = Describe("OffersHandler", func() {
	var (
		mockOffersCollection db.Offers
		handler              handlerFunc
	)

	BeforeEach(func() {
		mockOffersCollection = &mockOffers{}
	})

	JustBeforeEach(func() {
		handler = Offers(mockOffersCollection)
	})

	Describe("GetForTimeRange", func() {

		It("should return json", func(done Done) {
			defer close(done)
			handler(responseRecorder, request)
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
				mockOffersCollection = &mockOffers{
					func(startTime time.Time, endTime time.Time) (offers []*model.Offer, err error) {
						offers = mockResult
						return
					},
				}
			})

			It("should write the returned data to responsewriter", func(done Done) {
				defer close(done)
				handler(responseRecorder, request)
				// Expect(responseRecorder.Flushed).To(BeTrue()) // TODO check if this should be true
				var result []*model.Offer
				json.Unmarshal(responseRecorder.Body.Bytes(), &result)
				Expect(result).To(HaveLen(1))
				Expect(result[0].Title).To(Equal(mockResult[0].Title))
			})
		})
	})
})

type mockOffers struct {
	getForTimeRangeFunc func(time.Time, time.Time) ([]*model.Offer, error)
}

func (mock mockOffers) Insert(offersToInsert ...*model.Offer) (err error) {
	return
}

func (mock mockOffers) GetForTimeRange(startTime time.Time, endTime time.Time) (offers []*model.Offer, err error) {
	if mock.getForTimeRangeFunc != nil {
		offers, err = mock.getForTimeRangeFunc(startTime, endTime)
	}
	return
}
