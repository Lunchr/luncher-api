package db_test

import (
	"github.com/deiwin/praad-api/db"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	dbClient              *db.Client
	offersCollection      db.Offers
	tagsCollection        db.Tags
	restaurantsCollection db.Restaurants
	mocks                 *Mocks
)

func TestDb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Db Suite")
}

var _ = BeforeSuite(func(done Done) {
	defer close(done)
	mocks = createMocks()
	createClient()
	wipeDb()
	initCollections()
})

var _ = AfterSuite(func(done Done) {
	defer close(done)
	wipeDb()
	dbClient.Disconnect()
})

var _ = It("should work", func() {})

func createClient() {
	dbConfig := createTestDbConf()
	dbClient = db.NewClient(dbConfig)
	err := dbClient.Connect()
	Expect(err).NotTo(HaveOccurred())
}

func initCollections() {
	initOffersCollection()
	initTagsCollection()
	initRestaurantsCollection()
}

func initOffersCollection() {
	offersCollection = db.NewOffers(dbClient)
	err := insertOffers()
	Expect(err).NotTo(HaveOccurred())
}

func initTagsCollection() {
	tagsCollection = db.NewTags(dbClient)
	err := insertTags()
	Expect(err).NotTo(HaveOccurred())
}

func initRestaurantsCollection() {
	restaurantsCollection = db.NewRestaurants(dbClient)
	err := insertRestaurants()
	Expect(err).NotTo(HaveOccurred())
}

func createTestDbConf() (dbConfig *db.Config) {
	dbConfig = db.NewConfig()
	dbConfig.DbURL = "mongodb://localhost/test"
	dbConfig.DbName = "test"
	return
}

func insertTags() (err error) {
	return tagsCollection.Insert(mocks.tags...)
}

func insertRestaurants() (err error) {
	return restaurantsCollection.Insert(mocks.restaurants...)
}

func insertOffers() (err error) {
	return offersCollection.Insert(mocks.offers...)
}

func wipeDb() {
	err := dbClient.DropDb()
	Expect(err).NotTo(HaveOccurred())
}
