package wave

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	expectedCurrencyJSON = `{
		"url":"https://api.example.com/currencies/CAD/",
		"code":"CAD",
		"symbol":"$",
		"name":"Canadian dollar"
	}`
	expectedCurrenciesJSON = "[" + expectedCurrencyJSON + "]"
)

func TestCurrenciesService(t *testing.T) {
	expectedCurrencyStruct := new(Currency)
	json.Unmarshal([]byte(expectedCurrencyJSON), expectedCurrencyStruct)

	Convey("LISTing all currencies", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/currencies", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "GET")
			fmt.Fprint(w, expectedCurrenciesJSON)
		})

		currencies, _, err := client.Currencies.List()
		c := []Currency{*expectedCurrencyStruct}
		So(err, ShouldBeNil)
		So(currencies, ShouldResemble, c)
	})

	Convey("GET a specific currency", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/currencies/1", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "GET")
			fmt.Fprint(w, expectedCurrencyJSON)
		})

		currency, _, err := client.Currencies.Get("1")
		So(err, ShouldBeNil)
		So(currency, ShouldResemble, expectedCurrencyStruct)
	})

	Convey("GET a specific currency with an invalid ID", t, func() {
		currency, resp, err := client.Currencies.Get("%")
		checkInvalidURLError(currency, resp, err)
	})

	Convey("String method on Currency", t, func() {
		c := new(Currency)
		c.Name = "Canadian dollar"
		c.Code = "CAD"
		So(c.String(), ShouldEqual, "CAD (Canadian dollar)")
	})

}
