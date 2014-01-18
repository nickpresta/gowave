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
	expectedAccountJSON = `{
		"id": 1,
		"url":"url",
		"name":"Account Name",
		"active":true,
		"account_class":"bank",
		"account_type":"asset",
		"standard_account_number":2,
		"account_template_id":42,
		"account_number":3,
		"is_payment":false,
		"can_delete":true,
		"currency":{
			"url":"url",
			"code":"CAD",
			"symbol":"$",
			"name":"Canadian dollar"
		},
		"is_currency_editable":true,
		"is_name_editable":true,
		"is_payment_editable":true,
		"date_created":"2009-11-10T23:00:00+00:00",
		"date_modified":"2009-11-10T23:00:00+00:00"
	}`
	expectedAccountsJSON = "[" + expectedAccountJSON + "]"
)

func TestAccountsService(t *testing.T) {
	expectedAccountStruct := new(Account)
	json.Unmarshal([]byte(expectedAccountJSON), expectedAccountStruct)

	Convey("Testing JSON marshalling of an Account", t, func() {
		datetime := DateTime(time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC))
		a := &Account{
			ID:                    Int(1),
			URL:                   String("url"),
			Name:                  String("Account Name"),
			Active:                Bool(true),
			AccountClass:          String("bank"),
			AccountType:           String("asset"),
			StandardAccountNumber: Int(2),
			AccountTemplateID:     Int(42),
			AccountNumber:         Int(3),
			IsPayment:             Bool(false),
			CanDelete:             Bool(true),
			Currency: &Currency{
				URL:    String("url"),
				Code:   String("CAD"),
				Symbol: String("$"),
				Name:   String("Canadian dollar"),
			},
			IsCurrencyEditable: Bool(true),
			IsNameEditable:     Bool(true),
			IsPaymentEditable:  Bool(true),
			DateCreated:        &datetime,
			DateModified:       &datetime,
		}
		checkMarshalJSON(a, expectedAccountJSON)
	})

	Convey("LIST all Accounts for a business", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses/1/accounts", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "GET")
			fmt.Fprint(w, expectedAccountsJSON)
		})

		accounts, _, err := client.Accounts.List("1")
		a := []Account{*expectedAccountStruct}
		So(err, ShouldEqual, nil)
		So(accounts, ShouldResemble, a)
	})

	Convey("LIST all Accounts for a business invalid ID", t, func() {
		_, resp, err := client.Accounts.List("%")
		checkInvalidURLError(nil, resp, err)
	})

	Convey("GET a specific Account", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses/1/accounts/1", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "GET")
			fmt.Fprint(w, expectedAccountJSON)
		})

		account, _, err := client.Accounts.Get("1", 1)
		So(err, ShouldBeNil)
		So(account, ShouldResemble, expectedAccountStruct)
	})

	Convey("GET a specific Account with an invalid ID", t, func() {
		account, resp, err := client.Accounts.Get("%", 1)
		checkInvalidURLError(account, resp, err)
	})

	Convey("CREATE an Account", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses/1/accounts", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "POST")
			fmt.Fprint(w, expectedAccountJSON)
		})

		a := &Account{}
		account, _, err := client.Accounts.Create("1", a)
		So(err, ShouldBeNil)
		So(account, ShouldResemble, expectedAccountStruct)
	})

	Convey("CREATE an Account with an invalid ID", t, func() {
		account, resp, err := client.Accounts.Create("%", &Account{})
		checkInvalidURLError(account, resp, err)
	})

	Convey("REPLACE an Account", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses/1/accounts/1", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "PUT")
			fmt.Fprint(w, expectedAccountJSON)
		})

		a := &Account{}
		account, _, err := client.Accounts.Replace("1", 1, a)
		So(err, ShouldBeNil)
		So(account, ShouldResemble, expectedAccountStruct)
	})

	Convey("REPLACE an Account with an invalid ID", t, func() {
		account, resp, err := client.Accounts.Replace("%", 1, &Account{})
		checkInvalidURLError(account, resp, err)
	})

	Convey("UPDATE an Account", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses/1/accounts/1", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "PATCH")
			fmt.Fprint(w, expectedAccountJSON)
		})

		a := &Account{}
		account, _, err := client.Accounts.Update("1", 1, a)
		So(err, ShouldEqual, nil)
		So(account, ShouldResemble, expectedAccountStruct)
	})

	Convey("UPDATE an Account with an invalid ID", t, func() {
		account, resp, err := client.Accounts.Update("%", 1, &Account{})
		checkInvalidURLError(account, resp, err)
	})

	Convey("DELETE a specific Account", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses/1/accounts/1", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "DELETE")
			w.WriteHeader(http.StatusNoContent)
		})

		_, err := client.Accounts.Delete("1", 1)
		So(err, ShouldEqual, nil)
	})

	Convey("DELETE a specific Account with an invalid ID", t, func() {
		resp, err := client.Accounts.Delete("%", 1)
		checkInvalidURLError(nil, resp, err)
	})

	Convey("String method on Account", t, func() {
		a := new(Account)
		a.Name = String("Account Test")
		a.AccountType = String("expense")
		a.IsPayment = Bool(true)
		So(a.String(), ShouldEqual, "Account Test (type=expense, payment=true)")
	})
}
