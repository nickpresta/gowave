package wave

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	expectedBusinessJSON = `{
		"id":"",
		"url":"",
		"company_name":"",
		"primary_currency_code":"",
		"phone_number":null,
		"mobile_phone_number":null,
		"toll_free_phone_number":null,
		"fax_number":null,
		"website":null,
		"is_personal_business":false,
		"date_created":"0001-01-01T00:00:00+00:00",
		"date_modified":"0001-01-01T00:00:00+00:00",
		"business_type_info":{
			"business_type":null,
			"business_subtype":null,
			"organizational_type":null
		},
		"address_info":{
			"address1":null,
			"address2":null,
			"city":null,
			"postal_code":null,
			"province":{
				"name":null,"slug":null
			},
			"country":{
				"name":null,
				"country_code":null,
				"currency_code":null,
				"url":""
			}
		}
	}`
	expectedBusinessesJSON = "[" + expectedBusinessJSON + "]"
)

func TestBusinessesService(t *testing.T) {
	expectedBusinessStruct := new(Business)
	json.Unmarshal([]byte(expectedBusinessJSON), &expectedBusinessStruct)

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
		b.CompanyName = "Company Test"
		b.IsPersonalBusiness = true
		b.Id = "id"
		So(b.String(), ShouldEqual, "Company Test (id=id, personal=true)")
	})
}
