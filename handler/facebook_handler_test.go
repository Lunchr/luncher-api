package handler_test

import (
	"net/http"
	"net/http/httptest"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/facebook"
	. "github.com/deiwin/luncher-api/handler"
	"github.com/deiwin/luncher-api/session"
	"golang.org/x/oauth2"

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
		auther = &authenticator{}
		mockSessMgr = &mockManager{}
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

type mockManager struct{}

func (mock mockManager) GetOrInitSession(w http.ResponseWriter, r *http.Request) string {
	return "session"
}

type authenticator struct{}

func (a authenticator) AuthURL(session string) string {
	Expect(session).To(Equal("session"))
	return testURL
}

func (a authenticator) Token(code string, r *http.Request) (tok *oauth2.Token, err error) {
	return
}

func (a authenticator) PageAccessToken(tok *oauth2.Token, pageID string) (string, error) {
	return "", nil
}

func (a authenticator) APIConnection(tok *oauth2.Token) facebook.Connection {
	return nil
}
