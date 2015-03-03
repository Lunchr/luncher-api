package handler_test

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strconv"

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
	requestPath      string
	requestData      url.Values
	requestQuery     url.Values
)

var _ = BeforeEach(func(done Done) {
	defer close(done)
	responseRecorder = httptest.NewRecorder()
})

var _ = JustBeforeEach(func() {
	var err error
	requestURL := (&url.URL{
		Scheme:   "http",
		Host:     "localhost",
		Path:     requestPath,
		RawQuery: requestQuery.Encode(),
	}).String()
	request, err = http.NewRequest(requestMethod, requestURL, bytes.NewBufferString(requestData.Encode()))
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Add("Content-Length", strconv.Itoa(len(requestData.Encode())))
	Expect(err).NotTo(HaveOccurred())
})
