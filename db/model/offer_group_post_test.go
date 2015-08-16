package model_test

import (
	"time"

	"github.com/Lunchr/luncher-api/db/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("OfferGroupPost", func() {
	Describe("DateWithoutTime", func() {
		It("returns the date in the correct layout", func() {
			t := time.Date(2011, time.January, 2, 0, 0, 0, 1, time.UTC)
			result := model.DateWithoutTime(t)
			Expect(string(result)).To(Equal("2011-01-02"))
		})

		It("plays well with other timezones", func() {
			location, err := time.LoadLocation("Europe/Tallinn")
			Expect(err).NotTo(HaveOccurred())
			t := time.Date(2011, time.January, 2, 0, 0, 0, 1, location)
			result := model.DateWithoutTime(t)
			Expect(string(result)).To(Equal("2011-01-02"))
		})
	})
})
