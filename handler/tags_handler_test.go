package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/db/model"
	. "github.com/Lunchr/luncher-api/handler"
	"github.com/Lunchr/luncher-api/router"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("TagsHandler", func() {
	var (
		mockTagsCollection db.Tags
		handler            router.Handler
	)

	BeforeEach(func() {
		mockTagsCollection = &mockTags{
			getAllFunc: func() db.TagIter {
				return &mockTagIter{}
			},
		}
	})

	JustBeforeEach(func() {
		handler = Tags(mockTagsCollection)
	})

	Describe("Tags", func() {
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
				mockResult = []*model.Tag{
					&model.Tag{Name: "sometag"},
					&model.Tag{Name: "sometag2"},
				}
				mockTagsCollection = &mockTags{
					getAllFunc: func() db.TagIter {
						return &mockTagIter{mockResult: mockResult}
					},
				}
			})

			It("should write the returned data to responsewriter", func(done Done) {
				defer close(done)
				handler(responseRecorder, request)
				var result []*model.Tag
				json.Unmarshal(responseRecorder.Body.Bytes(), &result)
				Expect(result).To(HaveLen(2))
				Expect(result[0].Name).To(Equal(mockResult[0].Name))
				Expect(result[1].Name).To(Equal(mockResult[1].Name))
			})
		})

		Context("with an error returned from the DB", func() {
			var dbErr = errors.New("DB stuff failed")

			BeforeEach(func() {
				mockTagsCollection = &mockTags{
					getAllFunc: func() db.TagIter {
						return &mockTagIter{err: dbErr}
					},
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

func (m mockTags) GetAll() db.TagIter {
	if m.getAllFunc != nil {
		return m.getAllFunc()
	}
	return nil
}

type mockTagIter struct {
	mockResult []*model.Tag
	i          int
	err        error
	db.TagIter
}

func (m *mockTagIter) Next(tag *model.Tag) bool {
	if m.err != nil {
		return false
	}
	if m.i >= len(m.mockResult) {
		return false
	}
	*tag = *m.mockResult[m.i]
	m.i++
	return true
}

func (m mockTagIter) Close() error {
	return m.err
}

type mockTags struct {
	getAllFunc func() db.TagIter
	db.Tags
}
