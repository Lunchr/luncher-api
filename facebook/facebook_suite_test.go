package facebook_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestFacebook(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Facebook Suite")
}
