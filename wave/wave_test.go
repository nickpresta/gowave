package wave

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"reflect"
	"runtime"
	"strings"
	"testing"
	"time"

	"code.google.com/p/goauth2/oauth"
	. "github.com/smartystreets/goconvey/convey"
)

var (
	// mux is the HTTP request multiplexer used with the test server.
	mux *http.ServeMux

	// client is the Wave client being tested.
	client *Client

	// integrationClient is the client used to actually hit the Wave API
	integrationClient *Client

	// server is a test HTTP server used to provide mock API responses.
	server *httptest.Server
)

type CachedContent struct {
	Response *http.Response
	Error    error
}

type CachedResponseTransport struct {
	Transport http.RoundTripper
	cache     map[string]*CachedContent
}

func NewCachedResponseTransport() *CachedResponseTransport {
	return &CachedResponseTransport{cache: make(map[string]*CachedContent)}
}

func (t *CachedResponseTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	u, err := url.Parse(r.URL.String())
	if err != nil {
		return nil, err
	}
	key := fmt.Sprintf("%v %v", r.Method, u.String())

	cachedContent, inCache := t.cache[key]
	if inCache {
		log.Printf("Found '%v' in cache\n", key)
		return cachedContent.Response, cachedContent.Error
	}

	log.Printf("Not in cache; making round trip for '%v'\n", key)
	resp, err := t.Transport.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	if code := resp.StatusCode; code >= 400 || code < 300 {
		log.Printf("Storing in cache: '%v'\n", key)
		t.cache[key] = &CachedContent{Response: resp, Error: err}
	}

	return resp, err
}

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

func setUpIntegrations() {
	if os.Getenv("WAVE_API_ACCESS_TOKEN") == "" {
		panic("You must provide a WAVE_API_ACCESS_TOKEN environment variable for integration tests")
	}

	transport := &oauth.Transport{
		Config:    &oauth.Config{},
		Transport: http.DefaultTransport,
		Token:     &oauth.Token{AccessToken: os.Getenv("WAVE_API_ACCESS_TOKEN")},
	}
	cachedTransport := NewCachedResponseTransport()
	cachedTransport.Transport = transport
	integrationClient = NewClient(&http.Client{Transport: cachedTransport})
}

func tearDown() {
	server.Close()
}

func checkInvalidURLError(v interface{}, resp *http.Response, err error) {
	So(v, ShouldBeNil)
	So(resp, ShouldBeNil)
	So(err, ShouldNotBeNil)

	typeErr, ok := err.(*url.Error)
	So(ok, ShouldBeTrue)
	So(typeErr.Op, ShouldEqual, "parse")
}

func checkMarshalJSON(v interface{}, e string) {
	want, err := json.Marshal(v)
	So(err, ShouldBeNil)

	given := new(bytes.Buffer)
	err = json.Compact(given, []byte(e))
	So(err, ShouldBeNil)

	So(given.String(), ShouldEqual, string(want))
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
			err := timestamp.UnmarshalJSON([]byte(`"2013-08-19T09:18:32+00:00"`))
			So(err, ShouldBeNil)

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
			So(err, ShouldNotBeNil)
		})
	})
}

func TestDateMarshalJSON(t *testing.T) {
	Convey("Marshalling a Date to JSON", t, func() {
		Convey("Should marshal a valid Date", func() {
			date := Date(time.Date(2009, time.November, 10, 0, 0, 0, 0, time.UTC))
			json, err := date.MarshalJSON()

			So(err, ShouldBeNil)
			So(string(json), ShouldEqual, `"2009-11-10"`)
		})

		Convey("Marshalling an invalid Date should return an error", func() {
			date := Date(time.Date(-5, time.November, 10, 0, 0, 0, 0, time.UTC))
			_, err := date.MarshalJSON()
			So(err, ShouldNotBeNil)
		})
	})
}

func TestDateUnmarshalJSON(t *testing.T) {
	Convey("Unmarshalling JSON to a Date", t, func() {
		date := Date(time.Now())

		Convey("Should unmarshal valid JSON", func() {
			err := date.UnmarshalJSON([]byte(`"2013-08-19"`))
			So(err, ShouldBeNil)

			v := time.Time(date)
			So(v.Year(), ShouldEqual, 2013)
			So(v.Month().String(), ShouldEqual, "August")
			So(v.Day(), ShouldEqual, 19)
		})

		Convey("Unmarshalling JSON to Date should return an error", func() {
			err := date.UnmarshalJSON([]byte(`invalid`))
			So(err, ShouldNotBeNil)
		})
	})
}

func TestDateTimeMarshalJSON(t *testing.T) {
	Convey("Marshalling a DateTime to JSON", t, func() {
		Convey("Should marshal a valid DateTime", func() {
			timestamp := DateTime(time.Date(2009, time.November, 10, 23, 4, 20, 0, time.UTC))
			json, err := timestamp.MarshalJSON()

			So(err, ShouldBeNil)
			So(string(json), ShouldEqual, `"2009-11-10T23:04:20+00:00"`)
		})

		Convey("Marshalling an invalid DateTime should return an error", func() {
			timestamp := DateTime(time.Date(-5, time.November, 10, 23, 4, 20, 0, time.UTC))
			_, err := timestamp.MarshalJSON()
			So(err, ShouldNotBeNil)
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
			So(c.UserAgent, ShouldEqual, strings.Replace(userAgent, "$VERSION$", runtime.Version(), 1))
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
			So(err, ShouldBeNil)
		})

		Convey("Request URL should be correct", func() {
			So(req.URL.String(), ShouldEqual, outURL)
		})

		Convey("Request User-Agent should be correct", func() {
			So(req.Header.Get("User-Agent"), ShouldEqual, strings.Replace(userAgent, "$VERSION$", runtime.Version(), 1))
		})

		Convey("Request Body should be JSON serialized", func() {
			body, err := ioutil.ReadAll(req.Body)
			So(err, ShouldBeNil)
			So(string(body), ShouldEqual, outBody)
		})
	})

	Convey("Making a NewRequest with invalid data should return an error", t, func() {
		c := NewClient(nil)

		Convey("Passing an invalid URL", func() {
			_, err := c.NewRequest("GET", "%gh&%ij", nil)
			So(err, ShouldNotBeNil)
		})

		Convey("Passing invalid body", func() {
			_, err := c.NewRequest("GET", "foo", func() {})
			So(err, ShouldNotBeNil)
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
		So(err, ShouldBeNil)
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
		So(err, ShouldNotBeNil)
		So(resp.StatusCode, ShouldEqual, http.StatusNotFound)
	})
	Convey("Passing in a bad request", t, func() {
		_, err := client.Do(&http.Request{}, nil)
		So(err, ShouldNotBeNil)
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
			So(err, ShouldNotBeNil)
			So(err.Response.StatusCode, ShouldEqual, http.StatusNotFound)
			So(err.Err.Message, ShouldEqual, "Error")
		})
		Convey("With a good status code", func() {
			resp := &http.Response{
				StatusCode: 200,
				Body:       ioutil.NopCloser(strings.NewReader(`{}`)),
			}
			err := CheckResponse(resp)
			So(err, ShouldBeNil)
		})
	})
}

func TestEmbedArgs(t *testing.T) {
	Convey("EmbedArgs should convert to a querystring", t, func() {
		e := EmbedArgs{"foo": true, "bar": false}
		queryParams := e.BuildQueryString()
		So(queryParams, ShouldEqual, "bar=false&foo=true")

		Convey("Should return an empty string with args", func() {
			e := EmbedArgs{}
			queryParams := e.BuildQueryString()
			So(queryParams, ShouldBeBlank)
		})
	})
}

func TestTypeHelpers(t *testing.T) {
	Convey("String should return a pointer to a string", t, func() {
		v := reflect.ValueOf(String("example"))
		So(v.Kind(), ShouldEqual, reflect.Ptr)
		So(v.IsNil(), ShouldBeFalse)
		So(reflect.Indirect(v).String(), ShouldEqual, "example")

		v = reflect.ValueOf(String(""))
		So(v.Kind(), ShouldEqual, reflect.Ptr)
		So(v.IsNil(), ShouldBeFalse)
		So(reflect.Indirect(v).String(), ShouldEqual, "")
	})

	Convey("Bool should return a pointer to a bool", t, func() {
		v := reflect.ValueOf(Bool(true))
		So(v.Kind(), ShouldEqual, reflect.Ptr)
		So(v.IsNil(), ShouldBeFalse)
		So(reflect.Indirect(v).Bool(), ShouldEqual, true)

		v = reflect.ValueOf(Bool(false))
		So(v.Kind(), ShouldEqual, reflect.Ptr)
		So(v.IsNil(), ShouldBeFalse)
		So(reflect.Indirect(v).Bool(), ShouldEqual, false)
	})

	Convey("Int should return a pointer to an int", t, func() {
		v := reflect.ValueOf(Int(1))
		So(v.Kind(), ShouldEqual, reflect.Ptr)
		So(v.IsNil(), ShouldBeFalse)
		So(reflect.Indirect(v).Int(), ShouldEqual, 1)

		v = reflect.ValueOf(Int(0))
		So(v.Kind(), ShouldEqual, reflect.Ptr)
		So(v.IsNil(), ShouldBeFalse)
		So(reflect.Indirect(v).Int(), ShouldEqual, 0)

		v = reflect.ValueOf(Int(-1))
		So(v.Kind(), ShouldEqual, reflect.Ptr)
		So(v.IsNil(), ShouldBeFalse)
		So(reflect.Indirect(v).Int(), ShouldEqual, -1)

		v = reflect.ValueOf(Int(42))
		So(v.Kind(), ShouldEqual, reflect.Ptr)
		So(v.IsNil(), ShouldBeFalse)
		So(reflect.Indirect(v).Int(), ShouldEqual, 42)
	})
}
