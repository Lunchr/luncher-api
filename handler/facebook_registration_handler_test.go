package handler_test

import (
	"encoding/json"
	"net/http"

	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/db/model"
	. "github.com/Lunchr/luncher-api/handler"
	"github.com/Lunchr/luncher-api/handler/mocks"
	"github.com/Lunchr/luncher-api/router"
	"github.com/Lunchr/luncher-api/session"
	fbModel "github.com/deiwin/facebook/model"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Facebook registration handler", func() {
	var (
		auther          *mocks.Authenticator
		sessionManager  session.Manager
		usersCollection db.Users
		handler         router.Handler
	)

	BeforeEach(func() {
		auther = new(mocks.Authenticator)
		sessionManager = &mockSessionManager{}
	})

	Describe("Login", func() {
		var (
			accessTokens *mocks.RegistrationAccessTokens
			handler      router.HandlerWithParams
			params       httprouter.Params
		)

		BeforeEach(func() {
			accessTokens = new(mocks.RegistrationAccessTokens)
			auther.On("AuthURL", "session").Return(testURL)
		})

		JustBeforeEach(func() {
			handler = RedirectToFBForRegistration(sessionManager, auther, accessTokens)
		})

		Context("without an access token", func() {
			It("fails with StatusUnauthorized", func() {
				err := handler(responseRecorder, request, params)
				Expect(err.Code).To(Equal(http.StatusUnauthorized))
			})
		})

		Context("with a rubbish token", func() {
			BeforeEach(func() {
				params = httprouter.Params{httprouter.Param{
					Key:   "token",
					Value: "rubbish",
				}}
			})

			It("fails with StatusBadRequest", func() {
				err := handler(responseRecorder, request, params)
				Expect(err.Code).To(Equal(http.StatusBadRequest))
			})
		})

		Context("with a made up or expired token", func() {
			BeforeEach(func() {
				token, err := model.NewToken()
				Expect(err).NotTo(HaveOccurred())
				params = httprouter.Params{httprouter.Param{
					Key:   "token",
					Value: token.String(),
				}}
				accessTokens.On("Exists", token).Return(false, nil)
			})

			It("fails with StatusForbidden", func() {
				err := handler(responseRecorder, request, params)
				Expect(err.Code).To(Equal(http.StatusForbidden))
			})
		})

		Context("with a valid access token", func() {
			BeforeEach(func() {
				token, err := model.NewToken()
				Expect(err).NotTo(HaveOccurred())
				params = httprouter.Params{httprouter.Param{
					Key:   "token",
					Value: token.String(),
				}}
				accessTokens.On("Exists", token).Return(true, nil)
			})

			It("redirects", func() {
				err := handler(responseRecorder, request, params)
				Expect(err).NotTo(HaveOccurred())
				Expect(responseRecorder.Code).To(Equal(http.StatusSeeOther))
			})

			It("redirects to mocked URL", func() {
				handler(responseRecorder, request, params)
				ExpectLocationToBeMockedURL(responseRecorder, testURL)
			})
		})
	})

	Describe("ListPagesManagedByUser", func() {
		JustBeforeEach(func() {
			handler = ListPagesManagedByUser(sessionManager, auther, usersCollection)
		})

		ExpectUserToBeLoggedIn(func() *router.HandlerError {
			return handler(responseRecorder, request)
		}, func(mgr session.Manager, users db.Users) {
			sessionManager = mgr
			usersCollection = users
		})

		Context("with user logged in", func() {
			var (
				api *mocks.API
			)
			BeforeEach(func() {
				sessionManager = &mockSessionManager{isSet: true, id: "correctSession"}
				usersCollection = mockUsers{}
				api = new(mocks.API)
				auther.On("APIConnection", mock.AnythingOfType("*oauth2.Token")).Return(api)
				api.On("Accounts").Return(&fbModel.Accounts{
					Data: []fbModel.Page{
						fbModel.Page{
							ID:   "id1",
							Name: "name1",
						},
						fbModel.Page{
							ID:   "id2",
							Name: "name2",
						},
					},
				}, nil)
			})

			AfterEach(func() {
				auther.AssertExpectations(GinkgoT())
				api.AssertExpectations(GinkgoT())
			})

			It("should succeed", func() {
				err := handler(responseRecorder, request)
				Expect(err).To(BeNil())
			})

			It("should return json", func() {
				handler(responseRecorder, request)
				contentTypes := responseRecorder.HeaderMap["Content-Type"]
				Expect(contentTypes).To(HaveLen(1))
				Expect(contentTypes[0]).To(Equal("application/json"))
			})

			It("should respond with a list of pages returned from Facebook", func() {
				handler(responseRecorder, request)
				var result []*FacebookPage
				json.Unmarshal(responseRecorder.Body.Bytes(), &result)
				Expect(result).To(HaveLen(2))
				Expect(result[0].ID).To(Equal("id1"))
				Expect(result[0].Name).To(Equal("name1"))
				Expect(result[1].ID).To(Equal("id2"))
				Expect(result[1].Name).To(Equal("name2"))
			})
		})
	})

	Describe("GET Page", func() {
		var handler router.HandlerWithParams

		JustBeforeEach(func() {
			handler = Page(sessionManager, auther, usersCollection)
		})

		ExpectUserToBeLoggedIn(func() *router.HandlerError {
			return handler(responseRecorder, request, nil)
		}, func(mgr session.Manager, users db.Users) {
			sessionManager = mgr
			usersCollection = users
		})

		Context("with user logged in", func() {
			var (
				api    *mocks.API
				params httprouter.Params
			)

			BeforeEach(func() {
				sessionManager = &mockSessionManager{isSet: true, id: "correctSession"}
				usersCollection = mockUsers{}
				api = new(mocks.API)
				auther.On("APIConnection", mock.AnythingOfType("*oauth2.Token")).Return(api)
				api.On("Page", "a_page_id").Return(&fbModel.Page{
					ID:   "id1",
					Name: "name1",
					Location: fbModel.Location{
						Street:  "a_street 10",
						City:    "a_city",
						Country: "a_country",
					},
					Phone:   "a_phone",
					Website: "a_website",
					Emails:  []string{"an_email", "other_email"},
				}, nil)
				params = httprouter.Params{httprouter.Param{
					Key:   "id",
					Value: "a_page_id",
				}}
			})

			AfterEach(func() {
				auther.AssertExpectations(GinkgoT())
				api.AssertExpectations(GinkgoT())
			})

			It("should succeed", func() {
				err := handler(responseRecorder, request, params)
				Expect(err).To(BeNil())
			})

			It("should return json", func() {
				handler(responseRecorder, request, params)
				contentTypes := responseRecorder.HeaderMap["Content-Type"]
				Expect(contentTypes).To(HaveLen(1))
				Expect(contentTypes[0]).To(Equal("application/json"))
			})

			It("should respond with a list of pages returned from Facebook", func() {
				handler(responseRecorder, request, params)
				var result *FacebookPage
				json.Unmarshal(responseRecorder.Body.Bytes(), &result)
				Expect(result.ID).To(Equal("id1"))
				Expect(result.Name).To(Equal("name1"))
				Expect(result.Address).To(Equal("a_street 10, a_city, a_country"))
				Expect(result.Phone).To(Equal("a_phone"))
				Expect(result.Website).To(Equal("a_website"))
				Expect(result.Email).To(Equal("an_email"))
			})
		})
	})
})
