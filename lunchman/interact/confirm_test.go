package interact_test

import (
	"github.com/deiwin/luncher-api/lunchman/interact"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Confirm", func() {
	var (
		def     interact.ConfirmDefault
		message = "Are you sure?"
	)

	Context("with no default", func() {
		BeforeEach(func() {
			def = interact.ConfirmNoDefault
		})

		It("should ask with yes displayed as default", func() {
			actor.Confirm(message, def)
			Eventually(output).Should(gbytes.Say(`Are you sure\? \[y/n\]: `))
		})

		Context("with user answering yes", func() {
			BeforeEach(func() {
				userInput = "y\n"
			})

			It("should return true", func() {
				confirmed, err := actor.Confirm(message, def)
				Expect(err).NotTo(HaveOccurred())
				Expect(confirmed).To(BeTrue())
			})
		})

		Context("with user answering no", func() {
			BeforeEach(func() {
				userInput = "n\n"
			})

			It("should return false", func() {
				confirmed, err := actor.Confirm(message, def)
				Expect(err).NotTo(HaveOccurred())
				Expect(confirmed).To(BeFalse())
			})
		})

		Context("with user answering nothing", func() {
			BeforeEach(func() {
				userInput = "\n"
			})

			It("should return an error", func() {
				_, err := actor.Confirm(message, def)
				Expect(err).To(Equal(interact.ErrNoOptionSelected))
			})
		})

		Context("with user answering gibberish", func() {
			BeforeEach(func() {
				userInput = "asdfsadfa\n"
			})

			It("should return an error", func() {
				_, err := actor.Confirm(message, def)
				Expect(err).To(Equal(interact.ErrNoOptionSelected))
			})
		})
	})

	Context("with no as default", func() {
		BeforeEach(func() {
			def = interact.ConfirmDefaultToNo
		})

		It("should ask with yes displayed as default", func() {
			actor.Confirm(message, def)
			Eventually(output).Should(gbytes.Say(`Are you sure\? \[y/N\]: `))
		})

		Context("with user answering yes", func() {
			BeforeEach(func() {
				userInput = "y\n"
			})

			It("should return true", func() {
				confirmed, err := actor.Confirm(message, def)
				Expect(err).NotTo(HaveOccurred())
				Expect(confirmed).To(BeTrue())
			})
		})

		Context("with user answering no", func() {
			BeforeEach(func() {
				userInput = "n\n"
			})

			It("should return false", func() {
				confirmed, err := actor.Confirm(message, def)
				Expect(err).NotTo(HaveOccurred())
				Expect(confirmed).To(BeFalse())
			})
		})

		Context("with user answering nothing", func() {
			BeforeEach(func() {
				userInput = "\n"
			})

			It("should return false", func() {
				confirmed, err := actor.Confirm(message, def)
				Expect(err).NotTo(HaveOccurred())
				Expect(confirmed).To(BeFalse())
			})
		})

		Context("with user answering gibberish", func() {
			BeforeEach(func() {
				userInput = "asdfasdf\n"
			})

			It("should return an error", func() {
				_, err := actor.Confirm(message, def)
				Expect(err).To(Equal(interact.ErrNoOptionSelected))
			})
		})
	})

	Context("with yes as default", func() {
		BeforeEach(func() {
			def = interact.ConfirmDefaultToYes
		})

		It("should ask with yes displayed as default", func() {
			actor.Confirm(message, def)
			Eventually(output).Should(gbytes.Say(`Are you sure\? \[Y/n\]: `))
		})

		Context("with user answering yes", func() {
			BeforeEach(func() {
				userInput = "y\n"
			})

			It("should return true", func() {
				confirmed, err := actor.Confirm(message, def)
				Expect(err).NotTo(HaveOccurred())
				Expect(confirmed).To(BeTrue())
			})
		})

		Context("with user answering no", func() {
			BeforeEach(func() {
				userInput = "n\n"
			})

			It("should return false", func() {
				confirmed, err := actor.Confirm(message, def)
				Expect(err).NotTo(HaveOccurred())
				Expect(confirmed).To(BeFalse())
			})
		})

		Context("with user answering nothing", func() {
			BeforeEach(func() {
				userInput = "\n"
			})

			It("should return true", func() {
				confirmed, err := actor.Confirm(message, def)
				Expect(err).NotTo(HaveOccurred())
				Expect(confirmed).To(BeTrue())
			})
		})

		Context("with user answering gibberish", func() {
			BeforeEach(func() {
				userInput = "sadfasdf\n"
			})

			It("should return an error", func() {
				_, err := actor.Confirm(message, def)
				Expect(err).To(Equal(interact.ErrNoOptionSelected))
			})
		})
	})
})
