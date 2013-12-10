package wave

import (
	"fmt"
	"net/http"
)

// AccountsService handles communication with the acccounts related methods of the Wave API.
//
// Wave API docs: http://waveaccounting.github.io/api/endpoints/accounts.html
type AccountsService struct {
	client *Client
}

// Account represents a Wave business.
type Account struct {
	ID                    uint64    `json:"id,omitempty"`
	URL                   *string   `json:"url,omitempty"`
	Name                  *string   `json:"name,omitempty"`
	Active                *bool     `json:"active,omitempty"`
	AccountType           *string   `json:"account_type,omitempty"`
	AccountClass          *string   `json:"account_class,omitempty"`
	StandardAccountNumber *int      `json:"standard_account_number,omitempty"`
	AccountNumber         *int      `json:"account_number,omitempty"`
	IsPayment             *bool     `json:"is_payment,omitempty"`
	CanDelete             *bool     `json:"can_delete,omitempty"`
	IsCurrencyEditable    *bool     `json:"is_currency_editable,omitempty"`
	IsNameEditable        *bool     `json:"is_name_editable,omitempty"`
	IsPaymentEditable     *bool     `json:"is_payment_editable,omitempty"`
	DateCreated           *DateTime `json:"date_created,omitempty"`
	DateModified          *DateTime `json:"date_modified,omitempty"`
	Currency              *Currency `json:"currency,omitempty"`
}

func (a *Account) String() string {
	return fmt.Sprintf("%v (type=%v, payment=%v)", *a.Name, *a.AccountType, *a.IsPayment)
}

// List all accounts for a given business.
//
// Wave API docs: http://docs.waveapps.com/endpoints/accounts.html#get--businesses-{business_id}-accounts-
func (service *AccountsService) List(businessID string) ([]Account, *http.Response, error) {
	url := fmt.Sprintf("businesses/%v/accounts", businessID)
	req, err := service.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	accounts := new([]Account)
	resp, err := service.client.Do(req, accounts)
	if err != nil {
		return nil, resp, err
	}
	return *accounts, resp, nil
}

// Get an existing account for a given business.
//
// Wave API docs: http://docs.waveapps.com/endpoints/accounts.html#get--businesses-{business_id}-accounts-{account_id}-
func (service *AccountsService) Get(businessID string, accountID uint64) (*Account, *http.Response, error) {
	url := fmt.Sprintf("businesses/%v/accounts/%v", businessID, accountID)
	req, err := service.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	account := new(Account)
	resp, err := service.client.Do(req, account)
	if err != nil {
		return nil, resp, err
	}
	return account, resp, nil
}

// Create a new account according to a standard account template.
//
// Wave API docs: http://waveaccounting.github.io/api/endpoints/businesses.html#post--businesses-
func (service *AccountsService) Create(businessID string, account Account) (*Account, *http.Response, error) {
	url := fmt.Sprintf("businesses/%v/accounts", businessID)
	req, err := service.client.NewRequest("POST", url, account)
	if err != nil {
		return nil, nil, err
	}
	a := new(Account)
	resp, err := service.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}
	return a, resp, nil
}

// Replace an existing account. You cannot create an account using this method.
//
// Wave API docs: http://docs.waveapps.com/endpoints/accounts.html#put--businesses-{business_id}-accounts-{account_id}-
func (service *AccountsService) Replace(businessID string, accountID uint64, account Account) (*Account, *http.Response, error) {
	url := fmt.Sprintf("businesses/%v/accounts/%v", businessID, accountID)
	req, err := service.client.NewRequest("PUT", url, account)
	if err != nil {
		return nil, nil, err
	}
	a := new(Account)
	resp, err := service.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}
	return a, resp, nil
}

// Update an existing account. You cannot create an account using this method.
//
// Wave API docs: http://docs.waveapps.com/endpoints/accounts.html#patch--businesses-{business_id}-accounts-{account_id}-
func (service *AccountsService) Update(businessID string, accountID uint64, account Account) (*Account, *http.Response, error) {
	url := fmt.Sprintf("businesses/%v/accounts/%v", businessID, accountID)
	req, err := service.client.NewRequest("PATCH", url, account)
	if err != nil {
		return nil, nil, err
	}
	a := new(Account)
	resp, err := service.client.Do(req, a)
	if err != nil {
		return nil, resp, err
	}
	return a, resp, nil
}

// Delete an existing account. The `can_delete` attribute of an account determines if it can be deleted.
//
// Wave API docs: http://docs.waveapps.com/endpoints/accounts.html#delete--businesses-{business_id}-accounts-{account_id}-
func (service *AccountsService) Delete(businessID string, accountID uint64) (*http.Response, error) {
	url := fmt.Sprintf("businesses/%v/accounts/%v", businessID, accountID)
	req, err := service.client.NewRequest("DELETE", url, nil)
	if err != nil {
		return nil, err
	}
	return service.client.Do(req, nil)
}
