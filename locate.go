package main

import (
	"encoding/json"
	"fmt"
	"maps"
	"net/http"

	"github.com/vuhnger/tracemap/constants"
)

type Location struct {
	City      string  `json:"city"`
	Country   string  `json:"country"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type Locations map[string]*Location

type ipmapResponse struct {
	Location *struct {
		CityName          string  `json:"cityName"`
		CountryCodeAlpha2 string  `json:"countryCodeAlpha2"`
		Latitude          float64 `json:"latitude"`
		Longitude         float64 `json:"longitude"`
	} `json:"location"`
}

var ipmapClient = &http.Client{Timeout: constants.IPMapTimeout}

func locateIP(baseURL, ip string) (*Location, error) {
	response, err := ipmapClient.Get(fmt.Sprintf("%s/%s/best", baseURL, ip))
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode == http.StatusBadRequest {
		return nil, nil
	}
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("IPmap answered %d for %s", response.StatusCode, ip)
	}
	var body ipmapResponse
	if err := json.NewDecoder(response.Body).Decode(&body); err != nil {
		return nil, err
	}
	if body.Location == nil {
		return nil, nil
	}
	return &Location{
		City:      body.Location.CityName,
		Country:   body.Location.CountryCodeAlpha2,
		Latitude:  body.Location.Latitude,
		Longitude: body.Location.Longitude,
	}, nil
}

func locate(baseURL string, ips []string, known Locations) (Locations, error) {
	locations := maps.Clone(known)
	if locations == nil {
		locations = Locations{}
	}
	for _, ip := range ips {
		if _, ok := locations[ip]; ok {
			continue
		}
		location, err := locateIP(baseURL, ip)
		if err != nil {
			return nil, err
		}
		locations[ip] = location
	}
	return locations, nil
}
