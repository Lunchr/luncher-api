package db_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Restaurant", func() {
	Describe("Get", func() {
		It("should get all restaurants", func(done Done) {
			defer close(done)
			tags, err := restaurantsCollection.Get()
			Expect(err).NotTo(HaveOccurred())
			Expect(tags).To(HaveLen(3))
		})
	})
})
