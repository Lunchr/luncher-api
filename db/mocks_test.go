package db_test

import (
	"time"

	"github.com/deiwin/luncher-api/db/model"
	. "github.com/onsi/gomega"
	"gopkg.in/mgo.v2/bson"
)

type Mocks struct {
	offers       []*model.Offer
	tags         []*model.Tag
	regions      []*model.Region
	restaurants  []*model.Restaurant
	restaurantID bson.ObjectId
	users        []*model.User
}

func createMocks() *Mocks {
	restaurantID := bson.NewObjectId()
	return &Mocks{
		restaurantID: restaurantID,
		offers: []*model.Offer{&model.Offer{
			Restaurant: model.OfferRestaurant{
				Name:   "Asian Chef",
				Region: "Tartu",
			},
			Title:       "Sweet & Sour Chicken",
			Ingredients: []string{"Kana", "aedviljad", "tsillikaste"},
			FromTime:    parseTime("2014-11-10T09:00:00.000Z"),
			ToTime:      parseTime("2014-11-10T11:00:00.000Z"),
			Price:       3.4,
			Tags:        []string{"lind"},
			Image:       "08446744073709551615",
		},
			&model.Offer{
				Restaurant: model.OfferRestaurant{
					Name:   "Bulgarian Dude",
					Region: "Tallinn",
				},
				Title:       "Sweet & Sour Pork",
				Ingredients: []string{"Seafilee", "aedviljad", "mahushapu kaste"},
				FromTime:    parseTime("2014-11-10T09:00:00.000Z"),
				ToTime:      parseTime("2014-11-10T12:00:00.000Z"),
				Price:       3.3,
				Tags:        []string{"lind"},
				Image:       "07446744073709551615",
			},
			&model.Offer{
				Restaurant: model.OfferRestaurant{
					Name:   "Caesarian Kitchen",
					Region: "Tartu",
				},
				Title:       "Sweet & Sour Duck",
				Ingredients: []string{"Pardifilee", "aedviljad", "magushapu kaste"},
				FromTime:    parseTime("2014-11-12T09:00:00.000Z"),
				ToTime:      parseTime("2014-11-12T11:00:00.000Z"),
				Price:       3.6,
				Tags:        []string{"lind"},
				Image:       "06446744073709551615",
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
		regions: []*model.Region{
			&model.Region{
				Name:     "Tartu",
				Location: "Europe/Tallinn",
			},
			&model.Region{
				Name:     "Tallinn",
				Location: "Europe/Tallinn",
			},
			&model.Region{
				Name:     "London",
				Location: "Europe/London",
			},
		},
		restaurants: []*model.Restaurant{
			&model.Restaurant{
				Name:    "Bulgarian Dude",
				Address: "Võru 23, Tartu",
				Region:  "Tallinn",
			},
			&model.Restaurant{
				ID:      restaurantID,
				Name:    "Asian Chef",
				Address: "Võru 24, Tartu",
				Region:  "Tartu",
			},
			&model.Restaurant{
				Name:    "Caesarian Kitchen",
				Address: "Võru 25, Tartu",
				Region:  "Tartu",
			},
		},
		users: []*model.User{
			&model.User{
				RestaurantID:   restaurantID,
				FacebookUserID: facebookUserID,
				FacebookPageID: facebookPageID,
			},
			&model.User{
				RestaurantID:   restaurantID,
				FacebookUserID: "another user",
				FacebookPageID: facebookPageID,
			},
		},
	}
}

func parseTime(timeString string) time.Time {
	parsedTime, err := time.Parse(time.RFC3339, timeString)
	Expect(err).NotTo(HaveOccurred())
	return parsedTime
}
