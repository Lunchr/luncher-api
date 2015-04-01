package handler_test

import (
	"bytes"
	"encoding/json"
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
	requestData      interface{}
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
	data, err := json.Marshal(requestData)
	Expect(err).NotTo(HaveOccurred())
	request, err = http.NewRequest(requestMethod, requestURL, bytes.NewBuffer(data))
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Content-Length", strconv.Itoa(len(data)))
	Expect(err).NotTo(HaveOccurred())
})
