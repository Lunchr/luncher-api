package handler_test

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestHandler(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Handler Suite")
}

var (
	responseRecorder *httptest.ResponseRecorder
	request          *http.Request
)

var _ = BeforeEach(func(done Done) {
	defer close(done)
	responseRecorder = httptest.NewRecorder()
	var err error
	request, err = http.NewRequest("", "", nil)
	Expect(err).NotTo(HaveOccurred())
})
