package db_test

import (
	"github.com/deiwin/luncher-api/db"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

var (
	dbClient              *db.Client
	offersCollection      db.Offers
	tagsCollection        db.Tags
	restaurantsCollection db.Restaurants
	usersCollection       db.Users
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
	initUsersCollection()
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

func initUsersCollection() {
	usersCollection = db.NewUsers(dbClient)
	err := insertUsers()
	Expect(err).NotTo(HaveOccurred())
}

func createTestDbConf() (dbConfig *db.Config) {
	dbConfig = &db.Config{
		DbURL:  "mongodb://localhost/test",
		DbName: "test",
	}
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

func insertUsers() (err error) {
	return usersCollection.Insert(mocks.users...)
}

func wipeDb() {
	err := dbClient.DropDb()
	Expect(err).NotTo(HaveOccurred())
}
