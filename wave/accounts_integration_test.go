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
		bId := businesses[0].Id

		accounts, _, err := integrationClient.Accounts.List(bId)
		So(err, ShouldBeNil)
		So(accounts, ShouldNotBeNil)

		aId := accounts[0].Id
		account, _, err := integrationClient.Accounts.Get(bId, aId)
		So(err, ShouldBeNil)
		So(account, ShouldNotBeNil)

		a := Account{
			Name:                  String("Checking CREATE TEST Account"),
			IsPayment:             Bool(true),
			Currency:              &Currency{Code: "JPY"},
			StandardAccountNumber: Int(1006),
		}
		account, _, err = integrationClient.Accounts.Create(bId, a)
		So(err, ShouldBeNil)
		So(account, ShouldNotBeNil)
		So(*account.Name, ShouldEqual, *a.Name)
		So(*account.IsPayment, ShouldEqual, *a.IsPayment)
		So(account.Currency.Code, ShouldEqual, a.Currency.Code)
		So(*account.StandardAccountNumber, ShouldEqual, *a.StandardAccountNumber)

		resp, err := integrationClient.Accounts.Delete(bId, account.Id)
		So(err, ShouldBeNil)
		So(resp.StatusCode, ShouldEqual, http.StatusNoContent)
	})
}
