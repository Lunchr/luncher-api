package geo

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strings"
)

const apiEndpoint = "https://maps.googleapis.com/maps/api/geocode/json"

// Coder is an object that knows how to geocode an address
type Coder interface {
	Code(address string) (Location, error)
	CodeForRegion(address, region string) (Location, error)
}

type coder struct {
	conf *Config
}

func NewCoder(conf *Config) Coder {
	return coder{conf}
}

func (c coder) Code(address string) (Location, error) {
	parameters := c.paramsForAddress(address)
	url := urlFor(parameters)
	return fetch(url)
}

func (c coder) CodeForRegion(address, region string) (Location, error) {
	parameters := c.paramsForAddress(address)
	parameters.Add("region", escape(region))
	url := urlFor(parameters)
	return fetch(url)
}

// Location defines a geographical coordinate
type Location struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type response struct {
	Status       string   `json:"status"`
	ErrorMessage string   `json:"error_message"`
	Results      []result `json:"results"`
}

type result struct {
	Geometry     geometry `json:"geometry"`
	PartialMatch bool     `json:"partial_match"`
}

type geometry struct {
	Location Location `json:"location"`
}

func (c coder) paramsForAddress(address string) url.Values {
	return url.Values{
		"address": {escape(address)},
		"key":     {c.conf.APIKey},
	}
}

func escape(s string) string {
	return url.QueryEscape(strings.TrimSpace(s))
}

func urlFor(params url.Values) string {
	return fmt.Sprintf("%s?%s", apiEndpoint, params.Encode())
}

func fetch(url string) (Location, error) {
	var nilLocation Location
	httpResponse, err := http.Get(url)
	if err != nil {
		return nilLocation, err
	}
	defer httpResponse.Body.Close()

	var response = new(response)
	err = json.NewDecoder(httpResponse.Body).Decode(response)
	if err != nil {
		return nilLocation, err
	}

	if response.Status != "OK" {
		if response.ErrorMessage != "" {
			return nilLocation, fmt.Errorf("Geocoder service error!  (%s - %s)", response.Status, response.ErrorMessage)
		}
		return nilLocation, fmt.Errorf("Geocoder service error!  (%s)", response.Status)
	}
	if len(response.Results) > 1 {
		return nilLocation, errors.New("More than one response received from the Geocoder service")
	}
	result := response.Results[0]
	if result.PartialMatch {
		return nilLocation, errors.New("Geocoder returned a partial match. Check for typos")
	}

	return result.Geometry.Location, nil
}
