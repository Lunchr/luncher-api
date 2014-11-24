package db_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Offers", func() {

	Describe("GetForTimeRange", func() {
		var (
			startTime time.Time
			endTime   time.Time
		)

		JustBeforeEach(func() {
			endTime = startTime.AddDate(0, 0, 1)
		})

		Context("with simple day range", func() {
			BeforeEach(func() {
				startTime = time.Date(2014, 11, 10, 0, 0, 0, 0, time.UTC)
			})

			It("should get offers for the day", func(done Done) {
				defer close(done)
				offers, err := offersCollection.GetForTimeRange(startTime, endTime)
				Expect(err).NotTo(HaveOccurred())
				Expect(offers).To(HaveLen(2))
				Expect(offers).To(ContainOfferMock(0))
				Expect(offers).To(ContainOfferMock(1))
				Expect(offers).NotTo(ContainOfferMock(2))
			})
		})

		Context("with simple day range withour offers", func() {
			BeforeEach(func() {
				startTime = time.Date(2014, 11, 13, 0, 0, 0, 0, time.UTC)
			})

			It("should get 0 offers for non-existent date", func(done Done) {
				defer close(done)
				offers, err := offersCollection.GetForTimeRange(startTime, endTime)
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
						offers, err := offersCollection.GetForTimeRange(startTime, endTime)
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
						offers, err := offersCollection.GetForTimeRange(startTime, endTime)
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
						offers, err := offersCollection.GetForTimeRange(startTime, endTime)
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
						offers, err := offersCollection.GetForTimeRange(startTime, endTime)
						Expect(err).NotTo(HaveOccurred())
						Expect(offers).To(BeEmpty())
					})
				})
			})

		})
	})
})
