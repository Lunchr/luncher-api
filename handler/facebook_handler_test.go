package handler_test

import (
	"errors"
	"net/http"
	"net/http/httptest"

	"github.com/deiwin/facebook"
	"github.com/deiwin/luncher-api/db"
	. "github.com/deiwin/luncher-api/handler"
	"github.com/deiwin/luncher-api/session"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

const (
	testURL = "http://domain.extension/a/valid/url"
)

var _ = Describe("FacebookHandler", func() {
	var (
		auther              facebook.Authenticator
		mockSessMgr         session.Manager
		mockUsersCollection db.Users
		handlers            Facebook
	)

	BeforeEach(func() {
		auther = &mockAuthenticator{}
		mockSessMgr = &mockSessionManager{}
	})

	JustBeforeEach(func() {
		handlers = NewFacebook(auther, mockSessMgr, mockUsersCollection)
	})

	Describe("Login", func() {
		It("should redirect", func(done Done) {
			defer close(done)
			handlers.Login().ServeHTTP(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusSeeOther))
		})

		It("should redirect to mocked URL", func(done Done) {
			defer close(done)
			handlers.Login().ServeHTTP(responseRecorder, request)
			ExpectLocationToBeMockedURL(responseRecorder)
		})
	})
})

func ExpectLocationToBeMockedURL(responseRecorder *httptest.ResponseRecorder) {
	location := responseRecorder.HeaderMap["Location"]
	Expect(location).To(HaveLen(1))
	Expect(location[0]).To(Equal(testURL))
}

type mockSessionManager struct {
	isSet bool
	id    string
}

func (m mockSessionManager) Get(r *http.Request) (string, error) {
	if !m.isSet {
		return "", errors.New("no session")
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
