package db_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2"
)

var _ = Describe("Regions", func() {
	Describe("Get", func() {
		It("should get region by name", func(done Done) {
			defer close(done)
			region, err := regionsCollection.Get("Tartu")
			Expect(err).NotTo(HaveOccurred())
			Expect(region.Name).To(Equal("Tartu"))
			Expect(region.Location).To(Equal("Europe/Tallinn"))
		})

		It("should get region by name", func(done Done) {
			defer close(done)
			region, err := regionsCollection.Get("London")
			Expect(err).NotTo(HaveOccurred())
			Expect(region.Name).To(Equal("London"))
			Expect(region.Location).To(Equal("Europe/London"))
		})

		It("should return nothing if doesn't exist", func(done Done) {
			defer close(done)
			_, err := regionsCollection.Get("blablabla")
			Expect(err).To(Equal(mgo.ErrNotFound))
		})
	})
})
