package model_test

import (
	"github.com/Lunchr/luncher-api/db/model"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RegistrationAccessToken", func() {
	Describe("Token", func() {
		Describe("NewToken", func() {
			It("doesn't return duplicate items", func() {
				t1, err := model.NewToken()
				Expect(err).NotTo(HaveOccurred())
				t2, err := model.NewToken()
				Expect(err).NotTo(HaveOccurred())
				Expect(t1).NotTo(Equal(t2))
			})
		})

		Describe("TokenFromString", func() {
			It("creates a correct token from a known string and back", func() {
				t, err := model.TokenFromString("EF4120DA-0302-BCEE-712B-1C258D2FB6D4")
				Expect(err).NotTo(HaveOccurred())
				Expect(t.String()).To(Equal("EF4120DA-0302-BCEE-712B-1C258D2FB6D4"))
			})
		})
	})
})
