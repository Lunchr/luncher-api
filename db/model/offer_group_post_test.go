package model_test

import (
	"time"

	"github.com/Lunchr/luncher-api/db/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OfferGroupPost", func() {
	Describe("DateWithoutTime", func() {
		Describe("DateFromTime", func() {
			It("returns the date in the correct layout", func() {
				t := time.Date(2011, time.January, 2, 0, 0, 0, 1, time.UTC)
				result := model.DateFromTime(t)
				Expect(string(result)).To(Equal("2011-01-02"))
			})

			It("plays well with other timezones", func() {
				location, err := time.LoadLocation("Europe/Tallinn")
				Expect(err).NotTo(HaveOccurred())
				t := time.Date(2011, time.January, 2, 0, 0, 0, 1, location)
				result := model.DateFromTime(t)
				Expect(string(result)).To(Equal("2011-01-02"))
			})
		})

		Describe("IsValid", func() {
			It("returns true for a valid date", func() {
				date := model.DateWithoutTime("2015-11-18")
				Expect(date.IsValid()).To(BeTrue())
			})

			It("returns false for an invalid date", func() {
				date := model.DateWithoutTime("2015-18-11")
				Expect(date.IsValid()).To(BeFalse())
			})

			It("returns false for gibberish strings", func() {
				date := model.DateWithoutTime("asdfasdfasjksdlaf")
				Expect(date.IsValid()).To(BeFalse())
			})
		})

		Describe("TimeBounds", func() {
			It("returns correct bounds for valid data", func() {
				date := model.DateWithoutTime("2015-11-18")
				startTime, endTime, err := date.TimeBounds(time.UTC)
				Expect(err).NotTo(HaveOccurred())
				Expect(startTime).To(Equal(time.Date(2015, 11, 18, 0, 0, 0, 0, time.UTC)))
				Expect(endTime).To(Equal(time.Date(2015, 11, 19, 0, 0, 0, 0, time.UTC)))
			})

			It("behaves well in different timezones", func() {
				location, err := time.LoadLocation("America/New_York")
				Expect(err).NotTo(HaveOccurred())
				date := model.DateWithoutTime("2015-11-18")

				startTime, endTime, err := date.TimeBounds(location)

				Expect(err).NotTo(HaveOccurred())
				Expect(startTime).To(Equal(time.Date(2015, 11, 18, 0, 0, 0, 0, location)))
				Expect(endTime).To(Equal(time.Date(2015, 11, 19, 0, 0, 0, 0, location)))
			})

			It("fails for invalid dates", func() {
				date := model.DateWithoutTime("2015-71-18")
				_, _, err := date.TimeBounds(time.UTC)
				Expect(err).To(HaveOccurred())
			})
		})
	})
})
