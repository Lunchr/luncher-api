package handler

import (
	"errors"
	"net/http"
	"time"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/router"
	"github.com/deiwin/luncher-api/storage"
)

// Offers handles GET requests to /offers. It returns all current day's offers.
func Offers(offersCollection db.Offers, regionsCollection db.Regions, imageStorage storage.Images) router.Handler {
	return func(w http.ResponseWriter, r *http.Request) *router.HandlerError {
		regionName := r.FormValue("region")
		if regionName == "" {
			return &router.HandlerError{errors.New("Region not specified for GET /offers"), "Please specify a region", http.StatusBadRequest}
		}
		region, err := regionsCollection.GetName(regionName)
		if err != nil {
			return &router.HandlerError{err, "Unable to find the specified region", http.StatusNotFound}
		}
		loc, err := time.LoadLocation(region.Location)
		if err != nil {
			return &router.HandlerError{err, "The location of this region is misconfigured", http.StatusInternalServerError}
		}
		now := time.Now()
		startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		endTime := startTime.AddDate(0, 0, 1)
		offers, err := offersCollection.GetForRegion(regionName, startTime, endTime)
		if err != nil {
			return &router.HandlerError{err, "An error occured while trying to fetch todays offers", http.StatusInternalServerError}
		}
		for _, offer := range offers {
			if offer.Image != "" {
				offer.Image, err = imageStorage.PathForLarge(offer.Image)
				if err != nil {
					return &router.HandlerError{err, "Failed to find an image for an offer", http.StatusInternalServerError}
				}
			}
		}
		return writeJSON(w, offers)
	}
}
