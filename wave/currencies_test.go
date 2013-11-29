package wave

import (
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func currencyJSON() string {
	return `{"url":"https://api.example.com/currencies/CAD/","code":"CAD","symbol":"$","name":"Canadian dollar"}`
}

func currencyStruct() *Currency {
	var c Currency
	json.Unmarshal([]byte(currencyJSON()), &c)
	return &c
}

func TestCurrenciesService(t *testing.T) {
	Convey("LISTing all currencies", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/currencies", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "GET")
			fmt.Fprint(w, "["+currencyJSON()+"]")
		})

		currencies, _, err := client.Currencies.List()
		c := []Currency{*currencyStruct()}
		So(err, ShouldEqual, nil)
		So(currencies, ShouldResemble, c)
	})

	Convey("GET a specific currency", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/currencies/1", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "GET")
			fmt.Fprint(w, currencyJSON())
		})

		currency, _, err := client.Currencies.Get("1")
		So(err, ShouldEqual, nil)
		So(currency, ShouldResemble, currencyStruct())
	})

	Convey("String method on Currency", t, func() {
		c := new(Currency)
		c.Name = "Canadian dollar"
		c.Code = "CAD"
		So(c.String(), ShouldEqual, "CAD (Canadian dollar)")
	})

}
