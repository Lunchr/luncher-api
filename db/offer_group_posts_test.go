package db_test

import (
	"github.com/Lunchr/luncher-api/db/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2/bson"
)

var _ = Describe("OfferGroupPosts", func() {

	Describe("Insert", func() {
		RebuildDBAfterEach()

		var aPost = func() *model.OfferGroupPost {
			return &model.OfferGroupPost{
				RestaurantID: bson.NewObjectId(),
			}
		}

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
})
