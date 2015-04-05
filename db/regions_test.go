package db_test

import (
	"github.com/deiwin/luncher-api/db/model"
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

	Describe("GetAll", func() {
		It("should list all the regions", func(done Done) {
			defer close(done)
			iter := regionsCollection.GetAll()
			count := 0
			regionNames := map[string]int{
				"Tartu":   0,
				"Tallinn": 0,
				"London":  0,
			}
			var region model.Region
			for iter.Next(&region) {
				Expect(region).NotTo(BeNil())
				i, contains := regionNames[region.Name]
				Expect(contains).To(BeTrue())
				Expect(i).To(Equal(0))
				regionNames[region.Name]++
				count++
			}
			Expect(count).To(Equal(3))
			err := iter.Close()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
