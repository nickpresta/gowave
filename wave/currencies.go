// Copyright (c) 2013, Nick Presta
// All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wave

import "fmt"

// CurrenciesService handles communication with the currency related methods of the Wave API.
//
// Wave API docs: http://docs.waveapps.com/endpoints/currencies.html
type CurrenciesService struct {
	client *Client
}

// Currency represents a currency in ISO 4217 format (http://en.wikipedia.org/wiki/ISO_4217).
type Currency struct {
	URL    *string `json:"url,omitempty"`
	Code   *string `json:"code,omitempty"`
	Symbol *string `json:"symbol,omitempty"`
	Name   *string `json:"name,omitempty"`
}

func (c *Currency) String() string {
	return fmt.Sprintf("%v (%v)", *c.Code, *c.Name)
}

// List all currencies available in Wave.
//
// Wave API docs: http://docs.waveapps.com/endpoints/currencies.html#get--currencies-
func (service *CurrenciesService) List() ([]Currency, *Response, error) {
	req, err := service.client.NewRequest("GET", "currencies", nil)
	if err != nil {
		return nil, nil, err
	}
	currencies := new([]Currency)
	resp, err := service.client.Do(req, currencies, false)
	if err != nil {
		return nil, resp, err
	}
	return *currencies, resp, nil
}

// Get a specific currency.
//
// Wave API docs: http://docs.waveapps.com/endpoints/currencies.html#get--currencies-(code)-
func (service *CurrenciesService) Get(code string) (*Currency, *Response, error) {
	u := fmt.Sprintf("currencies/%v", code)
	req, err := service.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	currency := new(Currency)
	resp, err := service.client.Do(req, currency, false)
	if err != nil {
		return nil, resp, err
	}
	return currency, resp, nil
}
