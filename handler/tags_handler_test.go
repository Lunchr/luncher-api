package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	. "github.com/deiwin/luncher-api/handler"
	. "github.com/deiwin/luncher-api/router"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TagsHandler", func() {
	var (
		mockTagsCollection db.Tags
		handler            Handler
	)

	BeforeEach(func() {
		mockTagsCollection = &mockTags{}
	})

	JustBeforeEach(func() {
		handler = Tags(mockTagsCollection)
	})

	Describe("Get", func() {
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

		Context("with simple mocked result from DB", func() {
			var (
				mockResult []*model.Tag
			)
			BeforeEach(func() {
				mockResult = []*model.Tag{&model.Tag{Name: "sometag"}}
				mockTagsCollection = &mockTags{
					func() (tags []*model.Tag, err error) {
						tags = mockResult
						return
					},
					nil,
				}
			})

			It("should write the returned data to responsewriter", func(done Done) {
				defer close(done)
				handler(responseRecorder, request)
				var result []*model.Tag
				json.Unmarshal(responseRecorder.Body.Bytes(), &result)
				Expect(result).To(HaveLen(1))
				Expect(result[0].Name).To(Equal(mockResult[0].Name))
			})
		})

		Context("with an error returned from the DB", func() {
			var dbErr = errors.New("DB stuff failed")

			BeforeEach(func() {
				mockTagsCollection = &mockTags{
					func() (tags []*model.Tag, err error) {
						err = dbErr
						return
					},
					nil,
				}
			})

			It("should return error 500", func(done Done) {
				defer close(done)
				err := handler(responseRecorder, request)
				Expect(err.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})

type mockTags struct {
	getFunc func() ([]*model.Tag, error)
	db.Tags
}

func (mock mockTags) Get() (tags []*model.Tag, err error) {
	if mock.getFunc != nil {
		tags, err = mock.getFunc()
	}
	return
}
