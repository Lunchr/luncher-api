package db_test

import (
	"time"

	"github.com/Lunchr/luncher-api/db/model"
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
	userID       bson.ObjectId
}

func createMocks() *Mocks {
	restaurantID := bson.NewObjectId()
	userID := bson.NewObjectId()
	return &Mocks{
		restaurantID: restaurantID,
		userID:       userID,
		offers: []*model.Offer{
			&model.Offer{
				CommonOfferFields: model.CommonOfferFields{
					Restaurant: model.OfferRestaurant{
						ID:     restaurantID,
						Name:   "Asian Chef",
						Region: "Tartu",
						Location: model.Location{
							Type:        "Point",
							Coordinates: []float64{26.72, 58.37},
						},
						Phone: "+372 5678 910",
					},
					Title:       "Sweet & Sour Chicken",
					Ingredients: []string{"Kana", "aedviljad", "tsillikaste"},
					FromTime:    parseTime("2014-11-10T09:00:00.000Z"),
					ToTime:      parseTime("2014-11-10T11:00:00.000Z"),
					Price:       3.4,
					Tags:        []string{"lind"},
				},
				ImageChecksum: "08446744073709551615",
			},
			&model.Offer{
				CommonOfferFields: model.CommonOfferFields{
					Restaurant: model.OfferRestaurant{
						Name:   "Bulgarian Dude",
						Region: "Tallinn",
						Location: model.Location{
							Type:        "Point",
							Coordinates: []float64{24.74, 59.42},
						},
						Phone: "+372 5678 910",
					},
					Title:       "Sweet & Sour Pork",
					Ingredients: []string{"Seafilee", "aedviljad", "mahushapu kaste"},
					FromTime:    parseTime("2014-11-10T09:00:00.000Z"),
					ToTime:      parseTime("2014-11-10T12:00:00.000Z"),
					Price:       3.3,
					Tags:        []string{"lind"},
				},
				ImageChecksum: "07446744073709551615",
			},
			&model.Offer{
				CommonOfferFields: model.CommonOfferFields{
					Restaurant: model.OfferRestaurant{
						Name:   "Caesarian Kitchen",
						Region: "Tartu",
						Location: model.Location{
							Type:        "Point",
							Coordinates: []float64{26.73, 58.36},
						},
						Phone: "+372 5678 910",
					},
					Title:       "Sweet & Sour Duck",
					Ingredients: []string{"Pardifilee", "aedviljad", "magushapu kaste"},
					FromTime:    parseTime("2014-11-12T09:00:00.000Z"),
					ToTime:      parseTime("2014-11-12T11:00:00.000Z"),
					Price:       3.6,
					Tags:        []string{"lind"},
				},
				ImageChecksum: "06446744073709551615",
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
				CCTLD:    "ee",
			},
			&model.Region{
				Name:     "Tallinn",
				Location: "Europe/Tallinn",
				CCTLD:    "ee",
			},
			&model.Region{
				Name:     "London",
				Location: "Europe/London",
				CCTLD:    "uk",
			},
		},
		restaurants: []*model.Restaurant{
			&model.Restaurant{
				FacebookPageID: "1",
				Name:           "Bulgarian Dude",
				Address:        "Võru 23, Tallinn",
				Region:         "Tallinn",
				Location: model.Location{
					Type:        "Point",
					Coordinates: []float64{24.74, 59.42},
				},
				Phone: "+372 5678 910",
			},
			&model.Restaurant{
				FacebookPageID: facebookPageID,
				ID:             restaurantID,
				Name:           "Asian Chef",
				Address:        "Võru 24, Tartu",
				Region:         "Tartu",
				Location: model.Location{
					Type:        "Point",
					Coordinates: []float64{26.72, 58.37},
				},
				Phone: "+372 5678 910",
			},
			&model.Restaurant{
				FacebookPageID: "3",
				Name:           "Caesarian Kitchen",
				Address:        "Võru 25, Tartu",
				Region:         "Tartu",
				Location: model.Location{
					Type:        "Point",
					Coordinates: []float64{26.73, 58.36},
				},
				Phone: "+372 5678 910",
			},
		},
		users: []*model.User{
			&model.User{
				ID:             userID,
				RestaurantIDs:  []bson.ObjectId{restaurantID},
				FacebookUserID: facebookUserID,
			},
			&model.User{
				RestaurantIDs:  []bson.ObjectId{restaurantID},
				FacebookUserID: "another user",
			},
		},
	}
}

func parseTime(timeString string) time.Time {
	parsedTime, err := time.Parse(time.RFC3339, timeString)
	Expect(err).NotTo(HaveOccurred())
	return parsedTime
}
