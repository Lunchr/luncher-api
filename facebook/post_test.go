package facebook_test

import (
	"github.com/Lunchr/luncher-api/db/model"
	"github.com/Lunchr/luncher-api/facebook"
	"github.com/Lunchr/luncher-api/facebook/mocks"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Post", func() {
	var (
		facebookPost facebook.Post

		groupPosts *mocks.OfferGroupPosts
		offers     *mocks.Offers
		regions    *mocks.Regions
		fbAuth     *mocks.Authenticator

		user       *model.User
		restaurant *model.Restaurant
	)

	BeforeEach(func() {
		groupPosts = new(mocks.OfferGroupPosts)
		offers = new(mocks.Offers)
		regions = new(mocks.Regions)
		fbAuth = new(mocks.Authenticator)

		facebookPost = facebook.NewPost(groupPosts, offers, regions, fbAuth)
	})

	Describe("Update", func() {
		var post *model.OfferGroupPost

		Context("for restaurants without an associated FB page", func() {
			BeforeEach(func() {
				restaurant = &model.Restaurant{
					FacebookPageID: "",
				}
			})

			It("does nothing", func() {
				err := facebookPost.Update(post, user, restaurant)
				Expect(err).NotTo(HaveOccurred())
			})
		})
	})
})
