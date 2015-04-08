package handler

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	"github.com/deiwin/luncher-api/geo"
	"github.com/deiwin/luncher-api/router"
	"github.com/deiwin/luncher-api/storage"
)

// RegionOffers handles GET requests to /regions/:name/offers. It returns all
// current day's offers for the region.
func RegionOffers(offersCollection db.Offers, regionsCollection db.Regions, imageStorage storage.Images) router.HandlerWithParams {
	handler := func(w http.ResponseWriter, r *http.Request, region *model.Region) *router.HandlerError {
		timeLocation, err := time.LoadLocation(region.Location)
		if err != nil {
			return &router.HandlerError{err, "The location of this region is misconfigured", http.StatusInternalServerError}
		}
		startTime, endTime := getTodaysTimeRange(timeLocation)
		offers, err := offersCollection.GetForRegion(region.Name, startTime, endTime)
		if err != nil {
			return &router.HandlerError{err, "An error occured while trying to fetch today's offers", http.StatusInternalServerError}
		}
		if handlerError := changeOfferImageChecksumsToPaths(offers, imageStorage); handlerError != nil {
			return handlerError
		}
		return writeJSON(w, offers)
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
		// XXX how do we figure out the timezone for the request? i don't want to use
		// some API to guess the timezone from the location, as this will slow things
		// down. Probably best to expect something from the client
		timeLocation, err := time.LoadLocation("Local")
		if err != nil {
			return &router.HandlerError{err, "", http.StatusInternalServerError}
		}
		startTime, endTime := getTodaysTimeRange(timeLocation)
		offers, err := offersCollection.GetNear(loc, startTime, endTime)
		if err != nil {
			return &router.HandlerError{err, "An error occured while trying to fetch today's offers", http.StatusInternalServerError}
		}
		if handlerError := changeOfferImageChecksumsToPaths(offers, imageStorage); handlerError != nil {
			return handlerError
		}
		return writeJSON(w, offers)
	}
}

func getTodaysTimeRange(timeLocation *time.Location) (startTime, endTime time.Time) {
	now := time.Now()
	startTime = time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, timeLocation)
	endTime = startTime.AddDate(0, 0, 1)
	return startTime, endTime
}

func changeOfferImageChecksumsToPaths(offers []*model.Offer, imageStorage storage.Images) *router.HandlerError {
	var err error
	for _, offer := range offers {
		if offer.Image != "" {
			offer.Image, err = imageStorage.PathForLarge(offer.Image)
			if err != nil {
				return &router.HandlerError{err, "Failed to find an image for an offer", http.StatusInternalServerError}
			}
		}
	}
	return nil
}

func getLocFromRequest(r *http.Request) (geo.Location, *router.HandlerError) {
	latString := r.FormValue("lat")
	if latString == "" {
		return geo.Location{}, &router.HandlerError{errors.New("Latitude not specified"), "Please specify your latitude using the 'lat' attribute.", http.StatusBadRequest}
	}
	lat, err := strconv.ParseFloat(latString, 64)
	if err != nil {
		return geo.Location{}, &router.HandlerError{err, "Couldn't parse the latitude", http.StatusBadRequest}
	}
	lngString := r.FormValue("lng")
	if lngString == "" {
		return geo.Location{}, &router.HandlerError{errors.New("Longitude not specified"), "Please specify your longitude using the 'lng' attribute.", http.StatusBadRequest}
	}
	lng, err := strconv.ParseFloat(lngString, 64)
	if err != nil {
		return geo.Location{}, &router.HandlerError{err, "Couldn't parse the longitude", http.StatusBadRequest}
	}
	return geo.Location{
		Lat: lat,
		Lng: lng,
	}, nil

}
