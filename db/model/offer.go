package model

import (
	"time"

	"gopkg.in/mgo.v2/bson"
)

const OfferCollectionName = "offers"

type (
	CommonOfferFields struct {
		ID          bson.ObjectId   `json:"_id,omitempty"        bson:"_id,omitempty"`
		Restaurant  OfferRestaurant `json:"restaurant"           bson:"restaurant"`
		Title       string          `json:"title"                bson:"title"`
		FromTime    time.Time       `json:"from_time"            bson:"from_time"`
		ToTime      time.Time       `json:"to_time"              bson:"to_time"`
		Description string          `json:"description"          bson:"description"`
		Price       float64         `json:"price"                bson:"price"`
		Tags        []string        `json:"tags"                 bson:"tags"`
	}

	// Offer provides the mapping to the offers as represented in the DB
	Offer struct {
		// The bson marshaller (unlike the json one) doesn't automatically inline
		// embedded fields, so the inline tag has to be specified.
		CommonOfferFields `bson:",inline"`
		ImageChecksum     string `bson:"image_checksum,omitempty"`
	}

	// OfferJSON is the view of an offer that gets sent to the users
	OfferJSON struct {
		CommonOfferFields
		Image *OfferImagePaths `json:"image,omitempty"`
	}

	// OfferJSON is the view of an offer that gets sent to the users
	OfferPOST struct {
		CommonOfferFields
		ImageData string `json:"image_data,omitempty"`
	}

	// OfferImagePaths holds paths to the various sizes of the offer's image
	OfferImagePaths struct {
		Large     string `json:"large"`
		Thumbnail string `json:"thumbnail"`
	}

	// OfferRestaurant holds the information about the restaurant that gets included
	// in every offer
	OfferRestaurant struct {
		ID       bson.ObjectId `json:"id"       bson:"id"`
		Name     string        `json:"name"     bson:"name"`
		Region   string        `json:"region"   bson:"region"`
		Address  string        `json:"address"  bson:"address"`
		Location Location      `json:"location" bson:"location"`
		Phone    string        `json:"phone"    bson:"phone"`
	}

	// OfferRestaurantWithDistance wraps an OfferRestaurant and adds a distance field.
	// This struct can be used to respond to queries about nearby offers.
	OfferRestaurantWithDistance struct {
		OfferRestaurant
		Distance float64 `json:"distance"`
	}

	// OfferWithDistance wraps an offer and adds a distance field to the included
	// restaurant struct. This struct can be used to respond to queries about nearby offers.
	OfferWithDistance struct {
		Offer
		Restaurant OfferRestaurantWithDistance `json:"restaurant"`
	}

	// OfferWithDistance wraps an offer's JSON representation and adds a distance field to the included
	// restaurant struct. This struct can be used to respond to queries about nearby offers.
	OfferWithDistanceJSON struct {
		OfferJSON
		Restaurant OfferRestaurantWithDistance `json:"restaurant"`
	}
)

func MapOfferToJSON(offer *Offer, imageToPathMapper func(string) (*OfferImagePaths, error)) (*OfferJSON, error) {
	image, err := imageToPathMapper(offer.ImageChecksum)
	if err != nil {
		return nil, err
	}
	return &OfferJSON{
		CommonOfferFields: offer.CommonOfferFields,
		Image:             image,
	}, nil
}

func MapOfferPOSTToOffer(offer *OfferPOST, imageDataToChecksumMapper func(string) (string, error)) (*Offer, error) {
	imageChecksum, err := imageDataToChecksumMapper(offer.ImageData)
	if err != nil {
		return nil, err
	}
	return &Offer{
		CommonOfferFields: offer.CommonOfferFields,
		ImageChecksum:     imageChecksum,
	}, nil
}

func MapOfferWithDistanceToJSON(offer *OfferWithDistance, imageToPathMapper func(string) (*OfferImagePaths, error)) (*OfferWithDistanceJSON, error) {
	offerJSON, err := MapOfferToJSON(&offer.Offer, imageToPathMapper)
	if err != nil {
		return nil, err
	}
	return &OfferWithDistanceJSON{
		OfferJSON:  *offerJSON,
		Restaurant: offer.Restaurant,
	}, nil
}
