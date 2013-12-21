package wave

import (
	"fmt"
	"net/http"
)

// ProductsService handles communication with the product related methods of the Wave API.
//
// Wave API docs: http://docs.waveapps.com/endpoints/products.html
type ProductsService struct {
	client *Client
}

// Product represents an entity associated with an invoice or transaction.
type Product struct {
	ID             *uint64   `json:"id,omitempty"`
	URL            *string   `json:"url,omitempty"`
	Name           *string   `json:"name,omitempty"`
	Price          *float64  `json:"price,omitempty"`
	Description    *string   `json:"description,omitempty"`
	IsSold         *bool     `json:"is_sold,omitempty"`
	IsBought       *bool     `json:"is_bought,omitempty"`
	IncomeAccount  *Account  `json:"income_account,omitempty"`
	ExpenseAccount *Account  `json:"expense_account,omitempty"`
	DateCreated    *DateTime `json:"date_created,omitempty"`
	DateModified   *DateTime `json:"date_modified,omitempty"`
}

func (p *Product) String() string {
	return *p.Name
}

// List all products for a given business.
//
// Wave API docs: http://docs.waveapps.com/endpoints/products.html#get--businesses-{business_id}-products-
func (service *ProductsService) List(businessID string, embed *EmbedArgs) ([]Product, *http.Response, error) {
	url := fmt.Sprintf("businesses/%v/products", businessID)
	if embed != nil {
		queryParams := embed.BuildQueryString()
		if len(queryParams) > 1 {
			url = fmt.Sprintf("%v?%v", url, queryParams)
		}
	}
	req, err := service.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	products := new([]Product)
	resp, err := service.client.Do(req, products)
	if err != nil {
		return nil, resp, err
	}
	return *products, resp, nil
}

// Get an existing product for a given business.
//
// Wave API docs: http://docs.waveapps.com/endpoints/products.html#get--businesses-{business_id}-products-{product_id}-
func (service *ProductsService) Get(businessID string, productID uint64, embed EmbedArgs) (*Product, *http.Response, error) {
	url := fmt.Sprintf("businesses/%v/products/%v", businessID, productID)
	if embed != nil {
		queryParams := embed.BuildQueryString()
		if len(queryParams) > 1 {
			url = fmt.Sprintf("%v?%v", url, queryParams)
		}
	}
	req, err := service.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	product := new(Product)
	resp, err := service.client.Do(req, product)
	if err != nil {
		return nil, resp, err
	}
	return product, resp, nil
}

// Create a new product for a given business.
//
// Wave API docs: http://docs.waveapps.com/endpoints/products.html#post--businesses-{business_id}-products-
func (service *ProductsService) Create(businessID string, product Product) (*Product, *http.Response, error) {
	url := fmt.Sprintf("businesses/%v/products", businessID)
	req, err := service.client.NewRequest("POST", url, product)
	if err != nil {
		return nil, nil, err
	}
	p := new(Product)
	resp, err := service.client.Do(req, p)
	if err != nil {
		return nil, resp, err
	}
	return p, resp, nil
}

// Replace an existing product. You cannot create a product using this method.
//
// Wave API docs: http://docs.waveapps.com/endpoints/products.html#put--businesses-{business_id}-products-{product_id}-
func (service *ProductsService) Replace(businessID string, productID uint64, product Product) (*Product, *http.Response, error) {
	url := fmt.Sprintf("businesses/%v/products/%v", businessID, productID)
	req, err := service.client.NewRequest("PUT", url, product)
	if err != nil {
		return nil, nil, err
	}
	p := new(Product)
	resp, err := service.client.Do(req, p)
	if err != nil {
		return nil, resp, err
	}
	return p, resp, nil
}

// Update an existing product. You cannot create a product using this method.
//
// Wave API docs: http://docs.waveapps.com/endpoints/products.html#patch--businesses-{business_id}-products-{product_id}-
func (service *ProductsService) Update(businessID string, productID uint64, product Product) (*Product, *http.Response, error) {
	url := fmt.Sprintf("businesses/%v/products/%v", businessID, productID)
	req, err := service.client.NewRequest("PATCH", url, product)
	if err != nil {
		return nil, nil, err
	}
	p := new(Product)
	resp, err := service.client.Do(req, p)
	if err != nil {
		return nil, resp, err
	}
	return p, resp, nil
}

// Delete an existing Product.
//
// Wave API docs: http://docs.waveapps.com/endpoints/accounts.html#delete--businesses-{business_id}-accounts-{account_id}-
func (service *ProductsService) Delete(businessID string, productID uint64) (*http.Response, error) {
	url := fmt.Sprintf("businesses/%v/products/%v", businessID, productID)
	req, err := service.client.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	return service.client.Do(req, nil)
}
