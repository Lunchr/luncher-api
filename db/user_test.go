package db_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"
)

var _ = Describe("User", func() {
	Describe("Get", func() {
		It("should get by facebook user id", func(done Done) {
			defer close(done)
			user, err := usersCollection.Get("1")
			Expect(err).NotTo(HaveOccurred())
			Expect(user).NotTo(BeNil())
			Expect(user.FacebookPageID).To(Equal("2"))
		})

		It("should get nothing for wrong facebook id", func(done Done) {
			defer close(done)
			user, err := usersCollection.Get("2")
			Expect(err).To(HaveOccurred())
			Expect(user).To(BeNil())
		})

		It("should link to the restaurant", func(done Done) {
			defer close(done)
			user, err := usersCollection.Get("1")
			Expect(err).NotTo(HaveOccurred())
			Expect(user).NotTo(BeNil())
			// there's no get by ID method at the moment so just get all and see
			restaurants, err := restaurantsCollection.Get()
			Expect(err).NotTo(HaveOccurred())
			found := false
			for _, restaurant := range restaurants {
				if restaurant.ID == user.RestaurantID {
					found = true
					Expect(restaurant.Name).To(Equal("Asian Chef"))
				}
			}
			Expect(found).To(BeTrue())
		})
	})

	Describe("SetAccessToken", func() {
		Context("with access token set", func() {
			var token oauth2.Token

			BeforeEach(func(done Done) {
				defer close(done)
				token = oauth2.Token{
					AccessToken: "asd",
				}
				err := usersCollection.SetAccessToken("1", token)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should be included in the Get", func(done Done) {
				defer close(done)
				user, err := usersCollection.Get("1")
				Expect(err).NotTo(HaveOccurred())
				Expect(user).NotTo(BeNil())
				Expect(user.FacebookUserToken.AccessToken).To(Equal("asd"))
			})
		})
	})

	Describe("SetPageAccessToken", func() {
		Context("with access token set", func() {
			var token string

			BeforeEach(func(done Done) {
				defer close(done)
				token = "bsd"
				err := usersCollection.SetPageAccessToken("1", token)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should be included in the Get", func(done Done) {
				defer close(done)
				user, err := usersCollection.Get("1")
				Expect(err).NotTo(HaveOccurred())
				Expect(user).NotTo(BeNil())
				Expect(user.FacebookPageToken).To(Equal("bsd"))
			})
		})
	})
})
