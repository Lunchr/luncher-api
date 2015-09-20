package handler_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/Lunchr/luncher-api/handler"
	"github.com/Lunchr/luncher-api/router"
	"github.com/Lunchr/luncher-api/session"
	"github.com/deiwin/facebook"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	testURL = "http://domain.extension/a/valid/url"
)

var _ = Describe("FacebookHandler", func() {
	var (
		auther         facebook.Authenticator
		sessionManager session.Manager
		handler        router.Handler
	)

	BeforeEach(func() {
		auther = &mockAuthenticator{}
		sessionManager = &mockSessionManager{}
	})

	JustBeforeEach(func() {
		handler = RedirectToFBForLogin(sessionManager, auther)
	})

	Describe("Login", func() {
		It("should redirect", func() {
			handler(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusSeeOther))
		})

		It("should redirect to mocked URL", func() {
			handler(responseRecorder, request)
			ExpectLocationToBeMockedURL(responseRecorder, testURL)
		})
	})
})

func ExpectLocationToBeMockedURL(responseRecorder *httptest.ResponseRecorder, url string) {
	location := responseRecorder.HeaderMap["Location"]
	Expect(location).To(HaveLen(1))
	Expect(location[0]).To(Equal(url))
}

type mockSessionManager struct {
	isSet bool
	id    string
}

func (m mockSessionManager) Get(r *http.Request) (string, error) {
	if !m.isSet {
		return "", session.ErrNotFound
	}
	if m.id == "" {
		return "session", nil
	}
	return m.id, nil
}

func (m mockSessionManager) GetOrInit(w http.ResponseWriter, r *http.Request) string {
	return "session"
}

type mockAuthenticator struct {
	api facebook.API
	facebook.Authenticator
}

func (m mockAuthenticator) AuthURL(session string) string {
	Expect(session).To(Equal("session"))
	return testURL
}
