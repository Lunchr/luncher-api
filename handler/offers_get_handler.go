package handler

import (
	"net/http"
	"time"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	"github.com/deiwin/luncher-api/router"
	"github.com/deiwin/luncher-api/storage"
)

// RegionOffers handles GET requests to /regions/:name/offers. It returns all
// current day's offers for the region.
func RegionOffers(offersCollection db.Offers, regionsCollection db.Regions, imageStorage storage.Images) router.HandlerWithParams {
	handler := func(w http.ResponseWriter, r *http.Request, region *model.Region) *router.HandlerError {
		loc, err := time.LoadLocation(region.Location)
		if err != nil {
			return &router.HandlerError{err, "The location of this region is misconfigured", http.StatusInternalServerError}
		}
		now := time.Now()
		startTime := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, loc)
		endTime := startTime.AddDate(0, 0, 1)
		offers, err := offersCollection.GetForRegion(region.Name, startTime, endTime)
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
	return forRegion(regionsCollection, handler)
}
