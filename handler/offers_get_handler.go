package handler

import (
	"net/http"
	"strconv"
	"time"

	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/db/model"
	"github.com/Lunchr/luncher-api/geo"
	"github.com/Lunchr/luncher-api/router"
	"github.com/Lunchr/luncher-api/storage"
	"github.com/bradfitz/latlong"
)

// RegionOffers handles GET requests to /regions/:name/offers. It returns all
// current day's offers for the region.
func RegionOffers(offersCollection db.Offers, regionsCollection db.Regions, imageStorage storage.Images) router.HandlerWithParams {
	handler := func(w http.ResponseWriter, r *http.Request, region *model.Region) *router.HandlerError {
		timeLocation, err := time.LoadLocation(region.Location)
		if err != nil {
			return router.NewHandlerError(err, "The location of this region is misconfigured", http.StatusInternalServerError)
		}
		startTime, endTime := getTodaysTimeRange(timeLocation)
		offers, err := offersCollection.GetForRegion(region.Name, startTime, endTime)
		if err != nil {
			return router.NewHandlerError(err, "An error occured while trying to fetch today's offers", http.StatusInternalServerError)
		}
		offerJSONs, handlerError := mapOffersToJSON(offers, imageStorage)
		if handlerError != nil {
			return handlerError
		}
		return writeJSON(w, offerJSONs)
	}
	return forRegion(regionsCollection, handler)
}

// ProximalOffers handles requests that whish to know about offers near a certain
// location.
func ProximalOffers(offersCollection db.Offers, imageStorage storage.Images) router.Handler {
	return func(w http.ResponseWriter, r *http.Request) *router.HandlerError {
		loc, handlerError := getLocFromRequest(r)
		if handlerError != nil {
			return handlerError
		}
		zone := latlong.LookupZoneName(loc.Lat, loc.Lng)
		if zone == "" {
			message := "Failed to find a timezone for this location"
			return router.NewSimpleHandlerError(message, http.StatusInternalServerError)
		}
		timeLocation, err := time.LoadLocation(zone)
		if err != nil {
			return router.NewHandlerError(err, "", http.StatusInternalServerError)
		}
		startTime, endTime := getTodaysTimeRange(timeLocation)
		offers, err := offersCollection.GetNear(loc, startTime, endTime)
		if err != nil {
			return router.NewHandlerError(err, "An error occured while trying to fetch today's offers", http.StatusInternalServerError)
		}
		offerJSONs, handlerError := mapOffersWithDistanceToJSON(offers, imageStorage)
		if handlerError != nil {
			return handlerError
		}
		return writeJSON(w, offerJSONs)
	}
}

func getTodaysTimeRange(timeLocation *time.Location) (startTime, endTime time.Time) {
	now := time.Now()
	startTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, timeLocation)
	endTime = startTime.AddDate(0, 0, 1)
	return startTime, endTime
}

func mapOffersToJSON(offers []*model.Offer, imageStorage storage.Images) ([]*model.OfferJSON, *router.HandlerError) {
	offerJSONs := make([]*model.OfferJSON, len(offers))
	for i, offer := range offers {
		offerJSON, err := mapOfferToJSON(offer, imageStorage)
		if err != nil {
			return nil, err
		}
		offerJSONs[i] = offerJSON
	}
	return offerJSONs, nil
}

func mapOfferToJSON(offer *model.Offer, imageStorage storage.Images) (*model.OfferJSON, *router.HandlerError) {
	offerJSON, err := model.MapOfferToJSON(offer, imageStorage.PathsFor)
	if err != nil {
		return nil, router.NewHandlerError(err, "Failed to map offers to JSON", http.StatusInternalServerError)
	}
	return offerJSON, nil
}

func mapOffersWithDistanceToJSON(offers []*model.OfferWithDistance, imageStorage storage.Images) ([]*model.OfferWithDistanceJSON, *router.HandlerError) {
	offerJSONs := make([]*model.OfferWithDistanceJSON, len(offers))
	for i, offer := range offers {
		offerJSON, err := mapOfferWithDistanceToJSON(offer, imageStorage)
		if err != nil {
			return nil, err
		}
		offerJSONs[i] = offerJSON
	}
	return offerJSONs, nil
}

func mapOfferWithDistanceToJSON(offer *model.OfferWithDistance, imageStorage storage.Images) (*model.OfferWithDistanceJSON, *router.HandlerError) {
	offerJSON, err := model.MapOfferWithDistanceToJSON(offer, imageStorage.PathsFor)
	if err != nil {
		return nil, router.NewHandlerError(err, "Failed to map offers to JSON", http.StatusInternalServerError)
	}
	return offerJSON, nil
}

func getLocFromRequest(r *http.Request) (geo.Location, *router.HandlerError) {
	latString := r.FormValue("lat")
	if latString == "" {
		return geo.Location{}, router.NewStringHandlerError("Latitude not specified", "Please specify your latitude using the 'lat' attribute.", http.StatusBadRequest)
	}
	lat, err := strconv.ParseFloat(latString, 64)
	if err != nil {
		return geo.Location{}, router.NewHandlerError(err, "Couldn't parse the latitude", http.StatusBadRequest)
	}
	lngString := r.FormValue("lng")
	if lngString == "" {
		return geo.Location{}, router.NewStringHandlerError("Longitude not specified", "Please specify your longitude using the 'lng' attribute.", http.StatusBadRequest)
	}
	lng, err := strconv.ParseFloat(lngString, 64)
	if err != nil {
		return geo.Location{}, router.NewHandlerError(err, "Couldn't parse the longitude", http.StatusBadRequest)
	}
	return geo.Location{
		Lat: lat,
		Lng: lng,
	}, nil

}
