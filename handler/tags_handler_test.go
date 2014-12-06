package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	. "github.com/deiwin/luncher-api/handler"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TagsHandler", func() {
	var (
		mockTagsCollection db.Tags
		handler            handlerFunc
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
			handler(responseRecorder, request)
			Expect(responseRecorder.Code).To(Equal(http.StatusOK))
		})

		It("should return json", func(done Done) {
			defer close(done)
			handler(responseRecorder, request)
			contentTypes := responseRecorder.HeaderMap["Content-Type"]
			Expect(contentTypes).To(HaveLen(1))
			Expect(contentTypes[0]).To(Equal("application/json"))
			// TODO the header assertion could be made a custom matcher
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
				}
			})

			It("should return error 500", func(done Done) {
				defer close(done)
				handler(responseRecorder, request)
				Expect(responseRecorder.Code).To(Equal(http.StatusInternalServerError))
			})
		})
	})
})

type mockTags struct {
	getFunc func() ([]*model.Tag, error)
}

func (mock mockTags) Insert(tagsToInsert ...*model.Tag) (err error) {
	return
}

func (mock mockTags) Get() (tags []*model.Tag, err error) {
	if mock.getFunc != nil {
		tags, err = mock.getFunc()
	}
	return
}
