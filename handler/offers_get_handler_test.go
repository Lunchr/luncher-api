package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"time"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	"github.com/deiwin/luncher-api/geo"
	. "github.com/deiwin/luncher-api/handler"
	"github.com/deiwin/luncher-api/router"
	"github.com/deiwin/luncher-api/storage"
	"github.com/julienschmidt/httprouter"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RegionOffersHandler", func() {

	var (
		offersCollection db.Offers
		imageStorage     storage.Images
	)

	BeforeEach(func() {
		offersCollection = &mockOffers{}
		imageStorage = mockImageStorage{}
	})

	Describe("ProximalOffers", func() {
		var (
			handler router.Handler
		)

		JustBeforeEach(func() {
			handler = ProximalOffers(offersCollection, imageStorage)
		})

		Context("with no location specified", func() {
			It("should fail", func(done Done) {
				defer close(done)
				err := handler(responseRecorder, request)
				Expect(err.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("with no lat specified", func() {
			BeforeEach(func() {
				requestQuery = url.Values{
					"lng": {"25.55"},
				}
			})

			It("should fail", func(done Done) {
				defer close(done)
				err := handler(responseRecorder, request)
				Expect(err.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("with no lng specified", func() {
			BeforeEach(func() {
				requestQuery = url.Values{
					"lat": {"25.55"},
				}
			})

			It("should fail", func(done Done) {
				defer close(done)
				err := handler(responseRecorder, request)
				Expect(err.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("with a non-float lat specified", func() {
			BeforeEach(func() {
				requestQuery = url.Values{
					"lat": {"wut"},
					"lng": {"25.55"},
				}
			})

			It("should fail", func(done Done) {
				defer close(done)
				err := handler(responseRecorder, request)
				Expect(err.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("with a proper location specified", func() {
			BeforeEach(func() {
				requestQuery = url.Values{
					"lat": {"58.380094"},
					"lng": {"26.722691"},
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

	Describe("RegionOffers", func() {
		var (
			handler           router.HandlerWithParams
			regionsCollection db.Regions
			params            httprouter.Params
		)

		BeforeEach(func() {
			params = httprouter.Params{httprouter.Param{
				Key:   "name",
				Value: "",
			}}
			regionsCollection = &mockRegions{}
		})

		JustBeforeEach(func() {
			handler = RegionOffers(offersCollection, regionsCollection, imageStorage)
		})

		Context("with no region specified", func() {
			It("should fail", func(done Done) {
				defer close(done)
				err := handler(responseRecorder, request, params)
				Expect(err.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("with region specified", func() {
			BeforeEach(func() {
				params[0].Value = "Tartu"
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
				})

				It("should write the returned data to responsewriter", func(done Done) {
					defer close(done)
					handler(responseRecorder, request, params)
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
					err := handler(responseRecorder, request, params)
					Expect(err.Code).To(Equal(http.StatusInternalServerError))
				})
			})
		})
	})
})

func (m mockOffers) GetForRegion(region string, startTime, endTime time.Time) (offers []*model.Offer, err error) {
	Expect(region).To(Equal("Tartu"))
	if m.getForTimeRangeFunc != nil {
		offers, err = m.getForTimeRangeFunc(startTime, endTime)
	}
	return
}

func (m mockOffers) GetNear(loc geo.Location, startTime, endTime time.Time) (offers []*model.Offer, err error) {
	Expect(loc.Lat).To(BeNumerically("~", 58.380094))
	Expect(loc.Lng).To(BeNumerically("~", 26.722691))
	if m.getForTimeRangeFunc != nil {
		offers, err = m.getForTimeRangeFunc(startTime, endTime)
	}
	return
}
