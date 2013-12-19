package wave

import (
	"fmt"
	"net/http"
)

// CustomersService handles communication with the customer related methods of the Wave API.
//
// Wave API docs: http://docs.waveapps.com/endpoints/customers.html
type CustomersService struct {
	client *Client
}

// ShippingDetails represents details for shipping for a given Customer.
type ShippingDetails struct {
	ShipToContact        *string  `json:"ship_to_contact,omitempty"`
	DeliveryInstructions *string  `json:"delivery_instructions,omitempty"`
	PhoneNumber          *string  `json:"phone_number,omitempty"`
	Address              *Address `json:"address,omitempty"`
}

// Customer represents an entity associated with an invoice or transaction.
type Customer struct {
	ID              uint64           `json:"id,omitempty"`
	URL             *string          `json:"url,omitempty"`
	AccountNumber   *string          `json:"account_number,omitempty"`
	FirstName       *string          `json:"first_name,omitempty"`
	LastName        *string          `json:"last_name,omitempty"`
	Email           *string          `json:"email,omitempty"`
	FaxNumber       *string          `json:"fax_number,omitempty"`
	MobileNumber    *string          `json:"mobile_number,omitempty"`
	PhoneNumber     *string          `json:"phone_number,omitempty"`
	TollFreeNumber  *string          `json:"toll_free_number,omitempty"`
	Website         *string          `json:"website,omitempty"`
	Currency        *Currency        `json:"currency,omitempty"`
	ShippingDetails *ShippingDetails `json:"shipping_details,omitempty"`
	Address         *Address         `json:address,omitempty"`
	DateCreated     *DateTime        `json:"date_created,omitempty"`
	DateModified    *DateTime        `json:"date_modified,omitempty"`
}

// FullName returns the full name of a customer.
//
// Given a first and last name, FullName will return 'First Last'.
// Given either a first or last name, FullName will return whichever is non-empty.
func (c *Customer) FullName() string {
	if c.FirstName == nil && c.LastName == nil {
		return ""
	}
	if c.FirstName == nil {
		return *c.LastName
	}
	if c.LastName == nil {
		return *c.FirstName
	}

	return fmt.Sprintf("%v %v", *c.FirstName, *c.LastName)
}

func (c *Customer) String() string {
	fullName := c.FullName()
	if c.Email == nil {
		return fullName
	}
	if fullName != "" {
		return fmt.Sprintf("%v (%v)", fullName, *c.Email)
	}
	return *c.Email
}

// List all customers for a given business.
//
// Wave API docs: http://docs.waveapps.com/endpoints/customers.html#get--businesses-{business_id}-customers-
func (service *CustomersService) List(businessID string) ([]Customer, *http.Response, error) {
	url := fmt.Sprintf("businesses/%v/customers", businessID)
	req, err := service.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	customers := new([]Customer)
	resp, err := service.client.Do(req, customers)
	if err != nil {
		return nil, resp, err
	}
	return *customers, resp, nil
}

// Get an existing customer for a given business.
//
// Wave API docs: http://docs.waveapps.com/endpoints/customers.html#get--businesses-{business_id}-customers-{customer_id}-
func (service *CustomersService) Get(businessID string, customerID uint64) (*Customer, *http.Response, error) {
	url := fmt.Sprintf("businesses/%v/customers/%v", businessID, customerID)
	req, err := service.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	customer := new(Customer)
	resp, err := service.client.Do(req, customer)
	if err != nil {
		return nil, resp, err
	}
	return customer, resp, nil
}

// Create a new customer for a given business.
//
// Wave API docs: http://docs.waveapps.com/endpoints/customers.html#post--businesses-{business_id}-customers-
func (service *CustomersService) Create(businessID string, customer Customer) (*Customer, *http.Response, error) {
	url := fmt.Sprintf("businesses/%v/customers", businessID)
	req, err := service.client.NewRequest("POST", url, customer)
	if err != nil {
		return nil, nil, err
	}
	c := new(Customer)
	resp, err := service.client.Do(req, c)
	if err != nil {
		return nil, resp, err
	}
	return c, resp, nil
}

// Replace an existing customer. You cannot create a customer using this method.
//
// Wave API docs: http://docs.waveapps.com/endpoints/customers.html#put--businesses-{business_id}-customers-{customer_id}-
func (service *CustomersService) Replace(businessID string, customerID uint64, customer Customer) (*Customer, *http.Response, error) {
	url := fmt.Sprintf("businesses/%v/customers/%v", businessID, customerID)
	req, err := service.client.NewRequest("PUT", url, customer)
	if err != nil {
		return nil, nil, err
	}
	c := new(Customer)
	resp, err := service.client.Do(req, c)
	if err != nil {
		return nil, resp, err
	}
	return c, resp, nil
}

// Update an existing customer. You cannot create a customer using this method.
//
// Wave API docs: http://docs.waveapps.com/endpoints/customers.html#patch--businesses-{business_id}-customers-{customer_id}-
func (service *CustomersService) Update(businessID string, customerID uint64, customer Customer) (*Customer, *http.Response, error) {
	url := fmt.Sprintf("businesses/%v/customers/%v", businessID, customerID)
	req, err := service.client.NewRequest("PATCH", url, customer)
	if err != nil {
		return nil, nil, err
	}
	c := new(Customer)
	resp, err := service.client.Do(req, c)
	if err != nil {
		return nil, resp, err
	}
	return c, resp, nil
}

// Delete an existing Customer.
//
// Wave API docs: http://docs.waveapps.com/endpoints/accounts.html#delete--businesses-{business_id}-accounts-{account_id}-
func (service *CustomersService) Delete(businessID string, customerID uint64) (*http.Response, error) {
	url := fmt.Sprintf("businesses/%v/customers/%v", businessID, customerID)
	req, err := service.client.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	return service.client.Do(req, nil)
}
