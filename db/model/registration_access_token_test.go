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
			It("creates a correct token from a known string", func() {
				t, err := model.TokenFromString("EF4120DA-0302-BCEE-712B-1C258D2FB6D4")
				Expect(err).NotTo(HaveOccurred())
				Expect(t).To(Equal(model.Token{0xef, 0x41, 0x20, 0xda, 0x3, 0x2, 0xbc, 0xee, 0x71, 0x2b,
					0x1c, 0x25, 0x8d, 0x2f, 0xb6, 0xd4}))
			})
		})

		Describe("String", func() {
			It("creates a correct string from known token", func() {
				t := model.Token{0xef, 0x41, 0x20, 0xda, 0x3, 0x2, 0xbc, 0xee, 0x71, 0x2b, 0x1c, 0x25, 0x8d,
					0x2f, 0xb6, 0xd4}
				Expect(t.String()).To(Equal("EF4120DA-0302-BCEE-712B-1C258D2FB6D4"))
			})
		})
	})
})
