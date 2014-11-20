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
	dbClient         *db.Client
	offersCollection *db.Offers
)

func TestDb(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Db Suite")
}

var _ = BeforeSuite(func(done Done) {
	defer close(done)
	dbConfig := createTestDbConf()
	dbClient = db.NewClient(dbConfig)
	err := dbClient.Connect()
	Expect(err).NotTo(HaveOccurred())

	err = wipeDb()
	Expect(err).NotTo(HaveOccurred())

	offersCollection = db.NewOffers(dbClient)
	err = insertOffers()
	Expect(err).NotTo(HaveOccurred())
})

var _ = AfterSuite(func(done Done) {
	defer close(done)
	err := wipeDb()
	Expect(err).NotTo(HaveOccurred())

	dbClient.Disconnect()
})

var _ = It("should work", func() {

})

func createTestDbConf() (dbConfig *db.Config) {
	dbConfig = db.NewConfig()
	dbConfig.DbURL = "mongodb://localhost/test"
	dbConfig.DbName = "test"
	return
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

func wipeDb() (err error) {
	return dbClient.DropDb()
}

func parseTime(timeString string) time.Time {
	parsedTime, err := time.Parse(time.RFC3339, timeString)
	Expect(err).NotTo(HaveOccurred())
	return parsedTime
}
