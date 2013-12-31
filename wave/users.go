// Copyright (c) 2013, Nick Presta
// All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wave

import "fmt"

// UsersService handles communication with the user related methods of the Wave API.
//
// Wave API docs: http://docs.waveapps.com/endpoints/users.html
type UsersService struct {
	client *Client
}

// User represents a Wave user.
type User struct {
	ID           *string   `json:"id,omitempty"`
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
	} `json:"emails,omitempty"`
	Profile struct {
		DateOfBirth *Date `json:"date_of_birth,omitempty"`
	} `json:"profile,omitempty"`
	Businesses []struct {
		ID  *string `json:"id,omitempty"`
		URL *string `json:"url,omitempty"`
	} `json:"businesses,omitempty"`
}

// FullName returns the full name of a customer.
//
// Given a first and last name, FullName will return 'First Last'.
// Given either a first or last name, FullName will return whichever is non-empty.
func (u User) FullName() string {
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

func (u User) String() string {
	return u.FullName()
}

// Get a specific user. Accepts "current" to refer to the current user.
//
// Wave API docs: http://docs.waveapps.com/endpoints/users.html#get--user-
func (service *UsersService) Get() (*User, *Response, error) {
	req, err := service.client.NewRequest("GET", "user", nil)
	if err != nil {
		return nil, nil, err
	}
	user := new(User)
	resp, err := service.client.Do(req, user, false)
	if err != nil {
		return nil, resp, err
	}
	return user, resp, nil
}

// Replace an existing user. You cannot create a user using this method.
//
// Wave API docs: http://docs.waveapps.com/endpoints/users.html#put--user-
func (service *UsersService) Replace(user *User) (*User, *Response, error) {
	req, err := service.client.NewRequest("PUT", "user", user)
	if err != nil {
		return nil, nil, err
	}
	u := new(User)
	resp, err := service.client.Do(req, u, false)
	if err != nil {
		return nil, resp, err
	}
	return u, resp, nil
}

// Update an existing user. You cannot create a user using this method.
//
// Wave API docs: http://docs.waveapps.com/endpoints/users.html#patch--user-
func (service *UsersService) Update(user *User) (*User, *Response, error) {
	req, err := service.client.NewRequest("PATCH", "user", user)
	if err != nil {
		return nil, nil, err
	}
	u := new(User)
	resp, err := service.client.Do(req, u, false)
	if err != nil {
		return nil, resp, err
	}
	return u, resp, nil
}
