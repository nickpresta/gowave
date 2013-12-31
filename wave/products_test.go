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
	expectedProductJSON = `{
		"id": 3,
		"url": "url",
		"name": "Product",
		"price": 13.37,
		"description": "Product Description",
		"is_sold": false,
		"is_bought": true,
		"income_account": null,
		"expense_account": {
			"id": 27,
			"url": "url",
			"name": "Expense Account",
			"active": true,
			"standard_account_number": 6004,
			"account_class": "expense",
			"account_type": "expense",
			"account_number": 6004,
			"is_payment": false,
			"can_delete": true,
			"currency": {
				"url": "url",
				"code": "USD",
				"symbol": "$",
				"name": "U.S. dollar"
			},
			"is_currency_editable": false,
			"is_name_editable": false,
			"is_payment_editable": false,
			"date_created":"2009-11-10T23:00:00+00:00",
			"date_modified":"2009-11-10T23:00:00+00:00"
		},
		"date_created":"2009-11-10T23:00:00+00:00",
		"date_modified":"2009-11-10T23:00:00+00:00"
	}`
	expectedProductsJSON = `{
"next": null,
"previous": null,
"results": [` + expectedProductJSON + "]}"
)

func TestProductsService(t *testing.T) {
	expectedProductStruct := new(Product)
	json.Unmarshal([]byte(expectedProductJSON), expectedProductStruct)

	Convey("LIST all Products for a business", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses/1/products", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "GET")
			fmt.Fprint(w, expectedProductsJSON)
		})

		products, _, err := client.Products.List("1", nil)
		c := []Product{*expectedProductStruct}
		So(err, ShouldEqual, nil)
		So(products, ShouldResemble, c)
	})

	Convey("LIST all Products for a business with an embed arugment", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses/1/products", func(w http.ResponseWriter, r *http.Request) {
			So(r.URL.RawQuery, ShouldEqual, "embed_accounts=true")
			So(r.Method, ShouldEqual, "GET")
			fmt.Fprint(w, expectedProductsJSON)
		})

		products, _, err := client.Products.List("1", &ProductListOptions{EmbedAccounts: true})
		c := []Product{*expectedProductStruct}
		So(err, ShouldEqual, nil)
		So(products, ShouldResemble, c)
	})

	Convey("LIST all Products for a business invalid ID", t, func() {
		_, resp, err := client.Products.List("%", nil)
		checkInvalidURLError(nil, resp, err)
	})

	Convey("GET a specific Product", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses/1/products/1", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "GET")
			fmt.Fprint(w, expectedProductJSON)
		})

		product, _, err := client.Products.Get("1", 1, nil)
		So(err, ShouldBeNil)
		So(product, ShouldResemble, expectedProductStruct)
	})

	Convey("GET a specific Product with an embed argument", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses/1/products/1", func(w http.ResponseWriter, r *http.Request) {
			So(r.URL.RawQuery, ShouldEqual, "embed_accounts=true")
			So(r.Method, ShouldEqual, "GET")
			fmt.Fprint(w, expectedProductJSON)
		})

		product, _, err := client.Products.Get("1", 1, &ProductGetOptions{EmbedAccounts: true})
		So(err, ShouldBeNil)
		So(product, ShouldResemble, expectedProductStruct)
	})

	Convey("GET a specific Product with an invalid ID", t, func() {
		product, resp, err := client.Products.Get("%", 1, nil)
		checkInvalidURLError(product, resp, err)
	})

	Convey("CREATE a Product", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses/1/products", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "POST")
			fmt.Fprint(w, expectedProductJSON)
		})

		p := &Product{}
		products, _, err := client.Products.Create("1", p)
		So(err, ShouldBeNil)
		So(products, ShouldResemble, expectedProductStruct)
	})

	Convey("CREATE a Product with an invalid ID", t, func() {
		product, resp, err := client.Products.Create("%", &Product{})
		checkInvalidURLError(product, resp, err)
	})

	Convey("REPLACE a Product", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses/1/products/1", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "PUT")
			fmt.Fprint(w, expectedProductJSON)
		})

		p := &Product{}
		product, _, err := client.Products.Replace("1", 1, p)
		So(err, ShouldBeNil)
		So(product, ShouldResemble, expectedProductStruct)
	})

	Convey("REPLACE a Product with an invalid ID", t, func() {
		product, resp, err := client.Products.Replace("%", 1, &Product{})
		checkInvalidURLError(product, resp, err)
	})

	Convey("UPDATE a Product", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses/1/products/1", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "PATCH")
			fmt.Fprint(w, expectedProductJSON)
		})

		p := &Product{}
		product, _, err := client.Products.Update("1", 1, p)
		So(err, ShouldEqual, nil)
		So(product, ShouldResemble, expectedProductStruct)
	})

	Convey("UPDATE a Product with an invalid ID", t, func() {
		product, resp, err := client.Products.Update("%", 1, &Product{})
		checkInvalidURLError(product, resp, err)
	})

	Convey("DELETE a specific Product", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses/1/products/1", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "DELETE")
			w.WriteHeader(http.StatusNoContent)
		})

		_, err := client.Products.Delete("1", 1)
		So(err, ShouldEqual, nil)
	})

	Convey("DELETE a specific Product with an invalid ID", t, func() {
		resp, err := client.Products.Delete("%", 1)
		checkInvalidURLError(nil, resp, err)
	})

	Convey("String method on Product", t, func() {
		Convey("Name should return 'Name'", func() {
			p := new(Product)
			p.Name = String("Foo")
			So(p.String(), ShouldEqual, "Foo")
		})
	})
}
