// Copyright (c) 2013, Nick Presta
// All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wave

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"strconv"
	"time"

	"github.com/google/go-querystring/query"
)

const (
	version        = "0.0.1"
	defaultBaseURL = "https://api.waveapps.com/"
)

var userAgent = fmt.Sprintf("gowave/%v (Go %v; %v/%v)", version, runtime.Version(), runtime.GOOS, runtime.GOARCH)

// DateTime represents a time that can be unmarshalled from a JSON string,
// formatted as ISO-8601 (2006-01-02T15:04:05+00:00)
type DateTime time.Time

// Date represents a time that can be unmarshalled from a JSON string,
// formatted as ISO-8601 (2006-01-02)
type Date time.Time

// Client used to interact with the Wave API.
type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests. Defaults to the Wave API.
	// BaseURL should always be specified with a trailing slash.
	BaseURL *url.URL

	// User agent used when communicating with the Wave API.
	UserAgent string

	// Services used to communicate with different parts of the Wave API
	Accounts   *AccountsService
	Businesses *BusinessesService
	Countries  *CountriesService
	Currencies *CurrenciesService
	Customers  *CustomersService
	Products   *ProductsService
	Users      *UsersService
}

// PageOptions specifies the pagination options for methods that support pagination (mostly LIST and GET options)
type PageOptions struct {
	Page     int `url:"page,omitempty"`
	PageSize int `url:"page_size,omitempty"`
}

// addOptions adds the parameters in opt as URL query parameters to s. opt
// must be a struct whose fields may contain "url" tags.
func addOptions(s string, opt interface{}) (string, error) {
	v := reflect.ValueOf(opt)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return s, nil
	}

	u, err := url.Parse(s)
	if err != nil {
		return s, err
	}

	qs, err := query.Values(opt)
	if err != nil {
		return s, err
	}

	u.RawQuery = qs.Encode()
	return u.String(), nil
}

type paginatedResponse struct {
	Next       *string `json:"next,omitempty"`
	Previous   *string `json:"previous,omitempty"`
	TotalCount int     `json:"total_count,omitempty"`
}

// Response is a Wave API response. This embeds the standard http.Response and provides pagination information.
type Response struct {
	CurrentPage  int
	NextPage     int
	PreviousPage int
	TotalCount   int
	*http.Response
}

// ErrorResponse represents a single error returned from the API.
type ErrorResponse struct {
	Response *http.Response
	Err      struct {
		Message string `json:"message"`
	} `json:"error"`
}

func (resp *ErrorResponse) Error() string {
	return fmt.Sprintf("%v %v: %d %+v",
		resp.Response.Request.Method, resp.Response.Request.URL, resp.Response.StatusCode, resp.Err.Message)
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// DateTime is expected to be in ISO-8601 format.
func (t *DateTime) UnmarshalJSON(b []byte) error {
	v, err := time.Parse("2006-01-02T15:04:05+00:00", string(b[1:len(b)-1]))
	if err != nil {
		return err
	}
	*t = DateTime(v)
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
// DateTime will be formatted as ISO-8601.
func (t DateTime) MarshalJSON() ([]byte, error) {
	trueTime := time.Time(t)
	if y := trueTime.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("year outside of range [0,9999]")
	}
	return []byte(trueTime.Format(`"2006-01-02T15:04:05+00:00"`)), nil
}

// UnmarshalJSON implements the json.Unmarshaler interface.
// Date is expected to be in ISO-8601 format.
func (d *Date) UnmarshalJSON(b []byte) error {
	v, err := time.Parse("2006-01-02", string(b[1:len(b)-1]))
	if err != nil {
		return err
	}
	*d = Date(v)
	return nil
}

// MarshalJSON implements the json.Marshaler interface.
// Date will be formatted as ISO-8601.
func (d Date) MarshalJSON() ([]byte, error) {
	trueTime := time.Time(d)
	if y := trueTime.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("year outside of range [0,9999]")
	}
	return []byte(trueTime.Format(`"2006-01-02"`)), nil
}

// NewClient returns a new Wave API client.
// If a nil httpClient is provided, http.DefaultClient will be used.
// To use API methods which require
// authentication, provide an http.Client that will perform the authentication
// for you (such as that provided by the goauth2 library).
func NewClient(client *http.Client) *Client {
	if client == nil {
		client = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: client, BaseURL: baseURL, UserAgent: userAgent}
	c.Accounts = &AccountsService{client: c}
	c.Businesses = &BusinessesService{client: c}
	c.Countries = &CountriesService{client: c}
	c.Currencies = &CurrenciesService{client: c}
	c.Customers = &CustomersService{client: c}
	c.Products = &ProductsService{client: c}
	c.Users = &UsersService{client: c}

	return c
}

// NewRequest creates an API request. A relative URL can be provided in urlStr,
// in which case it is resolved relative to the BaseURL of the Client.
// Relative URLs should always be specified without a preceding slash.
// If specified, the value pointed to by body is JSON encoded and included as the
// request body.
func (c *Client) NewRequest(method string, urlStr string, body interface{}) (*http.Request, error) {
	rel, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	u := c.BaseURL.ResolveReference(rel)

	buf := new(bytes.Buffer)
	if body != nil {
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequest(method, u.String(), buf)
	if err != nil {
		return nil, err
	}

	req.Header.Set("User-Agent", c.UserAgent)
	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

// Do sends an API request and returns the API response.
// The API response is decoded and stored in the value pointed to by v, or returned
// as an error if an API error has occured.
func (c *Client) Do(request *http.Request, v interface{}, isPaginated bool) (*Response, error) {
	resp, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	response := &Response{Response: resp}

	if err = CheckResponse(response.Response); err != nil {
		return response, err
	}

	var buf bytes.Buffer
	buf.ReadFrom(resp.Body)
	bufBytes := buf.Bytes()
	reader := bytes.NewReader(bufBytes)
	// Put back the body into response.Body so it can be ready again by the consumer
	response.Response.Body = ioutil.NopCloser(reader)

	if v != nil {
		err = json.Unmarshal(bufBytes, &v)
	}

	if err == nil && isPaginated {
		// parse out the pagination information
		p := new(paginatedResponse)
		json.Unmarshal(bufBytes, &p)
		err = response.populatePageValues(p)
	}
	return response, err
}

func (r *Response) populatePageValues(opts *paginatedResponse) error {
	r.TotalCount = opts.TotalCount

	if opts.Next != nil {
		u, err := url.Parse(*opts.Next)
		if err != nil {
			return err
		}
		r.NextPage, _ = strconv.Atoi(u.Query().Get("page"))
	}

	if opts.Previous != nil {
		u, err := url.Parse(*opts.Previous)
		if err != nil {
			return err
		}
		r.PreviousPage, _ = strconv.Atoi(u.Query().Get("page"))
	}

	if opts.Next == nil && opts.Previous == nil {
		r.CurrentPage = 1
	} else if opts.Previous != nil {
		r.CurrentPage = r.PreviousPage + 1
	} else if opts.Next != nil {
		r.CurrentPage = r.NextPage - 1
	}

	return nil
}

// CheckResponse checks the API response for errors, and returns them if
// present. A response is considered an error if it has a status code outside
// the 200 range. API error responses are expected to have either no response
// body, or a JSON response body that maps to ErrorResponse. Any other
// response body will be silently ignored.
func CheckResponse(resp *http.Response) error {
	if code := resp.StatusCode; 200 <= code && code <= 299 {
		return nil
	}
	errorResponse := &ErrorResponse{Response: resp}
	data, err := ioutil.ReadAll(resp.Body)
	if err == nil && data != nil {
		json.Unmarshal(data, errorResponse)
	}
	return errorResponse
}

// Bool is a helper method that allocates a new bool value and returns a pointer to it.
func Bool(v bool) *bool {
	p := new(bool)
	*p = v
	return p
}

// Int is a helper method that allocates a new int value and returns a pointer to it.
func Int(v int) *int {
	p := new(int)
	*p = v
	return p
}

// Float64 is a helper method that allocates a new float64 value and returns a pointer to it.
func Float64(v float64) *float64 {
	p := new(float64)
	*p = v
	return p
}

// String is a helper method that allocates a new string value and returns a pointer to it.
func String(v string) *string {
	p := new(string)
	*p = v
	return p
}
