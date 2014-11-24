package db_test

import (
	"fmt"

	"github.com/deiwin/praad-api/db/model"
	"github.com/onsi/gomega/types"
)

func ContainOfferMock(expectedIndex interface{}) types.GomegaMatcher {
	return &mockOfferMatcher{
		expectedIndex: expectedIndex,
	}
}

type mockOfferMatcher struct {
	expectedIndex interface{}
}

func (matcher *mockOfferMatcher) Match(actual interface{}) (success bool, err error) {
	response, ok := actual.([]*model.Offer)
	if !ok {
		return false, fmt.Errorf("ContainOfferMock matcher expects an []*model.Offer for the actual value")
	}

	var expectedIndex int
	expectedIndex, ok = matcher.expectedIndex.(int)
	if !ok {
		return false, fmt.Errorf("ContainOfferMock matcher expects an int for the expected value")
	}
	expected := mocks.offers[expectedIndex]

	var contains = false
	for _, offer := range response {
		if offer.Title == expected.Title {
			contains = true
			break
		}
	}
	return contains, nil
}

func (matcher *mockOfferMatcher) FailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nto contain the mock offer with index %#v", actual, matcher.expectedIndex)
}

func (matcher *mockOfferMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	return fmt.Sprintf("Expected\n\t%#v\nnot to contain the mock offer with index %#v", actual, matcher.expectedIndex)
}
