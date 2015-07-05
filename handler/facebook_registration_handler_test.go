package handler_test

import (
	"encoding/json"
	"net/http"

	"github.com/Lunchr/luncher-api/db"
	. "github.com/Lunchr/luncher-api/handler"
	"github.com/Lunchr/luncher-api/handler/mocks"
	"github.com/Lunchr/luncher-api/router"
	"github.com/Lunchr/luncher-api/session"
	"github.com/deiwin/facebook/model"
	"github.com/julienschmidt/httprouter"
	"github.com/stretchr/testify/mock"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Facebook registration handler", func() {
	var (
		registrationAuther  *mocks.Authenticator
		mockSessMgr         session.Manager
		mockUsersCollection db.Users
		handlers            Facebook
	)

	BeforeEach(func() {
		registrationAuther = new(mocks.Authenticator)
		mockSessMgr = &mockSessionManager{}
	})

	JustBeforeEach(func() {
		handlers = NewFacebook(nil, registrationAuther, mockSessMgr, mockUsersCollection)
	})

	Describe("Login", func() {

		BeforeEach(func() {
			registrationAuther.On("AuthURL", "session").Return(testURL)
		})

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

	Describe("ListPagesManagedByUser", func() {
		var (
			handler router.Handler
		)

		JustBeforeEach(func() {
			handler = handlers.ListPagesManagedByUser()
		})

		ExpectUserToBeLoggedIn(func() *router.HandlerError {
			return handler(responseRecorder, request)
		}, func(mgr session.Manager, users db.Users) {
			mockSessMgr = mgr
			mockUsersCollection = users
		})

		Context("with user logged in", func() {
			var (
				api *mocks.API
			)
			BeforeEach(func() {
				mockSessMgr = &mockSessionManager{isSet: true, id: "correctSession"}
				mockUsersCollection = mockUsers{}
				api = new(mocks.API)
				registrationAuther.On("APIConnection", mock.AnythingOfType("*oauth2.Token")).Return(api)
				api.On("Accounts").Return(&model.Accounts{
					Data: []model.Page{
						model.Page{
							ID:   "id1",
							Name: "name1",
						},
						model.Page{
							ID:   "id2",
							Name: "name2",
						},
					},
				}, nil)
			})

			AfterEach(func() {
				registrationAuther.AssertExpectations(GinkgoT())
				api.AssertExpectations(GinkgoT())
			})

			It("should succeed", func(done Done) {
				defer close(done)
				err := handler(responseRecorder, request)
				Expect(err).To(BeNil())
			})

			It("should return json", func(done Done) {
				defer close(done)
				handler(responseRecorder, request)
				contentTypes := responseRecorder.HeaderMap["Content-Type"]
				Expect(contentTypes).To(HaveLen(1))
				Expect(contentTypes[0]).To(Equal("application/json"))
			})

			It("should respond with a list of pages returned from Facebook", func() {
				handler(responseRecorder, request)
				var result []*Page
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
		var (
			handler router.HandlerWithParams
		)

		JustBeforeEach(func() {
			handler = handlers.Page()
		})

		ExpectUserToBeLoggedIn(func() *router.HandlerError {
			return handler(responseRecorder, request, nil)
		}, func(mgr session.Manager, users db.Users) {
			mockSessMgr = mgr
			mockUsersCollection = users
		})

		Context("with user logged in", func() {
			var (
				api    *mocks.API
				params httprouter.Params
			)
			BeforeEach(func() {
				mockSessMgr = &mockSessionManager{isSet: true, id: "correctSession"}
				mockUsersCollection = mockUsers{}
				api = new(mocks.API)
				registrationAuther.On("APIConnection", mock.AnythingOfType("*oauth2.Token")).Return(api)
				api.On("Page", "a_page_id").Return(&model.Page{
					ID:   "id1",
					Name: "name1",
					Location: model.Location{
						Street:  "a_street 10",
						City:    "a_city",
						Country: "a_country",
					},
					Phone:   "a_phone",
					Website: "a_website",
				}, nil)
				params = httprouter.Params{httprouter.Param{
					Key:   "id",
					Value: "a_page_id",
				}}
			})

			AfterEach(func() {
				registrationAuther.AssertExpectations(GinkgoT())
				api.AssertExpectations(GinkgoT())
			})

			It("should succeed", func(done Done) {
				defer close(done)
				err := handler(responseRecorder, request, params)
				Expect(err).To(BeNil())
			})

			It("should return json", func(done Done) {
				defer close(done)
				handler(responseRecorder, request, params)
				contentTypes := responseRecorder.HeaderMap["Content-Type"]
				Expect(contentTypes).To(HaveLen(1))
				Expect(contentTypes[0]).To(Equal("application/json"))
			})

			It("should respond with a list of pages returned from Facebook", func() {
				handler(responseRecorder, request, params)
				var result *Page
				json.Unmarshal(responseRecorder.Body.Bytes(), &result)
				Expect(result.ID).To(Equal("id1"))
				Expect(result.Name).To(Equal("name1"))
				Expect(result.Address).To(Equal("a_street 10, a_city, a_country"))
				Expect(result.Phone).To(Equal("a_phone"))
				Expect(result.Website).To(Equal("a_website"))
			})
		})
	})
})
