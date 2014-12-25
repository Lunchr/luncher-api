package db_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Restaurant", func() {
	Describe("Get", func() {
		It("should get all restaurants", func(done Done) {
			defer close(done)
			restaurants, err := restaurantsCollection.Get()
			Expect(err).NotTo(HaveOccurred())
			Expect(restaurants).To(HaveLen(3))
		})
	})
})
