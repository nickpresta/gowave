package wave

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)

const (
	version        = "0.0.1"
	defaultBaseURL = "https://api.waveapps.com/"
	userAgent      = "gowave/" + version
)

type Timestamp time.Time

type Client struct {
	// HTTP client used to communicate with the API.
	client *http.Client

	// Base URL for API requests. Defaults to the Wave API.
	// BaseURL should always be specified with a trailing slash.
	BaseURL *url.URL

	// User agent used when communicating with the Wave API.
	UserAgent string

	// Services used to communicate with different parts of the Wave API
	Businesses *BusinessesService
	Currencies *CurrenciesService
	Countries  *CountriesService
}

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

func (t *Timestamp) UnmarshalJSON(b []byte) error {
	v, err := time.Parse("2006-01-02T15:04:05", string(b[1:len(b)-1]))
	if err != nil {
		return err
	}
	*t = Timestamp(v)
	return nil
}

func (t Timestamp) MarshalJSON() ([]byte, error) {
	trueTime := time.Time(t)
	if y := trueTime.Year(); y < 0 || y >= 10000 {
		return nil, errors.New("Time.MarshalJSON: year outside of range [0,9999]")
	}
	return []byte(trueTime.Format(`"2006-01-02T15:04:05"`)), nil
}

func NewClient(client *http.Client) *Client {
	if client == nil {
		client = http.DefaultClient
	}

	baseURL, _ := url.Parse(defaultBaseURL)

	c := &Client{client: client, BaseURL: baseURL, UserAgent: userAgent}
	c.Businesses = &BusinessesService{client: c}
	c.Currencies = &CurrenciesService{client: c}
	c.Countries = &CountriesService{client: c}

	return c
}

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

	req.Header.Add("User-Agent", c.UserAgent)
	return req, nil
}

func (c *Client) Do(request *http.Request, v interface{}) (*http.Response, error) {
	resp, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if err = CheckResponse(resp); err != nil {
		return resp, err
	}

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(v)
	}
	return resp, err
}

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
