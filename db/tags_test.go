package db_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Tags", func() {
	Describe("Get", func() {
		It("should get all tags", func(done Done) {
			defer close(done)
			tags, err := tagsCollection.Get()
			Expect(err).NotTo(HaveOccurred())
			Expect(tags).To(HaveLen(5))
		})
	})
})
