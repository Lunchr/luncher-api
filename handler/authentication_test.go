package handler_test

import (
	"net/http"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/router"
	"github.com/deiwin/luncher-api/session"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func ExpectUserToBeLoggedIn(handler func() *router.HandlerError, setDependencies func(session.Manager, db.Users)) {
	Describe("it expects the user to be logged in", func() {
		Context("with no session set", func() {
			BeforeEach(func() {
				setDependencies(&mockSessionManager{}, nil)
			})

			It("should be forbidden", func(done Done) {
				defer close(done)
				err := handler()
				Expect(err.Code).To(Equal(http.StatusForbidden))
			})
		})

		Context("with session set, but no matching user in DB", func() {
			BeforeEach(func() {
				setDependencies(&mockSessionManager{isSet: true}, mockUsers{})
			})

			It("should be forbidden", func(done Done) {
				defer close(done)
				err := handler()
				Expect(err.Code).To(Equal(http.StatusForbidden))
			})
		})
	})
}
