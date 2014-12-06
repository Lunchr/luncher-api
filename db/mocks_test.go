package db_test

import (
	"time"

	"github.com/deiwin/luncher-api/db/model"
	. "github.com/onsi/gomega"
)

type Mocks struct {
	offers      []*model.Offer
	tags        []*model.Tag
	restaurants []*model.Restaurant
}

func createMocks() *Mocks {
	return &Mocks{
		offers: []*model.Offer{&model.Offer{
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
				FromTime:    parseTime("2014-11-12T09:00:00.000Z"),
				ToTime:      parseTime("2014-11-12T11:00:00.000Z"),
				Price:       3.6,
				Tags:        []string{"lind"},
			},
		},
		tags: []*model.Tag{
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
		},
		restaurants: []*model.Restaurant{
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
		},
	}
}

func parseTime(timeString string) time.Time {
	parsedTime, err := time.Parse(time.RFC3339, timeString)
	Expect(err).NotTo(HaveOccurred())
	return parsedTime
}
