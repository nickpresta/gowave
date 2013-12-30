// Copyright (c) 2013, Nick Presta
// All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build integration

package wave

import (
	"net/http"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAccountsServiceIntegration(t *testing.T) {
	setUpIntegrations()

	// Skip until the API drops trailing slashes
	SkipConvey("Full end-to-end Account integration", t, func() {
		businesses, _, err := integrationClient.Businesses.List()
		So(err, ShouldBeNil)
		bID := businesses[0].ID

		accounts, _, err := integrationClient.Accounts.List(bID)
		So(err, ShouldBeNil)
		So(accounts, ShouldNotBeNil)

		aID := accounts[0].ID
		account, _, err := integrationClient.Accounts.Get(bId, aID)
		So(err, ShouldBeNil)
		So(account, ShouldNotBeNil)

		a := Account{
			Name:                  String("Checking CREATE TEST Account"),
			IsPayment:             Bool(true),
			Currency:              &Currency{Code: "JPY"},
			StandardAccountNumber: Int(1006),
		}
		account, _, err = integrationClient.Accounts.Create(bID, a)
		So(err, ShouldBeNil)
		So(account, ShouldNotBeNil)
		So(*account.Name, ShouldEqual, *a.Name)
		So(*account.IsPayment, ShouldEqual, *a.IsPayment)
		So(account.Currency.Code, ShouldEqual, a.Currency.Code)
		So(*account.StandardAccountNumber, ShouldEqual, *a.StandardAccountNumber)

		resp, err := integrationClient.Accounts.Delete(bID, account.ID)
		So(err, ShouldBeNil)
		So(resp.StatusCode, ShouldEqual, http.StatusNoContent)
	})
}
