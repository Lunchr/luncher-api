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
	requestMethod    = "GET"
	requestURL       string
)

var _ = BeforeEach(func(done Done) {
	defer close(done)
	responseRecorder = httptest.NewRecorder()

})

var _ = JustBeforeEach(func() {
	var err error
	request, err = http.NewRequest(requestMethod, requestURL, nil)
	Expect(err).NotTo(HaveOccurred())
})
