package wave

import (
	"fmt"
	"net/http"
)

// CurrenciesService handles communication with the currency related methods of the Wave API.
//
// Wave API docs: http://waveaccounting.github.io/api/endpoints/currencies.html
type CurrenciesService struct {
	client *Client
}

// Currency represents a currency in ISO 4217 format (http://en.wikipedia.org/wiki/ISO_4217).
type Currency struct {
	URL    string `json:"url,omitempty"`
	Code   string `json:"code,omitempty"`
	Symbol string `json:"symbol,omitempty"`
	Name   string `json:"name,omitempty"`
}

func (c *Currency) String() string {
	return fmt.Sprintf("%v (%v)", c.Code, c.Name)
}

// List all currencies available in Wave.
//
// Wave API docs: http://waveaccounting.github.io/api/endpoints/currencies.html#get--currencies-
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

// Get a specific currency.
//
// Wave API docs: http://waveaccounting.github.io/api/endpoints/currencies.html#get--currencies-(code)-
func (service *CurrenciesService) Get(code string) (*Currency, *http.Response, error) {
	u := fmt.Sprintf("currencies/%v", code)
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
