package interact_test

import (
	"errors"

	"github.com/deiwin/luncher-api/lunchman/interact"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Input", func() {
	var message = "Please answer"
	Describe("GetInput", func() {
		Context("with user input", func() {
			BeforeEach(func() {
				userInput = " user-input \n"
			})

			It("should return the trimmed input", func() {
				input, err := actor.GetInput(message)
				Expect(err).NotTo(HaveOccurred())
				Expect(input).To(Equal("user-input"))
				Eventually(output).Should(gbytes.Say(`Please answer: `))
			})

			Context("with a check", func() {
				var (
					checkErr error
					check    interact.InputCheck
				)

				JustBeforeEach(func() {
					check = func(input string) error {
						Expect(input).To(Equal("user-input"))
						return checkErr
					}
				})
				Context("with a failing check", func() {
					BeforeEach(func() {
						checkErr = errors.New("Check failed!")
					})

					It("should return the error from the check", func() {
						_, err := actor.GetInput(message, check)
						Expect(err).To(Equal(checkErr))
					})

					Context("with another check after the failed one", func() {
						It("should not call the second check", func() {
							actor.GetInput(message, check, func(input string) error {
								Fail("should not be called")
								return nil
							})
						})
					})
				})

				Context("with a passing check", func() {
					BeforeEach(func() {
						checkErr = nil
					})

					It("should not return an error", func() {
						_, err := actor.GetInput(message, check)
						Expect(err).NotTo(HaveOccurred())
					})
				})
			})
		})
	})

})
