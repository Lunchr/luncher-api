package handler

import (
	"errors"
	"net/http"

	"github.com/deiwin/luncher-api/db"
	"github.com/deiwin/luncher-api/db/model"
	"github.com/deiwin/luncher-api/router"
	"github.com/julienschmidt/httprouter"
)

type HandlerWithRegion func(w http.ResponseWriter, r *http.Request, region *model.Region) *router.HandlerError

func forRegion(regionsCollection db.Regions, handler HandlerWithRegion) router.HandlerWithParams {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) *router.HandlerError {
		regionName := ps.ByName("name")
		if regionName == "" {
			return &router.HandlerError{errors.New("Region not specified for GET /regions/:name"), "Please specify a region", http.StatusBadRequest}
		}
		region, err := regionsCollection.GetName(regionName)
		if err != nil {
			return &router.HandlerError{err, "Unable to find the specified region", http.StatusNotFound}
		}
		return handler(w, r, region)
	}
}
