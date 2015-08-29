package db_test

import (
	"time"

	"github.com/Lunchr/luncher-api/db/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2/bson"
)

var _ = Describe("RegistrationAccessTokens", func() {
	var aToken = func() *model.RegistrationAccessToken {
		return &model.RegistrationAccessToken{
			CreatedAt: time.Now(),
		}
	}

	Describe("Insert", func() {
		RebuildDBAfterEach()
		It("returns the token with new IDs", func() {
			token, err := registrationAccessTokensCollection.Insert(aToken())
			Expect(err).NotTo(HaveOccurred())
			Expect(token.ID).NotTo(BeEmpty())
		})

		It("keeps current ID if exists", func() {
			id := bson.NewObjectId()
			token := aToken()
			token.ID = id
			insertedToken, err := registrationAccessTokensCollection.Insert(token)
			Expect(err).NotTo(HaveOccurred())
			Expect(insertedToken.ID).To(Equal(id))
		})
	})

	Describe("Exists", func() {
		RebuildDBAfterEach()
		It("returns false if no such token in the DB", func() {
			token, err := model.NewToken()
			Expect(err).NotTo(HaveOccurred())
			exists, err := registrationAccessTokensCollection.Exists(token)
			Expect(err).NotTo(HaveOccurred())
			Expect(exists).To(BeFalse())
		})

		Context("with an known token inserted", func() {
			var token model.Token
			BeforeEach(func() {
				var err error
				token, err = model.NewToken()
				Expect(err).NotTo(HaveOccurred())
				regToken := aToken()
				regToken.Token = token
				_, err = registrationAccessTokensCollection.Insert(regToken)
				Expect(err).NotTo(HaveOccurred())
			})

			It("returns true for that token", func() {
				exists, err := registrationAccessTokensCollection.Exists(token)
				Expect(err).NotTo(HaveOccurred())
				Expect(exists).To(BeTrue())
			})
		})
	})
})
