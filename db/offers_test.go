package db_test

import (
	"time"

	"github.com/deiwin/luncher-api/db/model"
	"github.com/deiwin/luncher-api/geo"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2/bson"
)

var _ = Describe("Offers", func() {
	var (
		earliestTime = time.Date(2014, 11, 10, 0, 0, 0, 0, time.UTC)
		latestTime   = time.Date(2014, 11, 13, 0, 0, 0, 0, time.UTC)
	)

	var anOffer = func() *model.Offer {
		return &model.Offer{
			// The location is needed because otherwise the index will complain
			Restaurant: model.OfferRestaurant{
				Location: model.Location{
					Type: "Point",
					// Somewhere distant so these entries wouldn't affect $near tests
					Coordinates: []float64{89, 89},
				},
			},
		}
	}

	Describe("Insert", func() {
		RebuildDBAfterEach()
		It("should return the offers with new IDs", func(done Done) {
			defer close(done)
			offers, err := offersCollection.Insert(anOffer(), anOffer())
			Expect(err).NotTo(HaveOccurred())
			Expect(offers).To(HaveLen(2))
			Expect(offers[0].ID).NotTo(BeEmpty())
			Expect(offers[1].ID).NotTo(BeEmpty())
		})

		It("should keep current ID if exists", func(done Done) {
			defer close(done)
			id := bson.NewObjectId()
			offer := anOffer()
			offer.ID = id
			offers, err := offersCollection.Insert(offer, anOffer())
			Expect(err).NotTo(HaveOccurred())
			Expect(offers).To(HaveLen(2))
			Expect(offers[0].ID).To(Equal(id))
			Expect(offers[1].ID).NotTo(Equal(id))
		})
	})

	Describe("UpdateID", func() {
		RebuildDBAfterEach()
		It("should fail for a non-existent ID", func(done Done) {
			defer close(done)
			err := offersCollection.UpdateID(bson.NewObjectId(), anOffer())
			Expect(err).To(HaveOccurred())
		})

		Context("with an offer with known ID inserted", func() {
			var id bson.ObjectId
			BeforeEach(func(done Done) {
				defer close(done)
				id = bson.NewObjectId()
				offer := anOffer()
				offer.ID = id
				_, err := offersCollection.Insert(offer)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should update the offer in DB", func(done Done) {
				defer close(done)
				offer := anOffer()
				offer.Title = "an updated title"
				err := offersCollection.UpdateID(id, offer)
				Expect(err).NotTo(HaveOccurred())
				result, err := offersCollection.GetID(id)
				Expect(err).NotTo(HaveOccurred())
				Expect(result.Title).To(Equal("an updated title"))
			})
		})
	})

	var ItHandlesStartAndEndTime = func(getOffers func(startTime, endTime time.Time) ([]*model.Offer, error)) {
		var (
			startTime             time.Time
			endTime               time.Time
			allOffers             []*model.Offer
			allOffersContainsMock = func(n int) bool {
				mock := mocks.offers[n]
				for _, offer := range allOffers {
					if offer.Title == mock.Title {
						return true
					}
				}
				return false
			}
			ExpectToMaybeContainMock = func(offers []*model.Offer, n int) {
				if allOffersContainsMock(n) {
					Expect(offers).To(ContainOfferMock(n))
				} else {
					Expect(offers).NotTo(ContainOfferMock(n))
				}
			}
		)
		Describe("time range", func() {
			BeforeEach(func(done Done) {
				defer close(done)
				var err error
				allOffers, err = getOffers(earliestTime, latestTime)
				Expect(err).NotTo(HaveOccurred())
			})
			JustBeforeEach(func() {
				if startTime != (time.Time{}) {
					endTime = startTime.AddDate(0, 0, 1)
				}
			})

			Context("with simple day range", func() {
				BeforeEach(func() {
					startTime = time.Date(2014, 11, 10, 0, 0, 0, 0, time.UTC)
				})

				It("should get offers for the day", func(done Done) {
					defer close(done)
					defer GinkgoRecover()
					offers, err := getOffers(startTime, endTime)
					Expect(err).NotTo(HaveOccurred())
					ExpectToMaybeContainMock(offers, 0)
					Expect(offers).NotTo(ContainOfferMock(1))
					Expect(offers).NotTo(ContainOfferMock(2))
				})
			})

			Context("with simple day range without offers", func() {
				BeforeEach(func() {
					startTime = time.Date(2014, 11, 13, 0, 0, 0, 0, time.UTC)
				})

				It("should get 0 offers for non-existent date", func(done Done) {
					defer close(done)
					offers, err := getOffers(startTime, endTime)
					Expect(err).NotTo(HaveOccurred())
					Expect(offers).To(BeEmpty())
				})
			})

			Describe("limits", func() {
				Describe("lower limit", func() {
					Context("with date just before fromTime-24h", func() {
						BeforeEach(func() {
							startTime = time.Date(2014, 11, 11, 8, 59, 0, 0, time.UTC)
						})

						It("should not get any offers", func(done Done) {
							defer close(done)
							offers, err := getOffers(startTime, endTime)
							Expect(err).NotTo(HaveOccurred())
							Expect(offers).To(BeEmpty())
						})
					})

					Context("with date just after fromTime-24h", func() {
						BeforeEach(func() {
							startTime = time.Date(2014, 11, 11, 9, 01, 0, 0, time.UTC)
						})

						It("should get an offer", func(done Done) {
							defer close(done)
							offers, err := getOffers(startTime, endTime)
							Expect(err).NotTo(HaveOccurred())
							ExpectToMaybeContainMock(offers, 2)
						})
					})
				})

				Describe("higher limit", func() {
					Context("with date just before toTime", func() {
						BeforeEach(func() {
							startTime = time.Date(2014, 11, 12, 10, 59, 0, 0, time.UTC)
						})

						It("should get an offer", func(done Done) {
							defer close(done)
							offers, err := getOffers(startTime, endTime)
							Expect(err).NotTo(HaveOccurred())
							ExpectToMaybeContainMock(offers, 2)
						})
					})

					Context("with date just after fromTime-24h", func() {
						BeforeEach(func() {
							startTime = time.Date(2014, 11, 12, 11, 01, 0, 0, time.UTC)
						})

						It("should not get any offers", func(done Done) {
							defer close(done)
							offers, err := getOffers(startTime, endTime)
							Expect(err).NotTo(HaveOccurred())
							Expect(offers).To(BeEmpty())
						})
					})
				})
			})
		})
	}

	Describe("GetForRegion", func() {
		var (
			region string
		)

		Context("with region matching no offers", func() {
			BeforeEach(func() {
				region = "blablabla"
			})

			It("should get 0 offers", func(done Done) {
				defer close(done)
				offers, err := offersCollection.GetForRegion(region, earliestTime, latestTime)
				Expect(err).NotTo(HaveOccurred())
				Expect(offers).To(BeEmpty())
			})
		})

		Context("with region matching some of the offers", func() { // The 1st and 3rd offer
			BeforeEach(func() {
				region = "Tartu"
			})

			It("should get all offers for that region", func(done Done) {
				defer close(done)
				offers, err := offersCollection.GetForRegion(region, earliestTime, latestTime)
				Expect(err).NotTo(HaveOccurred())
				Expect(offers).To(HaveLen(2))
				Expect(offers).To(ContainOfferMock(0))
				Expect(offers).NotTo(ContainOfferMock(1))
				Expect(offers).To(ContainOfferMock(2))
			})

			ItHandlesStartAndEndTime(func(startTime, endTime time.Time) ([]*model.Offer, error) {
				return offersCollection.GetForRegion(region, startTime, endTime)
			})

			Context("with the region matching rest of the offers", func() {
				BeforeEach(func() {
					region = "Tallinn"
				})

				It("should get all offers for that region", func(done Done) {
					defer close(done)
					offers, err := offersCollection.GetForRegion(region, earliestTime, latestTime)
					Expect(err).NotTo(HaveOccurred())
					Expect(offers).To(HaveLen(1))
					Expect(offers).NotTo(ContainOfferMock(0))
					Expect(offers).To(ContainOfferMock(1))
					Expect(offers).NotTo(ContainOfferMock(2))
				})

			})
		})
	})

	Describe("GetNear", func() {
		var (
			loc geo.Location
		)

		Context("with location on top of one of the restaurants", func() {
			BeforeEach(func() {
				loc = geo.Location{
					Lat: 58.37,
					Lng: 26.72,
				}
			})

			It("should return close restaurants in order of proximity", func(done Done) {
				defer close(done)
				defer GinkgoRecover()
				offers, err := offersCollection.GetNear(loc, earliestTime, latestTime)
				Expect(err).NotTo(HaveOccurred())
				Expect(offers).To(HaveLen(2))
				Expect(offers[0].Title).To(Equal(mocks.offers[0].Title))
				Expect(offers[1].Title).To(Equal(mocks.offers[2].Title))
			})

			It("should include distances", func(done Done) {
				defer close(done)
				defer GinkgoRecover()
				offers, err := offersCollection.GetNear(loc, earliestTime, latestTime)
				Expect(err).NotTo(HaveOccurred())
				Expect(offers).To(HaveLen(2))
				Expect(offers[0].Distance).To(BeNumerically("~", 0, 1))
				Expect(offers[1].Distance).To(BeNumerically("~", 1257, 1))
			})

			ItHandlesStartAndEndTime(func(startTime, endTime time.Time) ([]*model.Offer, error) {
				offersWithDist, err := offersCollection.GetNear(loc, startTime, endTime)
				if err != nil {
					return nil, err
				}
				var offers = make([]*model.Offer, len(offersWithDist))
				for i, o := range offersWithDist {
					offers[i] = &o.Offer
				}
				return offers, nil
			})
		})

		Context("with location on top of a different restaurant", func() {
			BeforeEach(func() {
				loc = geo.Location{
					Lat: 58.36,
					Lng: 26.73,
				}
			})

			It("should return close restaurants in order of proximity", func(done Done) {
				defer close(done)
				defer GinkgoRecover()
				offers, err := offersCollection.GetNear(loc, earliestTime, latestTime)
				Expect(err).NotTo(HaveOccurred())
				Expect(offers).To(HaveLen(2))
				Expect(offers[0].Title).To(Equal(mocks.offers[2].Title))
				Expect(offers[1].Title).To(Equal(mocks.offers[0].Title))
			})
		})
	})

	Describe("GetForRestaurant", func() {
		var (
			startTime      time.Time
			restaurant     string
			offerStartTime = time.Date(2014, 11, 10, 9, 0, 0, 0, time.UTC)
			offerEndTime   = time.Date(2014, 11, 10, 11, 0, 0, 0, time.UTC)
		)

		Context("with an existing restaurant", func() {
			BeforeEach(func() {
				restaurant = "Asian Chef"
			})

			Context("with time right before an offer", func() {
				BeforeEach(func() {
					startTime = offerStartTime.Add(-1 * time.Second)
				})

				It("should include the offer", func() {
					offers, err := offersCollection.GetForRestaurant(restaurant, startTime)
					Expect(err).NotTo(HaveOccurred())
					Expect(offers).To(HaveLen(1))
					Expect(offers).To(ContainOfferMock(0))
				})
			})

			Context("with time right before an offer ends", func() {
				BeforeEach(func() {
					startTime = offerEndTime.Add(-1 * time.Second)
				})

				It("should include the offer", func() {
					offers, err := offersCollection.GetForRestaurant(restaurant, startTime)
					Expect(err).NotTo(HaveOccurred())
					Expect(offers).To(HaveLen(1))
					Expect(offers).To(ContainOfferMock(0))
				})
			})
			Context("with time right after an offer offer has ended", func() {
				BeforeEach(func() {
					startTime = offerEndTime.Add(1 * time.Second)
				})

				It("should NOT include the offer", func() {
					offers, err := offersCollection.GetForRestaurant(restaurant, startTime)
					Expect(err).NotTo(HaveOccurred())
					Expect(offers).To(HaveLen(0))
				})
			})
		})

		Context("with a non-existing restaurant", func() {
			BeforeEach(func() {
				restaurant = "something random"
				startTime = offerStartTime.Add(-24 * 10 * time.Hour)
			})

			It("should NOT include the offer", func() {
				offers, err := offersCollection.GetForRestaurant(restaurant, startTime)
				Expect(err).NotTo(HaveOccurred())
				Expect(offers).To(HaveLen(0))
			})
		})
	})
})
