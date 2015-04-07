package db_test

import (
	"time"

	"github.com/deiwin/luncher-api/db/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2/bson"
)

var _ = Describe("Offers", func() {
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

	Describe("GetForRegion", func() {
		var (
			startTime    time.Time
			endTime      time.Time
			region       string
			earliestTime = time.Date(2014, 11, 10, 0, 0, 0, 0, time.UTC)
			latestTime   = time.Date(2014, 11, 13, 0, 0, 0, 0, time.UTC)
		)

		JustBeforeEach(func() {
			if startTime != (time.Time{}) {
				endTime = startTime.AddDate(0, 0, 1)
			}
		})

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

			Context("with simple day range", func() {
				BeforeEach(func() {
					startTime = time.Date(2014, 11, 10, 0, 0, 0, 0, time.UTC)
				})

				It("should get offers for the day", func(done Done) {
					defer close(done)
					offers, err := offersCollection.GetForRegion(region, startTime, endTime)
					Expect(err).NotTo(HaveOccurred())
					Expect(offers).To(HaveLen(1))
					Expect(offers).To(ContainOfferMock(0))
					Expect(offers).NotTo(ContainOfferMock(1))
					Expect(offers).NotTo(ContainOfferMock(2))
				})
			})

			Context("with simple day range withour offers", func() {
				BeforeEach(func() {
					startTime = time.Date(2014, 11, 13, 0, 0, 0, 0, time.UTC)
				})

				It("should get 0 offers for non-existent date", func(done Done) {
					defer close(done)
					offers, err := offersCollection.GetForRegion(region, startTime, endTime)
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
							offers, err := offersCollection.GetForRegion(region, startTime, endTime)
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
							offers, err := offersCollection.GetForRegion(region, startTime, endTime)
							Expect(err).NotTo(HaveOccurred())
							Expect(offers).To(HaveLen(1))
							Expect(offers).To(ContainOfferMock(2))
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
							offers, err := offersCollection.GetForRegion(region, startTime, endTime)
							Expect(err).NotTo(HaveOccurred())
							Expect(offers).To(HaveLen(1))
							Expect(offers).To(ContainOfferMock(2))
						})
					})

					Context("with date just after fromTime-24h", func() {
						BeforeEach(func() {
							startTime = time.Date(2014, 11, 12, 11, 01, 0, 0, time.UTC)
						})

						It("should not get any offers", func(done Done) {
							defer close(done)
							offers, err := offersCollection.GetForRegion(region, startTime, endTime)
							Expect(err).NotTo(HaveOccurred())
							Expect(offers).To(BeEmpty())
						})
					})
				})
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
