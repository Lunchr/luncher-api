package db_test

import (
	"time"

	"github.com/Lunchr/luncher-api/db/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2/bson"
)

var _ = Describe("OfferGroupPosts", func() {
	var aPost = func() *model.OfferGroupPost {
		return &model.OfferGroupPost{
			RestaurantID: bson.NewObjectId(),
		}
	}

	Describe("Insert", func() {
		RebuildDBAfterEach()

		It("should return the posts with new IDs", func() {
			posts, err := offerGroupPostsCollection.Insert(aPost(), aPost())
			Expect(err).NotTo(HaveOccurred())
			Expect(posts).To(HaveLen(2))
			Expect(posts[0].ID).NotTo(BeEmpty())
			Expect(posts[1].ID).NotTo(BeEmpty())
		})

		It("should keep current ID if exists", func() {
			id := bson.NewObjectId()
			post := aPost()
			post.ID = id
			posts, err := offerGroupPostsCollection.Insert(post, aPost())
			Expect(err).NotTo(HaveOccurred())
			Expect(posts).To(HaveLen(2))
			Expect(posts[0].ID).To(Equal(id))
			Expect(posts[1].ID).NotTo(Equal(id))
		})
	})

	Describe("UpdateByID", func() {
		RebuildDBAfterEach()
		It("should fail for a non-existent ID", func() {
			err := offerGroupPostsCollection.UpdateByID(bson.NewObjectId(), aPost())
			Expect(err).To(HaveOccurred())
		})

		Context("with an post with known ID inserted", func() {
			var id bson.ObjectId
			BeforeEach(func() {
				id = bson.NewObjectId()
				post := aPost()
				post.ID = id
				_, err := offerGroupPostsCollection.Insert(post)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should update the post in DB", func() {
				post := aPost()
				post.MessageTemplate = "an updated message"
				err := offerGroupPostsCollection.UpdateByID(id, post)
				Expect(err).NotTo(HaveOccurred())
				result, err := offerGroupPostsCollection.GetByID(id)
				Expect(err).NotTo(HaveOccurred())
				Expect(result.MessageTemplate).To(Equal("an updated message"))
			})
		})
	})

	Describe("GetByDate", func() {
		var date = model.DateFromTime(time.Date(2115, 04, 03, 0, 0, 0, 0, time.UTC))
		var restaurantID = bson.NewObjectId()
		RebuildDBAfterEach()

		It("fails if not found", func() {
			_, err := offerGroupPostsCollection.GetByDate(date, restaurantID)
			Expect(err).To(HaveOccurred())
		})

		Context("with an post with known ID inserted", func() {
			BeforeEach(func() {
				post := aPost()
				post.Date = date
				post.RestaurantID = restaurantID
				_, err := offerGroupPostsCollection.Insert(post)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns the found post", func() {
				result, err := offerGroupPostsCollection.GetByDate(date, restaurantID)
				Expect(err).NotTo(HaveOccurred())
				Expect(result).NotTo(BeNil())
			})
		})
	})
})
