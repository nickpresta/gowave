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
	expectedCustomerJSON = `{
    "id": 1,
    "url": "https://api.waveapps.com/businesses/c2731e5c-5001-4fe3-87ae-477f9a48dfcc/customers/1/",
    "account_number": "12345",
    "customer_name": "Mike Ehrmantraut",
    "first_name": "Mike",
    "last_name": "Ehrmantraut",
    "email": "mike@hitmenforhire.com",
    "fax_number": "555-578-9112",
    "mobile_number": "555-553-1212",
    "phone_number": "555-578-9111",
    "toll_free_number": "1-800-555-5555",
    "website": "http://www.pi-by-mike.com/",
    "currency": {
        "url": "https://api.waveapps.com/currencies/USD/",
        "code": "USD",
        "symbol": "$",
        "name": "U.S. dollar"
    },
    "address": {
        "address1": "1100 Central Avenue Southeast",
        "address2": "",
        "city": "Albuquerque",
        "province": {
            "name": "New Mexico",
            "slug": "new-mexico"
        },
        "country": {
            "name": "United States",
            "country_code": "US",
            "currency_code": "USD",
            "url": "https://api.waveapps.com/countries/US/"
        },
        "postal_code": "87106"
    },
    "shipping_details": {
        "ship_to_contact": "Mike Ehrmantraut",
        "delivery_instructions": "Leave it in a black bag by the tree.",
        "phone_number": "555-553-1212",
        "address": {
            "address1": "1100 Central Avenue Southeast",
            "address2": "",
            "city": "Albuquerque",
            "province": {
                "name": "New Mexico",
                "slug": "new-mexico"
            },
            "country": {
                "name": "United States",
                "country_code": "US",
                "currency_code": "USD",
                "url": "https://api.waveapps.com/countries/US/"
            },
            "postal_code": "87106"
        }
    },
    "date_created": "2013-12-05T10:31:01+00:00",
    "date_modified": "2013-12-05T13:37:59+00:00"
}`
	expectedCustomersJSON = "[" + expectedCustomerJSON + "]"
)

func TestCustomersService(t *testing.T) {
	expectedCustomerStruct := new(Customer)
	json.Unmarshal([]byte(expectedCustomerJSON), expectedCustomerStruct)

	Convey("LIST all Customers for a business", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses/1/customers/", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "GET")
			fmt.Fprint(w, expectedCustomersJSON)
		})

		customers, _, err := client.Customers.List("1", nil)
		c := []Customer{*expectedCustomerStruct}
		So(err, ShouldEqual, nil)
		So(customers, ShouldResemble, c)
	})

	Convey("LIST all Customers for a business invalid ID", t, func() {
		_, resp, err := client.Customers.List("%", nil)
		checkInvalidURLError(nil, resp, err)
	})

	Convey("GET a specific Customer", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses/1/customers/1/", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "GET")
			fmt.Fprint(w, expectedCustomerJSON)
		})

		customer, _, err := client.Customers.Get("1", 1)
		So(err, ShouldBeNil)
		So(customer, ShouldResemble, expectedCustomerStruct)
	})

	Convey("GET a specific Customer with an invalid ID", t, func() {
		customer, resp, err := client.Customers.Get("%", 1)
		checkInvalidURLError(customer, resp, err)
	})

	Convey("CREATE a Customer", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses/1/customers/", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "POST")
			fmt.Fprint(w, expectedCustomerJSON)
		})

		c := &Customer{}
		customers, _, err := client.Customers.Create("1", c)
		So(err, ShouldBeNil)
		So(customers, ShouldResemble, expectedCustomerStruct)
	})

	Convey("CREATE a Customer with an invalid ID", t, func() {
		customer, resp, err := client.Customers.Create("%", &Customer{})
		checkInvalidURLError(customer, resp, err)
	})

	Convey("REPLACE a Customer", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses/1/customers/1/", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "PUT")
			fmt.Fprint(w, expectedCustomerJSON)
		})

		c := &Customer{}
		customer, _, err := client.Customers.Replace("1", 1, c)
		So(err, ShouldBeNil)
		So(customer, ShouldResemble, expectedCustomerStruct)
	})

	Convey("REPLACE a Customer with an invalid ID", t, func() {
		customer, resp, err := client.Customers.Replace("%", 1, &Customer{})
		checkInvalidURLError(customer, resp, err)
	})

	Convey("UPDATE a Customer", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses/1/customers/1/", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "PATCH")
			fmt.Fprint(w, expectedCustomerJSON)
		})

		c := &Customer{}
		customer, _, err := client.Customers.Update("1", 1, c)
		So(err, ShouldEqual, nil)
		So(customer, ShouldResemble, expectedCustomerStruct)
	})

	Convey("UPDATE a Customer with an invalid ID", t, func() {
		customer, resp, err := client.Customers.Update("%", 1, &Customer{})
		checkInvalidURLError(customer, resp, err)
	})

	Convey("DELETE a specific Customer", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses/1/customers/1/", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "DELETE")
			w.WriteHeader(http.StatusNoContent)
		})

		_, err := client.Customers.Delete("1", 1)
		So(err, ShouldEqual, nil)
	})

	Convey("DELETE a specific Customer with an invalid ID", t, func() {
		resp, err := client.Customers.Delete("%", 1)
		checkInvalidURLError(nil, resp, err)
	})

	Convey("FullName method on Customer", t, func() {
		Convey("No FirstName and LastName should return ''", func() {
			c := new(Customer)
			So(c.FullName(), ShouldBeBlank)
		})
		Convey("Only CustomerName should return 'CustomerName'", func() {
			c := new(Customer)
			c.CustomerName = String("Foo Bar")
			So(c.FullName(), ShouldEqual, "Foo Bar")
		})
		Convey("No FirstName should return 'LastName'", func() {
			c := new(Customer)
			c.LastName = String("Bar")
			So(c.FullName(), ShouldEqual, "Bar")
		})
		Convey("No LastName should return 'FirstName'", func() {
			c := new(Customer)
			c.FirstName = String("Foo")
			So(c.FullName(), ShouldEqual, "Foo")
		})
		Convey("FirstName and LastName should return 'FirstName LastName'", func() {
			c := new(Customer)
			c.FirstName = String("Foo")
			c.LastName = String("Bar")
			So(c.FullName(), ShouldEqual, "Foo Bar")
		})
	})

	Convey("String method on Customer", t, func() {
		Convey("FirstName, LastName, Email should return 'FirstName LastName (Email)'", func() {
			c := new(Customer)
			c.FirstName = String("Foo")
			c.LastName = String("Bar")
			c.Email = String("foo@example.com")
			So(c.String(), ShouldEqual, "Foo Bar (foo@example.com)")
		})

		Convey("No Email should return 'FirstName LastName'", func() {
			c := new(Customer)
			c.FirstName = String("Foo")
			c.LastName = String("Bar")
			So(c.String(), ShouldEqual, "Foo Bar")
		})

		Convey("No FirstName, LastName should return Email", func() {
			c := new(Customer)
			c.Email = String("foo@example.com")
			So(c.String(), ShouldEqual, "foo@example.com")
		})

		Convey("No FirstName, LastName, Email should return ''", func() {
			c := new(Customer)
			So(c.String(), ShouldBeBlank)
		})
	})
}
