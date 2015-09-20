package db_test

import (
	"github.com/Lunchr/luncher-api/db/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"golang.org/x/oauth2"
	"gopkg.in/mgo.v2/bson"
)

const (
	facebookUserID = "1387231118239727"
	facebookPageID = "1556442521239635"
)

var _ = Describe("User", func() {
	Describe("GetFbID", func() {
		It("should get by facebook user id", func() {
			user, err := usersCollection.GetFbID(facebookUserID)
			Expect(err).NotTo(HaveOccurred())
			Expect(user).NotTo(BeNil())
		})

		It("should get nothing for wrong facebook id", func() {
			_, err := usersCollection.GetFbID(facebookPageID)
			Expect(err).To(HaveOccurred())
		})

		It("should link to the restaurant", func() {
			user, err := usersCollection.GetFbID(facebookUserID)
			Expect(err).NotTo(HaveOccurred())
			Expect(user).NotTo(BeNil())
			restaurant, err := restaurantsCollection.GetID(user.RestaurantIDs[0])
			Expect(err).NotTo(HaveOccurred())
			Expect(restaurant).NotTo(BeNil())
			Expect(restaurant.Name).To(Equal("Asian Chef"))
		})
	})

	Describe("GetAll", func() {
		It("should list all the users", func() {
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

				BeforeEach(func() {
					token = oauth2.Token{
						AccessToken: "asd",
					}
					err := usersCollection.SetAccessToken(facebookUserID, token)
					Expect(err).NotTo(HaveOccurred())
				})

				It("should be included in the Get", func() {
					user, err := usersCollection.GetFbID(facebookUserID)
					Expect(err).NotTo(HaveOccurred())
					Expect(user).NotTo(BeNil())
					Expect(user.Session.FacebookUserToken.AccessToken).To(Equal("asd"))
				})
			})
		})

		Describe("SetPageAccessToken", func() {
			Context("with access token set", func() {
				var tokens []model.FacebookPageToken

				BeforeEach(func() {
					tokens = []model.FacebookPageToken{model.FacebookPageToken{
						PageID: "pageid",
						Token:  "atoken",
					}}
					err := usersCollection.SetPageAccessTokens(facebookUserID, tokens)
					Expect(err).NotTo(HaveOccurred())
				})

				It("should be included in the Get", func() {
					user, err := usersCollection.GetFbID(facebookUserID)
					Expect(err).NotTo(HaveOccurred())
					Expect(user).NotTo(BeNil())
					Expect(user.Session.FacebookPageTokens).To(Equal(tokens))
				})
			})
		})

		Describe("Update", func() {
			Context("with user updated with a facebook user id change", func() {
				var newID bson.ObjectId
				BeforeEach(func() {
					updatedUser := *mocks.users[0]
					newID = bson.NewObjectId()
					updatedUser.RestaurantIDs = []bson.ObjectId{newID}
					// this isn't strictly necessary, but otherwise this test seems to fail
					// on older MongoDB versions
					updatedUser.ID = ""
					err := usersCollection.Update(facebookUserID, &updatedUser)
					Expect(err).NotTo(HaveOccurred())
				})

				It("should be reflected in the Get", func() {
					user, err := usersCollection.GetFbID(facebookUserID)
					Expect(err).NotTo(HaveOccurred())
					Expect(user).NotTo(BeNil())
					Expect(user.RestaurantIDs[0]).To(Equal(newID))
				})
			})
		})

		Describe("SessionID", func() {
			Context("with SessionID set", func() {
				var id string

				BeforeEach(func() {
					id = "someid"
					err := usersCollection.SetSessionID(mocks.userID, id)
					Expect(err).NotTo(HaveOccurred())
				})

				Describe("SetSessionID", func() {
					It("should be included in the Get", func() {
						user, err := usersCollection.GetFbID(facebookUserID)
						Expect(err).NotTo(HaveOccurred())
						Expect(user).NotTo(BeNil())
						Expect(user.Session.ID).To(Equal("someid"))
					})
				})

				Describe("UnsetSessionID", func() {
					It("should remove the session ID", func() {
						err := usersCollection.UnsetSessionID(mocks.userID)
						Expect(err).NotTo(HaveOccurred())

						user, err := usersCollection.GetFbID(facebookUserID)
						Expect(err).NotTo(HaveOccurred())
						Expect(user).NotTo(BeNil())
						Expect(user.Session.ID).To(BeEmpty())
					})
				})

				Describe("GetBySessionID", func() {
					It("should be included in the Get", func() {
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
