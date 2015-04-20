package handler_test

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	. "github.com/deiwin/luncher-api/handler"
	"github.com/deiwin/luncher-api/router"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RegionsHandler", func() {
	var (
		mockRegionsCollection db.Regions
		handler               router.Handler
	)

	BeforeEach(func() {
		mockRegionsCollection = &mockRegions{
			getAllFunc: func() db.RegionIter {
				return &mockRegionIter{}
			},
		}
	})

	JustBeforeEach(func() {
		handler = Regions(mockRegionsCollection)
	})

	Describe("Regions", func() {
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
				mockResult []*model.Region
			)
			BeforeEach(func() {
				mockResult = []*model.Region{
					&model.Region{Name: "someregion"},
					&model.Region{Name: "someregion2"},
				}
				mockRegionsCollection = &mockRegions{
					getAllFunc: func() db.RegionIter {
						return &mockRegionIter{mockResult: mockResult}
					},
				}
			})

			It("should write the returned data to responsewriter", func(done Done) {
				defer close(done)
				handler(responseRecorder, request)
				var result []*model.Region
				json.Unmarshal(responseRecorder.Body.Bytes(), &result)
				Expect(result).To(HaveLen(2))
				Expect(result[0].Name).To(Equal(mockResult[0].Name))
				Expect(result[1].Name).To(Equal(mockResult[1].Name))
			})
		})

		Context("with an error returned from the DB", func() {
			var dbErr = errors.New("DB stuff failed")

			BeforeEach(func() {
				mockRegionsCollection = &mockRegions{
					getAllFunc: func() db.RegionIter {
						return &mockRegionIter{err: dbErr}
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

func (m mockRegions) GetAll() db.RegionIter {
	if m.getAllFunc != nil {
		return m.getAllFunc()
	}
	return nil
}

type mockRegionIter struct {
	mockResult []*model.Region
	i          int
	err        error
	db.RegionIter
}

func (m *mockRegionIter) Next(region *model.Region) bool {
	if m.err != nil {
		return false
	}
	if m.i >= len(m.mockResult) {
		return false
	}
	*region = *m.mockResult[m.i]
	m.i++
	return true
}

func (m mockRegionIter) Close() error {
	return m.err
}
