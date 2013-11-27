package wave

import (
	"fmt"
	"net/http"
)

type CurrenciesService struct {
	client *Client
}

type Currency struct {
	URL    string `json:"url"`
	Code   string `json:"code"`
	Symbol string `json:"symbol"`
	Name   string `json:"name"`
}

func (service *CurrenciesService) List() ([]Currency, *http.Response, error) {
	req, err := service.client.NewRequest("GET", "currencies", nil)
	if err != nil {
		return nil, nil, err
	}
	currencies := new([]Currency)
	resp, err := service.client.Do(req, currencies)
	if err != nil {
		return nil, resp, err
	}
	return *currencies, resp, nil
}

func (service *CurrenciesService) Get(code string) (*Currency, *http.Response, error) {
	u := fmt.Sprintf("currencies/%s", code)
	req, err := service.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	currency := new(Currency)
	resp, err := service.client.Do(req, currency)
	if err != nil {
		return nil, resp, err
	}
	return currency, resp, nil
}
