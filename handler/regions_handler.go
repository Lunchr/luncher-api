package handler

import (
	"net/http"

	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/db/model"
	"github.com/Lunchr/luncher-api/router"
	"github.com/julienschmidt/httprouter"
)

func Regions(regionsCollection db.Regions) router.Handler {
	return func(w http.ResponseWriter, r *http.Request) *router.HandlerError {
		regionsIter := regionsCollection.GetAll()
		var regions []model.Region
		var region model.Region
		for regionsIter.Next(&region) {
			regions = append(regions, region)
		}
		if err := regionsIter.Close(); err != nil {
			return router.NewHandlerError(err, "An error occured while fetching the regions from the DB", http.StatusInternalServerError)
		}
		return writeJSON(w, regions)
	}
}

type HandlerWithRegion func(w http.ResponseWriter, r *http.Request, region *model.Region) *router.HandlerError

func forRegion(regionsCollection db.Regions, handler HandlerWithRegion) router.HandlerWithParams {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) *router.HandlerError {
		regionName := ps.ByName("name")
		if regionName == "" {
			return router.NewStringHandlerError("Region not specified for GET /regions/:name", "Please specify a region", http.StatusBadRequest)
		}
		region, err := regionsCollection.GetName(regionName)
		if err != nil {
			return router.NewHandlerError(err, "Unable to find the specified region", http.StatusNotFound)
		}
		return handler(w, r, region)
	}
}
