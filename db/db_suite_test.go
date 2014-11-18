package db_test

import (
	"github.com/deiwin/praad-api/db"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var dbClient *db.Client

func TestDb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Db Suite")
}

var _ = BeforeSuite(func() {
	dbClient = db.NewClient()
	err := dbClient.Connect()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func() {
	dbClient.Disconnect()
})

var _ = It("should work", func() {

})
