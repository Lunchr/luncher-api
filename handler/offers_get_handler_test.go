package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	. "github.com/deiwin/luncher-api/handler"
	"github.com/deiwin/luncher-api/router"
	"github.com/deiwin/luncher-api/storage"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RegionOffersHandler", func() {

	var (
		offersCollection db.Offers
	)

	BeforeEach(func() {
		offersCollection = &mockOffers{}
	})

	Describe("Offers", func() {
		var (
			handler           router.Handler
			regionsCollection db.Regions
			imageStorage      storage.Images
		)

		BeforeEach(func() {
			regionsCollection = &mockRegions{}
		})

		JustBeforeEach(func() {
			handler = Offers(offersCollection, regionsCollection, imageStorage)
		})

		Context("with no region specified", func() {
			It("should fail", func(done Done) {
				defer close(done)
				err := handler(responseRecorder, request)
				Expect(err.Code).To(Equal(http.StatusBadRequest))
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

			Context("with simple mocked result from DB", func() {
				var (
					mockResult []*model.Offer
				)
				BeforeEach(func() {
					mockResult = []*model.Offer{
						&model.Offer{
							Title: "sometitle",
							Image: "image checksum",
						},
					}
					offersCollection = &mockOffers{
						getForTimeRangeFunc: func(startTime time.Time, endTime time.Time) (offers []*model.Offer, err error) {
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
					}
					imageStorage = mockImageStorage{}
				})

				It("should write the returned data to responsewriter", func(done Done) {
					defer close(done)
					handler(responseRecorder, request)
					var result []*model.Offer
					json.Unmarshal(responseRecorder.Body.Bytes(), &result)
					Expect(result).To(HaveLen(1))
					Expect(result[0].Title).To(Equal(mockResult[0].Title))
					Expect(result[0].Image).To(Equal("images/a large image path"))
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
					err := handler(responseRecorder, request)
					Expect(err.Code).To(Equal(http.StatusInternalServerError))
				})
			})
		})
	})
})
