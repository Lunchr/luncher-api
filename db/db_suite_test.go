package db_test

import (
	"github.com/deiwin/praad-api/db"
	"github.com/deiwin/praad-api/db/model"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
	"time"
)

var (
	dbClient              *db.Client
	offersCollection      *db.Offers
	tagsCollection        *db.Tags
	restaurantsCollection *db.Restaurants
)

func TestDb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Db Suite")
}

var _ = BeforeSuite(func(done Done) {
	defer close(done)
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
	return tagsCollection.Insert(
		&model.Tag{
			Name:        "kala",
			DisplayName: "Kalast",
		},
		&model.Tag{
			Name:        "lind",
			DisplayName: "Linnust",
		},
		&model.Tag{
			Name:        "siga",
			DisplayName: "Seast",
		},
		&model.Tag{
			Name:        "veis",
			DisplayName: "Veisest",
		},
		&model.Tag{
			Name:        "lammas",
			DisplayName: "Lambast",
		},
	)
}

func insertRestaurants() (err error) {
	return restaurantsCollection.Insert(
		&model.Restaurant{
			Name:    "Bulgarian Dude",
			Address: "Võru 23, Tartu",
		},
		&model.Restaurant{
			Name:    "Asian Chef",
			Address: "Võru 24, Tartu",
		},
		&model.Restaurant{
			Name:    "Caesarian Kitchen",
			Address: "Võru 25, Tartu",
		},
	)
}

func insertOffers() (err error) {
	return offersCollection.Insert(
		&model.Offer{
			Restaurant: model.OfferRestaurant{
				Name: "Asian Chef",
			},
			Title:       "Sweet & Sour Chicken",
			Description: "Kanafilee aedviljadega rikkalikus magushapus kastmes.",
			FromTime:    parseTime("2014-11-10T09:00:00.000Z"),
			ToTime:      parseTime("2014-11-10T11:00:00.000Z"),
			Price:       3.4,
			Tags:        []string{"lind"},
		},
		&model.Offer{
			Restaurant: model.OfferRestaurant{
				Name: "Bulgarian Dude",
			},
			Title:       "Sweet & Sour Pork",
			Description: "Seafilee aedviljadega rikkalikus magushapus kastmes.",
			FromTime:    parseTime("2014-11-10T09:00:00.000Z"),
			ToTime:      parseTime("2014-11-10T12:00:00.000Z"),
			Price:       3.3,
			Tags:        []string{"lind"},
		},
		&model.Offer{
			Restaurant: model.OfferRestaurant{
				Name: "Caesarian Kitchen",
			},
			Title:       "Sweet & Sour Duck",
			Description: "Pardifilee aedviljadega rikkalikus magushapus kastmes.",
			FromTime:    parseTime("2014-11-11T09:00:00.000Z"),
			ToTime:      parseTime("2014-11-11T11:00:00.000Z"),
			Price:       3.6,
			Tags:        []string{"lind"},
		},
	)
}

func wipeDb() {
	err := dbClient.DropDb()
	Expect(err).NotTo(HaveOccurred())
}

func parseTime(timeString string) time.Time {
	parsedTime, err := time.Parse(time.RFC3339, timeString)
	Expect(err).NotTo(HaveOccurred())
	return parsedTime
}
