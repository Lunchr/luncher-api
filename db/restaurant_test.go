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

	Describe("GetByID", func() {
		It("should get by id", func(done Done) {
			defer close(done)
			res, err := restaurantsCollection.GetByID(mocks.restaurantID)
			Expect(err).NotTo(HaveOccurred())
			Expect(res).NotTo(BeNil())
			Expect(res.ID).To(Equal(mocks.restaurantID))
		})
	})

	Describe("Exists", func() {
		It("should return true for an existing restaurant", func(done Done) {
			defer close(done)
			exists, err := restaurantsCollection.Exists("Asian Chef")
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeTrue())
		})

		It("should return false for a nonexisting restaurant", func(done Done) {
			defer close(done)
			exists, err := restaurantsCollection.Exists("bla bla bla")
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
		})
	})
})
