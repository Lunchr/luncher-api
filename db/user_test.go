package db_test

import (
	"github.com/deiwin/luncher-api/db/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"
)

const (
	facebookUserID = "1387231118239727"
	facebookPageID = "1556442521239635"
)

var _ = Describe("User", func() {
	Describe("GetFbID", func() {
		It("should get by facebook user id", func(done Done) {
			defer close(done)
			user, err := usersCollection.GetFbID(facebookUserID)
			Expect(err).NotTo(HaveOccurred())
			Expect(user).NotTo(BeNil())
			Expect(user.FacebookPageID).To(Equal(facebookPageID))
		})

		It("should get nothing for wrong facebook id", func(done Done) {
			defer close(done)
			_, err := usersCollection.GetFbID(facebookPageID)
			Expect(err).To(HaveOccurred())
		})

		It("should link to the restaurant", func(done Done) {
			defer close(done)
			user, err := usersCollection.GetFbID(facebookUserID)
			Expect(err).NotTo(HaveOccurred())
			Expect(user).NotTo(BeNil())
			restaurant, err := restaurantsCollection.GetID(user.RestaurantID)
			Expect(err).NotTo(HaveOccurred())
			Expect(restaurant).NotTo(BeNil())
			Expect(restaurant.Name).To(Equal("Asian Chef"))
		})
	})

	Describe("GetAll", func() {
		It("should list all the users", func(done Done) {
			defer close(done)
			iter := usersCollection.GetAll()
			count := 0
			fbIDs := map[string]int{
				facebookUserID: 0,
				"another user": 0,
			}
			var user model.User
			for iter.Next(&user) {
				Expect(user).NotTo(BeNil())
				i, contains := fbIDs[user.FacebookUserID]
				Expect(contains).To(BeTrue())
				Expect(i).To(Equal(0))
				fbIDs[user.FacebookUserID]++
				count++
			}
			Expect(count).To(Equal(2))
			err := iter.Close()
			Expect(err).NotTo(HaveOccurred())
		})
	})

	Describe("updates", func() {
		RebuildDBAfterEach()
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
					user, err := usersCollection.GetFbID(facebookUserID)
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
					user, err := usersCollection.GetFbID(facebookUserID)
					Expect(err).NotTo(HaveOccurred())
					Expect(user).NotTo(BeNil())
					Expect(user.Session.FacebookPageToken).To(Equal("bsd"))
				})
			})
		})

		Describe("Update", func() {
			Context("with user updated with a page id change", func() {
				BeforeEach(func(done Done) {
					defer close(done)
					updatedUser := mocks.users[0]
					updatedUser.FacebookPageID = "bsd"
					err := usersCollection.Update(facebookUserID, updatedUser)
					Expect(err).NotTo(HaveOccurred())
				})

				It("should be reflected in the Get", func(done Done) {
					defer close(done)
					user, err := usersCollection.GetFbID(facebookUserID)
					Expect(err).NotTo(HaveOccurred())
					Expect(user).NotTo(BeNil())
					Expect(user.FacebookPageID).To(Equal("bsd"))
				})
			})
		})

		Describe("SessionID", func() {
			Context("with SessionID set", func() {
				var id string

				BeforeEach(func(done Done) {
					defer close(done)
					id = "someid"
					err := usersCollection.SetSessionID(mocks.userID, id)
					Expect(err).NotTo(HaveOccurred())
				})

				Describe("SetSessionID", func() {
					It("should be included in the Get", func(done Done) {
						defer close(done)
						user, err := usersCollection.GetFbID(facebookUserID)
						Expect(err).NotTo(HaveOccurred())
						Expect(user).NotTo(BeNil())
						Expect(user.Session.ID).To(Equal("someid"))
					})
				})

				Describe("GetBySessionID", func() {
					It("should be included in the Get", func(done Done) {
						defer close(done)
						user, err := usersCollection.GetSessionID(id)
						Expect(err).NotTo(HaveOccurred())
						Expect(user).NotTo(BeNil())
						Expect(user.FacebookUserID).To(Equal(facebookUserID))
					})
				})
			})
		})
	})
})
