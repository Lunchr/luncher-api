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

const (
	// postUpdateDebounceDuration is the period we leave for the users to modify a day's offers in before
	// publishing the offers to Facebook. That is, the last modification time to today's offers + this time
	// will be when the post goes live. NB: May not be less than 10 minutes, as per Facebook's documentation.
	// That is, Facebook doesn't allow posts to be scheduled less than 10 minutes (or more than 6 months, but
	// I don't think that will be an issue) from now.
	postUpdateDebounceDuration       = 11 * time.Minute
	publishDurationBeforeOfferActive = 30 * time.Minute
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
	fbPost, handlerErr := formFBPost(post, offersForDate)
	if handlerErr != nil {
		return handlerErr
	}

	fbAPI := f.fbAuth.APIConnection(&user.Session.FacebookUserToken)
	fbPostResponse, err := fbAPI.PagePublish(user.Session.FacebookPageToken, restaurant.FacebookPageID, fbPost)
	if err != nil {
		return router.NewHandlerError(err, "Failed to post the offers to Facebook", http.StatusBadGateway)
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
	fbPost, handlerErr := formFBPostForUpdate(post, offersForDate, user, fbAPI)
	if handlerErr != nil {
		return handlerErr
	}

	err := fbAPI.PostUpdate(user.Session.FacebookPageToken, post.FBPostID, fbPost)
	if err != nil {
		return router.NewHandlerError(err, "Failed updating the offers in Facebook", http.StatusBadGateway)
	}
	return nil
}

func formFBPostForUpdate(post *model.OfferGroupPost, offersForDate []*model.Offer, user *model.User, fbAPI facebook.API) (*fbmodel.Post, *router.HandlerError) {
	currentPost, err := fbAPI.Post(user.Session.FacebookPageToken, post.FBPostID)
	if err != nil {
		return nil, router.NewHandlerError(err, "Failed to retrieve the current post from FB", http.StatusBadGateway)
	}
	if currentPost.IsPublished {
		return &fbmodel.Post{
			Message:   formFBMessage(post, offersForDate),
			Published: true,
		}, nil
	}
	return formFBPost(post, offersForDate)
}

func formFBPost(post *model.OfferGroupPost, offersForDate []*model.Offer) (*fbmodel.Post, *router.HandlerError) {
	publishTime, handlerErr := calculatePublishTime(offersForDate)
	if handlerErr != nil {
		return nil, handlerErr
	}
	return &fbmodel.Post{
		Message:              formFBMessage(post, offersForDate),
		ScheduledPublishTime: publishTime,
		Published:            publishTime.Before(time.Now()),
	}, nil
}

func formFBMessage(post *model.OfferGroupPost, offers []*model.Offer) string {
	offerMessages := make([]string, len(offers))
	for i, offer := range offers {
		offerMessages[i] = formFBOfferMessage(offer)
	}
	offersMessage := strings.Join(offerMessages, "\n")
	return fmt.Sprintf("%s\n\n%s", post.MessageTemplate, offersMessage)
}

func formFBOfferMessage(o *model.Offer) string {
	// TODO get rid of the hard-coded €
	return fmt.Sprintf("%s - %.2f€", o.Title, o.Price)
}

// calculatePublishTime returns a time either 5 minutes (the debounce period) from now or the earliest FromTime of an offer,
// whichever is later.
func calculatePublishTime(offers []*model.Offer) (time.Time, *router.HandlerError) {
	if len(offers) == 0 {
		return time.Time{}, router.NewSimpleHandlerError("Cannot calculate a publish time for 0 offers", http.StatusInternalServerError)
	}
	var earliestTime time.Time
	for i, offer := range offers {
		if i == 0 {
			earliestTime = offer.FromTime
			continue
		}
		if offer.FromTime.Before(earliestTime) {
			earliestTime = offer.FromTime
		}
	}
	earliestPublishTime := earliestTime.Add(-publishDurationBeforeOfferActive)
	afterDebounceBuffer := time.Now().Add(postUpdateDebounceDuration)
	if earliestPublishTime.Before(afterDebounceBuffer) {
		return afterDebounceBuffer, nil
	}
	return earliestPublishTime, nil
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
