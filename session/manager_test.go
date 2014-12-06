package session_test

import (
	"net/http"
	"net/http/httptest"
	"strings"

	. "github.com/deiwin/luncher-api/session"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	sessionCookieName = "luncher_session"
	cookieHeader      = "Set-Cookie"
)

var _ = Describe("Manager", func() {
	Describe("GetOrInitSession", func() {
		var manager Manager

		BeforeEach(func() {
			manager = NewManager()
		})

		It("shoud return a non-empty string", func() {
			cookieValue := manager.GetOrInitSession(responseRecorder, request)
			Expect(cookieValue).NotTo(HaveLen(0))
			verifySingleSessionCookie(responseRecorder, func(cookieValue string) {
				Expect(cookieValue).NotTo(HaveLen(0))
			})
		})

		Context("with session cookie in request", func() {
			var requestCookieValue = "k_bV590l1T7mkhmwQgAIDA=="

			BeforeEach(func() {
				cookie := &http.Cookie{
					Name:  sessionCookieName,
					Value: requestCookieValue,
				}
				request.AddCookie(cookie)
			})

			It("should return the same cookie", func() {
				cookieValue := manager.GetOrInitSession(responseRecorder, request)
				Expect(cookieValue).To(Equal(requestCookieValue))
			})
		})
	})
})

func verifySingleSessionCookie(responseRecorder *httptest.ResponseRecorder, verify func(string)) {
	cookies := responseRecorder.HeaderMap[cookieHeader]
	Expect(cookies).To(HaveLen(1))
	cookie := cookies[0]
	Expect(cookie).To(HavePrefix(sessionCookieName + "="))
	cookieValue := strings.TrimPrefix(cookie, sessionCookieName+"=")
	verify(cookieValue)
}
