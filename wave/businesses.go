package wave

import (
	"fmt"
	"net/http"
)

type BusinessesService struct {
	client *Client
}

type Business struct {
	Id                  string    `json:"id"`
	URL                 string    `json:"url"`
	CompanyName         string    `json:"company_name"`
	PrimaryCurrencyCode string    `json:"primary_currency_code"`
	PhoneNumber         *string   `json:"phone_number"`
	MobilePhoneNumber   *string   `json:"mobile_phone_number"`
	TollFreePhoneNumber *string   `json:"toll_free_phone_number"`
	FaxNumber           *string   `json:"fax_number"`
	Website             *string   `json:"website"`
	IsPersonalBusiness  bool      `json:"is_personal_business"`
	DateCreated         Timestamp `json:"date_created"`
	dateCreated         string    `json:"date_created"`
	DateModified        Timestamp `json:"date_modified"`
	dateModified        string    `json:"date_modified"`
	BusinessTypeInfo    struct {
		BusinessType       *string `json:"business_type"`
		BusinessSubtype    *string `json:"business_subtype"`
		OrganizationalType *string `json:"organizational_type"`
	} `json:"business_type_info"`
	AddressInfo struct {
		Address1   *string `json:"address1"`
		Address2   *string `json:"address2"`
		City       *string `json:"city"`
		PostalCode *string `json:"postal_code"`
		Province   struct {
			Name *string `json:"name"`
			Slug *string `json:"slug"`
		} `json:"province"`
		Country struct {
			Name         *string `json:"name"`
			CountryCode  *string `json:"country_code"`
			CurrencyCode *string `json:"currency_code"`
			URL          string  `json:"url"`
		} `json:"country"`
	} `json:"address_info"`
}

func (service *BusinessesService) List() ([]Business, *http.Response, error) {
	req, err := service.client.NewRequest("GET", "businesses", nil)
	if err != nil {
		return nil, nil, err
	}
	businesses := new([]Business)
	resp, err := service.client.Do(req, businesses)
	if err != nil {
		return nil, resp, err
	}
	return *businesses, resp, nil
}

func (service *BusinessesService) Get(id string) (*Business, *http.Response, error) {
	u := fmt.Sprintf("businesses/%s", id)
	req, err := service.client.NewRequest("GET", u, nil)
	if err != nil {
		return nil, nil, err
	}
	business := new(Business)
	resp, err := service.client.Do(req, business)
	if err != nil {
		return nil, resp, err
	}
	return business, resp, nil
}
