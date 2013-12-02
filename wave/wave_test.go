package wave

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the Wave client being tested.
	client *Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

type data struct {
	I int
}

func setUp() {
	// test server
	mux = http.NewServeMux()
	server = httptest.NewServer(mux)

	// wave client configured to use test server
	client = NewClient(nil)
	url, _ := url.Parse(server.URL)
	client.BaseURL = url
}

func tearDown() {
	server.Close()
}

func TestErrorResponseFormat(t *testing.T) {
	Convey("ErrorResponse should output with a specific format", t, func() {
		url, _ := url.Parse("http://example.com/")
		errResp := &ErrorResponse{
			Response: &http.Response{
				Request:    &http.Request{Method: "GET", URL: url},
				StatusCode: http.StatusNotFound,
			},
			Err: struct {
				Message string `json:"message"`
			}{Message: "Error message"},
		}
		So(errResp.Error(), ShouldEqual, "GET http://example.com/: 404 Error message")
	})
}

func TestDateTimeUnmarshalJSON(t *testing.T) {
	Convey("Unmarshalling JSON to a DateTime", t, func() {
		timestamp := DateTime(time.Now())

		Convey("Should unmarshal valid JSON", func() {
			err := timestamp.UnmarshalJSON([]byte(`"2013-08-19T09:18:32"`))
			So(err, ShouldEqual, nil)

			v := time.Time(timestamp)
			So(v.Year(), ShouldEqual, 2013)
			So(v.Month().String(), ShouldEqual, "August")
			So(v.Day(), ShouldEqual, 19)
			So(v.Hour(), ShouldEqual, 9)
			So(v.Minute(), ShouldEqual, 18)
			So(v.Second(), ShouldEqual, 32)
		})

		Convey("Unmarshalling JSON to DateTime should return an error", func() {
			err := timestamp.UnmarshalJSON([]byte(`invalid`))
			So(err, ShouldNotEqual, nil)
		})
	})
}

func TestDateTimeMarshalJSON(t *testing.T) {
	Convey("Marshalling a DateTime to JSON", t, func() {
		Convey("Should marshal a valid DateTime", func() {
			timestamp := DateTime(time.Date(2009, time.November, 10, 23, 4, 20, 0, time.UTC))
			json, err := timestamp.MarshalJSON()

			So(err, ShouldEqual, nil)
			So(string(json), ShouldEqual, `"2009-11-10T23:04:20"`)
		})

		Convey("Marshalling an invalid DateTime should return an error", func() {
			timestamp := DateTime(time.Date(-5, time.November, 10, 23, 4, 20, 0, time.UTC))
			_, err := timestamp.MarshalJSON()
			So(err, ShouldNotEqual, nil)
		})
	})
}

func TestNewClientHasDefaultClient(t *testing.T) {
	Convey("If no client is passed, the default http client is used", t, func() {
		client := NewClient(nil)
		So(client.client, ShouldEqual, http.DefaultClient)
	})
}

func TestNewClientDefaults(t *testing.T) {
	Convey("NewClient should return proper defaults", t, func() {
		c := NewClient(nil)

		Convey("The BaseURL should have a default", func() {
			So(c.BaseURL.String(), ShouldEqual, defaultBaseURL)
		})

		Convey("The UserAgent should have a default", func() {
			So(c.UserAgent, ShouldEqual, userAgent)
		})
	})
}

func TestNewRequest(t *testing.T) {
	Convey("Making a NewRequest should set up the correct data", t, func() {
		c := NewClient(nil)

		inURL, outURL := "foo", defaultBaseURL+"foo"
		inBody, outBody := &data{1}, `{"I":1}`+"\n"
		req, err := c.NewRequest("GET", inURL, inBody)

		Convey("Error should be nil", func() {
			So(err, ShouldEqual, nil)
		})

		Convey("Request URL should be correct", func() {
			So(req.URL.String(), ShouldEqual, outURL)
		})

		Convey("Request User-Agent should be correct", func() {
			So(req.Header.Get("User-Agent"), ShouldEqual, userAgent)
		})

		Convey("Request Body should be JSON serialized", func() {
			body, err := ioutil.ReadAll(req.Body)
			So(err, ShouldEqual, nil)
			So(string(body), ShouldEqual, outBody)
		})
	})

	Convey("Making a NewRequest with invalid data should return an error", t, func() {
		c := NewClient(nil)

		Convey("Passing an invalid URL", func() {
			_, err := c.NewRequest("GET", "%gh&%ij", nil)
			So(err, ShouldNotEqual, nil)
		})

		Convey("Passing invalid body", func() {
			_, err := c.NewRequest("GET", "foo", func() {})
			So(err, ShouldNotEqual, nil)
		})
	})
}

func TestDo(t *testing.T) {
	Convey("Making a request with a good response", t, func() {
		setUp()
		defer tearDown()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			So("GET", ShouldEqual, r.Method)
			fmt.Fprint(w, `{"I":1}`)
		})
		req, _ := client.NewRequest("GET", "/", nil)
		body := new(data)
		_, err := client.Do(req, body)
		So(err, ShouldEqual, nil)
		So(body, ShouldResemble, &data{I: 1})
	})
	Convey("Making a request with a bad response", t, func() {
		setUp()
		defer tearDown()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, "Not Found", http.StatusNotFound)
		})
		req, _ := client.NewRequest("GET", "/", nil)
		resp, err := client.Do(req, nil)
		So(err, ShouldNotEqual, nil)
		So(resp.StatusCode, ShouldEqual, http.StatusNotFound)
	})
	Convey("Passing in a bad request", t, func() {
		_, err := client.Do(&http.Request{}, nil)
		So(err, ShouldNotEqual, nil)
	})
}

func TestCheckResponse(t *testing.T) {
	Convey("Checking the response of a request", t, func() {
		Convey("With a bad status code", func() {
			resp := &http.Response{
				StatusCode: http.StatusNotFound,
				Body:       ioutil.NopCloser(strings.NewReader(`{"error": {"message":"Error"}}`)),
			}
			err := CheckResponse(resp).(*ErrorResponse)
			So(err, ShouldNotEqual, nil)
			So(err.Response.StatusCode, ShouldEqual, http.StatusNotFound)
			So(err.Err.Message, ShouldEqual, "Error")
		})
		Convey("With a good status code", func() {
			resp := &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader(`{}`)),
			}
			err := CheckResponse(resp)
			So(err, ShouldEqual, nil)
		})
	})
}
