package wave

import (
	"encoding/json"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"net/http"
	"testing"
)

func businessJSON() string {
	return `{"id":"","url":"","company_name":"","primary_currency_code":"","phone_number":null,"mobile_phone_number":null,"toll_free_phone_number":null,"fax_number":null,"website":null,"is_personal_business":false,"date_created":"0001-01-01T00:00:00","date_modified":"0001-01-01T00:00:00","business_type_info":{"business_type":null,"business_subtype":null,"organizational_type":null},"address_info":{"address1":null,"address2":null,"city":null,"postal_code":null,"province":{"name":null,"slug":null},"country":{"name":null,"country_code":null,"currency_code":null,"url":""}}}`
}

func businessStruct() *Business {
	var b Business
	json.Unmarshal([]byte(businessJSON()), &b)
	return &b
}

func TestBusinessesService(t *testing.T) {
	Convey("LISTing all owned businesses", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "GET")
			fmt.Fprint(w, "["+businessJSON()+"]")
		})

		businesses, _, err := client.Businesses.List()
		b := []Business{*businessStruct()}
		So(err, ShouldEqual, nil)
		So(businesses, ShouldResemble, b)
	})

	Convey("GET a specific business", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/businesses/1", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "GET")
			fmt.Fprint(w, businessJSON())
		})

		business, _, err := client.Businesses.Get("1")
		So(err, ShouldEqual, nil)
		So(business, ShouldResemble, businessStruct())
	})
}
