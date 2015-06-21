package handler_test

import (
	"net/http"

	"github.com/Lunchr/luncher-api/db"
	. "github.com/Lunchr/luncher-api/handler"
	"github.com/Lunchr/luncher-api/session"
	"github.com/deiwin/facebook"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Facebook registration handler", func() {
	var (
		registrationAuther  facebook.Authenticator
		mockSessMgr         session.Manager
		mockUsersCollection db.Users
		handlers            Facebook
	)

	BeforeEach(func() {
		registrationAuther = &mockAuthenticator{}
		mockSessMgr = &mockSessionManager{}
	})

	JustBeforeEach(func() {
		handlers = NewFacebook(nil, registrationAuther, mockSessMgr, mockUsersCollection)
	})

	Describe("Login", func() {
		It("should redirect", func(done Done) {
			defer close(done)
			handlers.RedirectToFBForRegistration()(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusSeeOther))
		})

		It("should redirect to mocked URL", func(done Done) {
			defer close(done)
			handlers.RedirectToFBForRegistration()(responseRecorder, request)
			ExpectLocationToBeMockedURL(responseRecorder, testURL)
		})
	})
})
