// Copyright (c) 2013, Nick Presta
// All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package wave

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	expectedProvincesJSON = `[{"name":"Ontario","slug":"ontario"}]`
	expectedCountryJSON   = `{
		"name":"Canada",
		"country_code":"CA",
		"currency_code":"CAD",
		"provinces":` + expectedProvincesJSON + `,
		"url":"url"
	}`
	expectedCountriesJSON = "[" + expectedCountryJSON + "]"
)

func countryStructHelper() *Country {
	var c Country
	json.Unmarshal([]byte(expectedCountryJSON), &c)
	return &c
}

func provincesStruct() []Province {
	p := new([]Province)
	json.Unmarshal([]byte(expectedProvincesJSON), &p)
	return *p
}

func TestCountriesService(t *testing.T) {
	countryStruct := countryStructHelper()

	Convey("Testing JSON marshalling of a Country", t, func() {
		p := []Province{
			Province{
				Name: String("Ontario"),
				Slug: String("ontario"),
			},
		}
		c := &Country{
			Name:         String("Canada"),
			CountryCode:  String("CA"),
			CurrencyCode: String("CAD"),
			URL:          String("url"),
			Provinces:    p,
		}
		checkMarshalJSON(c, expectedCountryJSON)
	})

	Convey("LISTing all countries", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/countries/", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "GET")
			fmt.Fprint(w, expectedCountriesJSON)
		})

		countries, _, err := client.Countries.List()
		c := []Country{*countryStruct}
		So(err, ShouldBeNil)
		So(countries, ShouldResemble, c)
	})

	Convey("GET a specific country", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/countries/CA/", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "GET")
			fmt.Fprint(w, expectedCountryJSON)
		})

		country, _, err := client.Countries.Get("CA")
		So(err, ShouldBeNil)
		So(country, ShouldResemble, countryStruct)
	})

	Convey("GET a specific country with an invalid code", t, func() {
		country, resp, err := client.Countries.Get("%")
		checkInvalidURLError(country, resp, err)
	})

	Convey("String method on Country", t, func() {
		c := new(Country)
		c.Name = String("Canada")
		c.CountryCode = String("CA")
		So(c.String(), ShouldEqual, "Canada (CA)")
	})

	Convey("String method on Province", t, func() {
		p := new(Province)
		p.Name = String("Ontario")
		So(p.String(), ShouldEqual, "Ontario")
	})
}
