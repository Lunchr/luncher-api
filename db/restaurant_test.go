package db_test

import (
	"github.com/Lunchr/luncher-api/db/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2/bson"
)

var _ = Describe("Restaurant", func() {
	Describe("Insert", func() {
		RebuildDBAfterEach()
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

	Describe("GetAll", func() {
		It("should list all the restaurants", func(done Done) {
			defer close(done)
			iter := restaurantsCollection.GetAll()
			count := 0
			restaurantNames := map[string]int{
				"Asian Chef":        0,
				"Bulgarian Dude":    0,
				"Caesarian Kitchen": 0,
			}
			var restaurant model.Restaurant
			for iter.Next(&restaurant) {
				Expect(restaurant).NotTo(BeNil())
				i, contains := restaurantNames[restaurant.Name]
				Expect(contains).To(BeTrue())
				Expect(i).To(Equal(0))
				restaurantNames[restaurant.Name]++
				count++
			}
			Expect(count).To(Equal(3))
			err := iter.Close()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("GetID", func() {
		It("should get by id", func(done Done) {
			defer close(done)
			res, err := restaurantsCollection.GetID(mocks.restaurantID)
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

	Describe("UpdateID", func() {
		RebuildDBAfterEach()
		It("should fail for a non-existent ID", func(done Done) {
			defer close(done)
			err := restaurantsCollection.UpdateID(bson.NewObjectId(), &model.Restaurant{})
			Expect(err).To(HaveOccurred())
		})

		Context("with a restaurant with known ID inserted", func() {
			var id bson.ObjectId
			BeforeEach(func(done Done) {
				defer close(done)
				id = bson.NewObjectId()
				_, err := restaurantsCollection.Insert(&model.Restaurant{
					ID: id,
				})
				Expect(err).NotTo(HaveOccurred())
			})

			It("should update the restaurant in DB", func(done Done) {
				defer close(done)
				err := restaurantsCollection.UpdateID(id, &model.Restaurant{
					Name: "an updated name",
				})
				Expect(err).NotTo(HaveOccurred())
				restaurant, err := restaurantsCollection.GetID(id)
				Expect(err).NotTo(HaveOccurred())
				Expect(restaurant.Name).To(Equal("an updated name"))
			})
		})
	})
})
