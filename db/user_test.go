package db_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"
)

const (
	facebookUserID = "1387231118239727"
	facebookPageID = "1556442521239635"
)

var _ = Describe("User", func() {
	Describe("Get", func() {
		It("should get by facebook user id", func(done Done) {
			defer close(done)
			user, err := usersCollection.Get(facebookUserID)
			Expect(err).NotTo(HaveOccurred())
			Expect(user).NotTo(BeNil())
			Expect(user.FacebookPageID).To(Equal(facebookPageID))
		})

		It("should get nothing for wrong facebook id", func(done Done) {
			defer close(done)
			user, err := usersCollection.Get(facebookPageID)
			Expect(err).To(HaveOccurred())
			Expect(user).To(BeNil())
		})

		It("should link to the restaurant", func(done Done) {
			defer close(done)
			user, err := usersCollection.Get(facebookUserID)
			Expect(err).NotTo(HaveOccurred())
			Expect(user).NotTo(BeNil())
			// there's no get by ID method at the moment so just get all and see
			restaurant, err := restaurantsCollection.GetByID(user.RestaurantID)
			Expect(err).NotTo(HaveOccurred())
			Expect(restaurant).NotTo(BeNil())
			Expect(restaurant.Name).To(Equal("Asian Chef"))
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
				err := usersCollection.SetAccessToken(facebookUserID, token)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should be included in the Get", func(done Done) {
				defer close(done)
				user, err := usersCollection.Get(facebookUserID)
				Expect(err).NotTo(HaveOccurred())
				Expect(user).NotTo(BeNil())
				Expect(user.Session.FacebookUserToken.AccessToken).To(Equal("asd"))
			})
		})
	})

	Describe("SetPageAccessToken", func() {
		Context("with access token set", func() {
			var token string

			BeforeEach(func(done Done) {
				defer close(done)
				token = "bsd"
				err := usersCollection.SetPageAccessToken(facebookUserID, token)
				Expect(err).NotTo(HaveOccurred())
			})

			It("should be included in the Get", func(done Done) {
				defer close(done)
				user, err := usersCollection.Get(facebookUserID)
				Expect(err).NotTo(HaveOccurred())
				Expect(user).NotTo(BeNil())
				Expect(user.Session.FacebookPageToken).To(Equal("bsd"))
			})
		})
	})

	Describe("SessionID", func() {
		Context("with SessionID set", func() {
			var id string

			BeforeEach(func(done Done) {
				defer close(done)
				id = "someid"
				err := usersCollection.SetSessionID(facebookUserID, id)
				Expect(err).NotTo(HaveOccurred())
			})

			Describe("SetSessionID", func() {
				It("should be included in the Get", func(done Done) {
					defer close(done)
					user, err := usersCollection.Get(facebookUserID)
					Expect(err).NotTo(HaveOccurred())
					Expect(user).NotTo(BeNil())
					Expect(user.Session.ID).To(Equal("someid"))
				})
			})

			Describe("GetBySessionID", func() {
				It("should be included in the Get", func(done Done) {
					defer close(done)
					user, err := usersCollection.GetBySessionID(id)
					Expect(err).NotTo(HaveOccurred())
					Expect(user).NotTo(BeNil())
					Expect(user.FacebookUserID).To(Equal(facebookUserID))
				})
			})
		})
	})
})
