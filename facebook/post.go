package facebook

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"gopkg.in/mgo.v2"

	"github.com/Lunchr/luncher-api/db"
	"github.com/Lunchr/luncher-api/db/model"
	"github.com/Lunchr/luncher-api/router"
	"github.com/deiwin/facebook"
)

type Post interface {
	UpdateForDate(model.DateWithoutTime, *model.User, *model.Restaurant) *router.HandlerError
	Update(*model.OfferGroupPost, *model.User, *model.Restaurant) *router.HandlerError
}

func NewPost(groupPosts db.OfferGroupPosts, offers db.Offers, regions db.Regions, fbAuth facebook.Authenticator) Post {
	return &facebookPost{
		groupPosts: groupPosts,
		offers:     offers,
		regions:    regions,
		fbAuth:     fbAuth,
	}
}

type facebookPost struct {
	groupPosts db.OfferGroupPosts
	offers     db.Offers
	regions    db.Regions
	fbAuth     facebook.Authenticator
}

func (f *facebookPost) UpdateForDate(date model.DateWithoutTime, user *model.User, restaurant *model.Restaurant) *router.HandlerError {
	post, err := f.groupPosts.GetByDate(date, restaurant.ID)
	if err == mgo.ErrNotFound {
		postToInsert := &model.OfferGroupPost{
			RestaurantID:    restaurant.ID,
			Date:            date,
			MessageTemplate: restaurant.DefaultGroupPostMessageTemplate,
		}
		insertedPosts, err := f.groupPosts.Insert(postToInsert)
		if err != nil {
			router.NewHandlerError(err, "Failed to create a group post with restaurant defaults", http.StatusInternalServerError)
		}
		post = insertedPosts[0]
	} else if err != nil {
		return router.NewHandlerError(err, "Failed to fetch a group post for that date", http.StatusInternalServerError)
	}
	return f.Update(post, user, restaurant)
}

func (f *facebookPost) Update(post *model.OfferGroupPost, user *model.User, restaurant *model.Restaurant) *router.HandlerError {
	if restaurant.FacebookPageID == "" {
		return nil
	}
	fbAPI := f.fbAuth.APIConnection(&user.Session.FacebookUserToken)
	// Remove the current post from FB, if it's already there
	if post.FBPostID != "" {
		err := fbAPI.PostDelete(user.Session.FacebookPageToken, post.FBPostID)
		if err != nil {
			return router.NewHandlerError(err, "Failed to delete the current post from Facebook", http.StatusBadGateway)
		}
	}
	offersForDate, handlerErr := f.getOffersForDate(post.Date, restaurant)
	if handlerErr != nil {
		return handlerErr
	} else if len(offersForDate) == 0 {
		return nil
	}
	message := f.formFBMessage(post, offersForDate)
	// Add the new version
	fbPost, err := fbAPI.PagePublish(user.Session.FacebookPageToken, restaurant.FacebookPageID, message)
	if err != nil {
		return router.NewHandlerError(err, "Failed to post the offer to Facebook", http.StatusBadGateway)
	}
	post.FBPostID = fbPost.ID
	if err = f.groupPosts.UpdateByID(post.ID, post); err != nil {
		return router.NewHandlerError(err, "Failed to update a group post in the DB", http.StatusInternalServerError)
	}
	return nil
}

func (f *facebookPost) formFBMessage(post *model.OfferGroupPost, offers []*model.Offer) string {
	offerMessages := make([]string, len(offers))
	for i, offer := range offers {
		offerMessages[i] = f.formFBOfferMessage(offer)
	}
	offersMessage := strings.Join(offerMessages, "\n")
	return fmt.Sprintf("%s\n\n%s", post.MessageTemplate, offersMessage)
}

func (f *facebookPost) formFBOfferMessage(o *model.Offer) string {
	// TODO get rid of the hard-coded €
	return fmt.Sprintf("%s - %.2f€", o.Title, o.Price)
}

func (f *facebookPost) getOffersForDate(date model.DateWithoutTime, restaurant *model.Restaurant) ([]*model.Offer, *router.HandlerError) {
	region, err := f.regions.GetName(restaurant.Region)
	if err != nil {
		return nil, router.NewHandlerError(err, "Failed to find the restaurant's region", http.StatusInternalServerError)
	}
	location, err := time.LoadLocation(region.Location)
	if err != nil {
		return nil, router.NewHandlerError(err, "Failed to load region's location", http.StatusInternalServerError)
	}
	startTime, endTime, err := date.TimeBounds(location)
	if err != nil {
		return nil, router.NewHandlerError(err, "Failed to parse a date", http.StatusInternalServerError)
	}
	offersForDate, err := f.offers.GetForRestaurantWithinTimeBounds(restaurant.ID, startTime, endTime)
	if err != nil {
		return nil, router.NewHandlerError(err, "Failed to find offers for this date", http.StatusInternalServerError)
	}
	return offersForDate, nil
}
