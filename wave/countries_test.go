package wave

import (
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

const (
	countryJSON = `{"name":"Canada","country_code":"CA","currency_code":"CAD","url":"https://api.example.com/countries/CA/","provinces":[{"name":"Alberta","slug":"alberta"}]}`
	countriesJSON = "["+countryJSON+"]"
	provincesJSON = `[{"name":"Ontario","slug":"ontario"}]`
)

func countryStruct() *Country {
	var c Country
	json.Unmarshal([]byte(countryJSON), &c)
	return &c
}

func provincesStruct() []Province {
	p := new([]Province)
	json.Unmarshal([]byte(provincesJSON), &p)
	return *p
}

func TestCountriesService(t *testing.T) {
	Convey("LISTing all countries", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/countries", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "GET")
			fmt.Fprint(w, countriesJSON)
		})

		countries, _, err := client.Countries.List()
		c := []Country{*countryStruct()}
		So(err, ShouldEqual, nil)
		So(countries, ShouldResemble, c)
	})

	Convey("GET a specific country", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/countries/CA", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "GET")
			fmt.Fprint(w, countryJSON)
		})

		country, _, err := client.Countries.Get("CA")
		So(err, ShouldEqual, nil)
		So(country, ShouldResemble, countryStruct())
	})

	Convey("LISTing all provinces for a given country", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/countries/CA/provinces", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "GET")
			fmt.Fprint(w, provincesJSON)
		})

		provinces, _, err := client.Countries.Provinces("CA")
		p := provincesStruct()
		So(err, ShouldEqual, nil)
		So(provinces, ShouldResemble, p)
	})

	Convey("String method on Country", t, func() {
		c := new(Country)
		c.Name = "Canada"
		c.CountryCode = "CA"
		So(c.String(), ShouldEqual, "Canada (CA)")
	})

	Convey("String method on Province", t, func() {
		p := new(Province)
		p.Name = "Ontario"
		So(p.String(), ShouldEqual, "Ontario")
	})
}
