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
	"time"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	expectedBusinessJSON = `{
		"id":"id",
		"url":"url",
		"company_name":"Company Name",
		"primary_currency_code":"CAD",
		"business_type":"business type",
		"business_subtype":"business sub-type",
		"organizational_type":"organizational type",
		"address1":"Address 1",
		"address2":"Address 2",
		"city":"City",
		"province":{
			"name":"Ontario",
			"slug":"ontario"
		},
		"country":{
			"name":"Canada",
			"country_code":"CA",
			"currency_code":"CAD",
			"url":"url"
		},
		"postal_code":"A1B 2C3",
		"phone_number":"416-555-5555",
		"mobile_phone_number":"416-555-5556",
		"toll_free_phone_number":"416-555-5557",
		"fax_number":"416-555-5558",
		"website":"https://example.com",
		"is_personal_business":false,
		"date_created":"2009-11-10T23:00:00+00:00",
		"date_modified":"2009-11-10T23:00:00+00:00"
	}`
	expectedBusinessesJSON = "[" + expectedBusinessJSON + "]"
)

func TestBusinessesService(t *testing.T) {
	expectedBusinessStruct := new(Business)
	json.Unmarshal([]byte(expectedBusinessJSON), &expectedBusinessStruct)

	Convey("Testing JSON marshalling of a Business", t, func() {
		datetime := DateTime(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC))
		b := &Business{
			ID:                  String("id"),
			URL:                 String("url"),
			CompanyName:         String("Company Name"),
			PrimaryCurrencyCode: String("CAD"),
			BusinessType:        String("business type"),
			BusinessSubtype:     String("business sub-type"),
			OrganizationalType:  String("organizational type"),
			Address1:            String("Address 1"),
			Address2:            String("Address 2"),
			City:                String("City"),
			Province: &Province{
				Name: String("Ontario"),
				Slug: String("ontario"),
			},
			Country: &Country{
				Name:         String("Canada"),
				CountryCode:  String("CA"),
				CurrencyCode: String("CAD"),
				URL:          String("url"),
			},
			PostalCode:          String("A1B 2C3"),
			PhoneNumber:         String("416-555-5555"),
			MobilePhoneNumber:   String("416-555-5556"),
			TollFreePhoneNumber: String("416-555-5557"),
			FaxNumber:           String("416-555-5558"),
			Website:             String("https://example.com"),
			IsPersonalBusiness:  Bool(false),
			DateCreated:         &datetime,
			DateModified:        &datetime,
		}
		checkMarshalJSON(b, expectedBusinessJSON)
	})

	Convey("LIST all owned businesses", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "GET")
			fmt.Fprint(w, expectedBusinessesJSON)
		})

		businesses, _, err := client.Businesses.List()
		b := []Business{*expectedBusinessStruct}
		So(err, ShouldBeNil)
		So(businesses, ShouldResemble, b)
	})

	Convey("GET a specific business", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses/1", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "GET")
			fmt.Fprint(w, expectedBusinessJSON)
		})

		business, _, err := client.Businesses.Get("1")
		So(err, ShouldBeNil)
		So(business, ShouldResemble, expectedBusinessStruct)
	})

	Convey("GET a specific Business with an invalid ID", t, func() {
		business, resp, err := client.Businesses.Get("%")
		checkInvalidURLError(business, resp, err)
	})

	Convey("CREATE a Business", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "POST")
			fmt.Fprint(w, expectedBusinessJSON)
		})

		b := Business{}
		business, _, err := client.Businesses.Create(b)
		So(err, ShouldBeNil)
		So(business, ShouldResemble, expectedBusinessStruct)
	})

	Convey("REPLACE a Business", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses/1", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "PUT")
			fmt.Fprint(w, expectedBusinessJSON)
		})

		b := Business{}
		business, _, err := client.Businesses.Replace("1", b)
		So(err, ShouldBeNil)
		So(business, ShouldResemble, expectedBusinessStruct)
	})

	Convey("REPLACE a Business with an invalid ID", t, func() {
		business, resp, err := client.Businesses.Replace("%", Business{})
		checkInvalidURLError(business, resp, err)
	})

	Convey("UPDATE a Business", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses/1", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "PATCH")
			fmt.Fprint(w, expectedBusinessJSON)
		})

		b := Business{}
		business, _, err := client.Businesses.Update("1", b)
		So(err, ShouldBeNil)
		So(business, ShouldResemble, expectedBusinessStruct)
	})

	Convey("UPDATE a Business with an invalid ID", t, func() {
		business, resp, err := client.Businesses.Update("%", Business{})
		checkInvalidURLError(business, resp, err)
	})

	Convey("String method on Business", t, func() {
		b := new(Business)
		b.CompanyName = String("Company Test")
		b.IsPersonalBusiness = Bool(true)
		b.ID = String("id")
		So(b.String(), ShouldEqual, "Company Test (id=id, personal=true)")
	})
}
