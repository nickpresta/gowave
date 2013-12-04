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
	Id           *string   `json:"id,omitempty"`
	URL          *string   `json:"url,omitempty"`
	FirstName    *string   `json:"first_name,omitempty"`
	LastName     *string   `json:"first_name,omitempty"`
	DateCreated  *DateTime `json:"date_created,omitempty"`
	DateModified *DateTime `json:"date_modified,omitempty"`
	LastLogin    *DateTime `json:"last_login,omitempty"`
	Emails       []struct {
		Email      *string `json:"email,omitempty"`
		IsVerified *bool   `json:"is_verified,omitempty"`
		IsDefault  *bool   `json:"is_default,omitempty"`
	} `json:"emails"`
	Profile struct {
		DateOfBirth *Date `json:"date_of_birth,omitempty"`
	} `json:"profile"`
	Businesses []struct {
		Id  *string `json:"id,omitempty"`
		URL *string `json:"url,omitempty"`
	} `json:"businesses"`
}

func (u *User) FullName() string {
	if u.FirstName == nil && u.LastName == nil {
		return ""
	}
	if u.FirstName == nil {
		return *u.LastName
	}
	if u.LastName == nil {
		return *u.FirstName
	}

	return fmt.Sprintf("%v %v", *u.FirstName, *u.LastName)
}

func (u *User) String() string {
	return u.FullName()
}

// Get a specific user. Accepts "current" to refer to the current user.
//
// Wave API docs: http://waveaccounting.github.io/api/endpoints/users.html#get--users-(identity_user_id)-
func (service *UsersService) Get() (*User, *http.Response, error) {
	req, err := service.client.NewRequest("GET", "user", nil)
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
func (service *UsersService) Replace(user User) (*User, *http.Response, error) {
	req, err := service.client.NewRequest("PUT", "user", user)
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

// Update an existing user. You cannot create a user using this method.
//
// Wave API docs: http://waveaccounting.github.io/api/endpoints/users.html#patch--users-(identity_user_id)-
func (service *UsersService) Update(user User) (*User, *http.Response, error) {
	req, err := service.client.NewRequest("PATCH", "user", user)
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
