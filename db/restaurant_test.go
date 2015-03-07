package db_test

import (
	"github.com/deiwin/luncher-api/db/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2/bson"
)

var _ = Describe("Restaurant", func() {
	Describe("Insert", func() {
		// Return the DB to the original state after these tests
		AfterEach(func(done Done) {
			defer close(done)
			wipeDb()
			initCollections()
		})

		It("should return the restaurants with new IDs", func(done Done) {
			defer close(done)
			restaurants, err := restaurantsCollection.Insert(&model.Restaurant{}, &model.Restaurant{})
			Expect(err).NotTo(HaveOccurred())
			Expect(restaurants).To(HaveLen(2))
			Expect(restaurants[0].ID).NotTo(BeEmpty())
			Expect(restaurants[1].ID).NotTo(BeEmpty())
		})

		It("should keep current ID if exists", func(done Done) {
			defer close(done)
			id := bson.NewObjectId()
			restaurants, err := restaurantsCollection.Insert(&model.Restaurant{
				ID: id,
			}, &model.Restaurant{})
			Expect(err).NotTo(HaveOccurred())
			Expect(restaurants).To(HaveLen(2))
			Expect(restaurants[0].ID).To(Equal(id))
			Expect(restaurants[1].ID).NotTo(Equal(id))
		})
	})

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
