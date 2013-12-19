package wave

import (
	"fmt"
	"net/http"
)

// BusinessesService handles communication with the business related methods of the Wave API.
//
// Wave API docs: http://docs.waveapps.com/endpoints/businesses.html
type BusinessesService struct {
	client *Client
}

// BusinessTypeInfo represents the type of Wave business.
type BusinessTypeInfo struct {
}

// Business represents a Wave business.
type Business struct {
	ID                  *string   `json:"id,omitempty"`
	URL                 *string   `json:"url,omitempty"`
	CompanyName         *string   `json:"company_name,omitempty"`
	PrimaryCurrencyCode *string   `json:"primary_currency_code,omitempty"`
	BusinessType        *string   `json:"business_type,omitempty"`
	BusinessSubtype     *string   `json:"business_subtype,omitempty"`
	OrganizationalType  *string   `json:"organizational_type,omitempty"`
	Address1            *string   `json:"address1,omitempty"`
	Address2            *string   `json:"address2,omitempty"`
	City                *string   `json:"city,omitempty"`
	Province            *Province `json:"province,omitempty"`
	Country             *Country  `json:"country,omitempty"`
	PostalCode          *string   `json:"postal_code,omitempty"`
	PhoneNumber         *string   `json:"phone_number,omitempty"`
	MobilePhoneNumber   *string   `json:"mobile_phone_number,omitempty"`
	TollFreePhoneNumber *string   `json:"toll_free_phone_number,omitempty"`
	FaxNumber           *string   `json:"fax_number,omitempty"`
	Website             *string   `json:"website,omitempty"`
	IsPersonalBusiness  *bool     `json:"is_personal_business,omitempty"`
	DateCreated         *DateTime `json:"date_created,omitempty"`
	DateModified        *DateTime `json:"date_modified,omitempty"`
}

func (b *Business) String() string {
	return fmt.Sprintf("%v (id=%v, personal=%v)", *b.CompanyName, *b.ID, *b.IsPersonalBusiness)
}

// List all businesses owned by the authenticated user.
//
// Wave API docs: http://docs.waveapps.com/endpoints/businesses.html#get--businesses-
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

// Get an existing business.
//
// Wave API docs: http://docs.waveapps.com/endpoints/businesses.html#get--businesses-(identity_business_id)-
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

// Create a new business
//
// Wave API docs: http://docs.waveapps.com/endpoints/businesses.html#post--businesses-
func (service *BusinessesService) Create(business Business) (*Business, *http.Response, error) {
	req, err := service.client.NewRequest("POST", "businesses", business)
	if err != nil {
		return nil, nil, err
	}
	b := new(Business)
	resp, err := service.client.Do(req, b)
	if err != nil {
		return nil, resp, err
	}
	return b, resp, nil
}

// Replace an existing business. You cannot create a business using this method.
//
// Wave API docs: http://docs.waveapps.com/endpoints/businesses.html#put--businesses-(identity_business_id)-
func (service *BusinessesService) Replace(id string, business Business) (*Business, *http.Response, error) {
	url := fmt.Sprintf("businesses/%v", id)
	req, err := service.client.NewRequest("PUT", url, business)
	if err != nil {
		return nil, nil, err
	}
	b := new(Business)
	resp, err := service.client.Do(req, b)
	if err != nil {
		return nil, resp, err
	}
	return b, resp, nil
}

// Update an existing business. You cannot create a business using this method.
//
// Wave API docs: http://docs.waveapps.com/endpoints/businesses.html#patch--businesses-(identity_business_id)-
func (service *BusinessesService) Update(id string, business Business) (*Business, *http.Response, error) {
	url := fmt.Sprintf("businesses/%v", id)
	req, err := service.client.NewRequest("PATCH", url, business)
	if err != nil {
		return nil, nil, err
	}
	b := new(Business)
	resp, err := service.client.Do(req, b)
	if err != nil {
		return nil, resp, err
	}
	return b, resp, nil
}
