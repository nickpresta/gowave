package wave

import (
	"fmt"
	"net/http"
)

type CountriesService struct {
	client *Client
}

type Province struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type Country struct {
	URL    string `json:"url"`
	CountryCode   string `json:"country_code"`
	CurrencyCode string `json:"currency_code"`
	Provinces []Province `json:"provinces"`
}

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

func (service *CountriesService) Get(code string) (*Country, *http.Response, error) {
	u := fmt.Sprintf("countries/%s", code)
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

func (service *CountriesService) Provinces(code string) ([]Province, *http.Response, error) {
	u := fmt.Sprintf("countries/%s/provinces", code)
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
