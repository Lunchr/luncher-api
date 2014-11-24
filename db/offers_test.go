package db_test

import (
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Offers", func() {

	It("should get offers for date", func(done Done) {
		defer close(done)
		date := time.Date(2014, 11, 10, 0, 0, 0, 0, time.UTC)
		offers, err := offersCollection.GetForDate(date)
		Expect(err).NotTo(HaveOccurred())
		Expect(len(offers)).To(Equal(2))
	})
})
