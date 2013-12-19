package wave

import (
	"fmt"
	"net/http"
)

// CountriesService handles communication with the country related methods of the Wave API.
//
// Wave API docs: http://docs.waveapps.com/endpoints/geography.html
type CountriesService struct {
	client *Client
}

// Province represents a province for a given country.
type Province struct {
	Name *string `json:"name"`
	Slug *string `json:"slug"`
}

func (p *Province) String() string {
	return *p.Name
}

// Country represents a country in ISO 3166-1 alpha-2 format (http://en.wikipedia.org/wiki/ISO_3166-1_alpha-2).
type Country struct {
	Name         *string    `json:"name,omitempty"`
	CountryCode  *string    `json:"country_code,omitempty"`
	CurrencyCode *string    `json:"currency_code,omitempty"`
	Provinces    []Province `json:"provinces,omitempty"`
	URL          *string    `json:"url,omitempty"`
}

func (c *Country) String() string {
	return fmt.Sprintf("%v (%v)", *c.Name, *c.CountryCode)
}

// List all countries available in Wave.
//
// Wave API docs: http://docs.waveapps.com/endpoints/geography.html#get--countries-
func (service *CountriesService) List() ([]Country, *http.Response, error) {
	req, err := service.client.NewRequest("GET", "countries", nil)
	if err != nil {
		return nil, nil, err
	}
	countries := new([]Country)
	resp, err := service.client.Do(req, countries)
	if err != nil {
		return nil, resp, err
	}
	return *countries, resp, nil
}

// Get a specific currency.
//
// Wave API docs: http://docs.waveapps.com/endpoints/geography.html#get--countries-(country_code)-
func (service *CountriesService) Get(code string) (*Country, *http.Response, error) {
	u := fmt.Sprintf("countries/%v", code)
	req, err := service.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	country := new(Country)
	resp, err := service.client.Do(req, country)
	if err != nil {
		return nil, resp, err
	}
	return country, resp, nil
}

// Provinces gets a specific province for a given country.
//
// Wave API docs: http://docs.waveapps.com/endpoints/geography.html#get--countries-(country_code)-provinces-
func (service *CountriesService) Provinces(code string) ([]Province, *http.Response, error) {
	u := fmt.Sprintf("countries/%v/provinces", code)
	req, err := service.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	provinces := new([]Province)
	resp, err := service.client.Do(req, provinces)
	if err != nil {
		return nil, resp, err
	}
	return *provinces, resp, nil
}
