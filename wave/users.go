package wave

import (
	"fmt"
	"net/http"
)

// UsersService handles communication with the user related methods of the Wave API.
//
// Wave API docs: http://waveaccounting.github.io/api/endpoints/users.html
type UsersService struct {
	client *Client
}

// User represents a Wave user.
type User struct {
	Id string `json:"id,omitempty"`
	URL string `json:"url,omitempty"`
	Email string `json:"email,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName string `json:"first_name,omitempty"`
	DateCreated DateTime `json:"date_created,omitempty"`
	LastLogin DateTime `json:"last_login,omitempty"`
	Emails []struct {
		Email string `json:"email,omitempty"`
		IsVerified bool `json:"is_verified,omitempty"`
		IsDefault bool `json:"is_default,omitempty"`
	} `json:"emails,omitempty"`
	Profile struct {
		DateOfBirth Date `json:"date_of_birth,omitempty"`
		EmailPreferences struct {
			SpecialOffers bool `json:"special_offers,omitempty"`
		}
	} `json:"profile,omitempty"`
	Businesses []struct {
		Id string `json:"id,omitempty"`
		URL string `json:"url,omitempty"`
	} `json:"businesses,omitempty"`
}

func (u *User) FullName() string {
	if u.FirstName == "" {
		return u.LastName
	}
	if u.LastName == "" {
		return u.FirstName
	}
	if u.FirstName == "" && u.LastName == "" {
		return ""
	}
	return fmt.Sprintf("%v %v", u.FirstName, u.LastName)
}

func (u *User) String() string {
	return fmt.Sprintf("%v (%v)", u.FullName, u.Email)
}

// Get a specific user. Accepts "current" to refer to the current user.
//
// Wave API docs: http://waveaccounting.github.io/api/endpoints/users.html#get--users-(identity_user_id)-
func (service *UsersService) Get(id string) (*User, *http.Response, error) {
	url := fmt.Sprintf("users/%s", id)
	req, err := service.client.NewRequest("GET", url, nil)
	if err != nil {
		return nil, nil, err
	}
	user := new(User)
	resp, err := service.client.Do(req, user)
	if err != nil {
		return nil, resp, err
	}
	return user, resp, nil
}

// Replace an existing user. You cannot create a user using this method.
//
// Wave API docs: http://waveaccounting.github.io/api/endpoints/users.html#put--users-(identity_user_id)-
func (service *UsersService) Replace(user User, id string) (*User, *http.Response, error) {
	url := fmt.Sprintf("users/%s", id)
	req, err := service.client.NewRequest("PUT", url, user)
	if err != nil {
		return nil, nil, err
	}
	u := new(User)
	resp, err := service.client.Do(req, u)
	if err != nil {
		return nil, resp, err
	}
	return u, resp, nil
}
