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
	fbmodel "github.com/deiwin/facebook/model"
)

type Post interface {
	Update(model.DateWithoutTime, *model.User, *model.Restaurant) *router.HandlerError
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

func (f *facebookPost) Update(date model.DateWithoutTime, user *model.User, restaurant *model.Restaurant) *router.HandlerError {
	if restaurant.FacebookPageID == "" {
		return nil
	}
	post, err := f.groupPosts.GetByDate(date, restaurant.ID)
	if err == mgo.ErrNotFound {
		postToInsert := &model.OfferGroupPost{
			RestaurantID:    restaurant.ID,
			Date:            date,
			MessageTemplate: restaurant.DefaultGroupPostMessageTemplate,
		}
		insertedPosts, err := f.groupPosts.Insert(postToInsert)
		if err != nil {
			return router.NewHandlerError(err, "Failed to create a group post with restaurant defaults", http.StatusInternalServerError)
		}
		post = insertedPosts[0]
	} else if err != nil {
		return router.NewHandlerError(err, "Failed to fetch a group post for that date", http.StatusInternalServerError)
	}
	return f.updatePost(post, user, restaurant)
}

func (f *facebookPost) updatePost(post *model.OfferGroupPost, user *model.User, restaurant *model.Restaurant) *router.HandlerError {
	if restaurant.FacebookPageID == "" {
		return nil
	}

	offersForDate, handlerErr := f.getOffersForDate(post.Date, restaurant)
	if handlerErr != nil {
		return handlerErr
	}
	if post.FBPostID == "" {
		if len(offersForDate) == 0 {
			return nil
		}
		return f.publishNewPost(post, offersForDate, user, restaurant)
	}
	if len(offersForDate) == 0 {
		return f.deleteExistingPost(post, user, restaurant)
	}
	return f.updateExistingPost(post, offersForDate, user, restaurant)
}

func (f *facebookPost) publishNewPost(post *model.OfferGroupPost, offersForDate []*model.Offer, user *model.User, restaurant *model.Restaurant) *router.HandlerError {
	fbAPI := f.fbAuth.APIConnection(&user.Session.FacebookUserToken)
	fbPost := f.formFBPost(post, offersForDate)
	fbPostResponse, err := fbAPI.PagePublish(user.Session.FacebookPageToken, restaurant.FacebookPageID, fbPost)
	if err != nil {
		return router.NewHandlerError(err, "Failed to post the offer to Facebook", http.StatusBadGateway)
	}
	post.FBPostID = fbPostResponse.ID
	if err = f.groupPosts.UpdateByID(post.ID, post); err != nil {
		return router.NewHandlerError(err, "Failed to update a group post in the DB", http.StatusInternalServerError)
	}
	return nil
}

func (f *facebookPost) deleteExistingPost(post *model.OfferGroupPost, user *model.User, restaurant *model.Restaurant) *router.HandlerError {
	fbAPI := f.fbAuth.APIConnection(&user.Session.FacebookUserToken)
	err := fbAPI.PostDelete(user.Session.FacebookPageToken, post.FBPostID)
	if err != nil {
		return router.NewHandlerError(err, "Failed to delete the current post from Facebook", http.StatusBadGateway)
	}
	post.FBPostID = ""
	if err = f.groupPosts.UpdateByID(post.ID, post); err != nil {
		return router.NewHandlerError(err, "Failed to update a group post in the DB", http.StatusInternalServerError)
	}
	return nil
}

func (f *facebookPost) updateExistingPost(post *model.OfferGroupPost, offersForDate []*model.Offer, user *model.User, restaurant *model.Restaurant) *router.HandlerError {
	fbAPI := f.fbAuth.APIConnection(&user.Session.FacebookUserToken)
	fbPost := f.formFBPost(post, offersForDate)
	err := fbAPI.PostUpdate(user.Session.FacebookPageToken, post.FBPostID, fbPost)
	if err != nil {
		return router.NewHandlerError(err, "Failed to post the offer to Facebook", http.StatusBadGateway)
	}
	return nil
}

func (f *facebookPost) formFBPost(post *model.OfferGroupPost, offers []*model.Offer) *fbmodel.Post {
	return &fbmodel.Post{
		Message: f.formFBMessage(post, offers),
	}
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
