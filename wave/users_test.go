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
	expectedUserJSON = `{
		"id":"9c9956e5-e995-4779-97ac-3d220b56152e",
		"url":"https://api.example.com/user/",
		"first_name":"Jane",
		"last_name":"Smith",
		"emails":[
			{
				"email":"jane@example.com",
				"is_verified":false,
				"is_default":true
			}
		],
		"profile":{
			"date_of_birth":"1981-02-17"
		},
		"businesses":[
			{
				"id":"c2731e5c-5001-4fe3-87ae-477f9a48dfcc",
				"url":"https://api.example.com/businesses/c2731e5c-5001-4fe3-87ae-477f9a48dfcc/"
			},
			{
				"id":"f99948d3-a16e-4082-8f3c-824f8cba6377",
				"url":"https://api.example.com/businesses/f99948d3-a16e-4082-8f3c-824f8cba6377/"
			}
		],
		"date_created":"2013-11-28T15:58:02+00:00",
		"date_modified":"2013-11-28T15:58:02+00:00",
		"last_login":"2013-11-28T15:58:03+00:00"
	}`
)

func TestUsersService(t *testing.T) {
	expectedUserStruct := new(User)
	json.Unmarshal([]byte(expectedUserJSON), expectedUserStruct)

	Convey("GET a specific User", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "GET")
			fmt.Fprint(w, expectedUserJSON)
		})

		user, _, err := client.Users.Get()
		So(err, ShouldBeNil)
		So(user, ShouldResemble, expectedUserStruct)
	})

	Convey("REPLACE a User", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "PUT")
			fmt.Fprint(w, expectedUserJSON)
		})

		u := User{}
		user, _, err := client.Users.Replace(u)
		So(err, ShouldBeNil)
		So(user, ShouldResemble, expectedUserStruct)
	})

	Convey("UPDATE a User", t, func() {
		setUp()
		defer tearDown()

		mux.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
			So(r.Method, ShouldEqual, "PATCH")
			fmt.Fprint(w, expectedUserJSON)
		})

		u := User{}
		user, _, err := client.Users.Update(u)
		So(err, ShouldBeNil)
		So(user, ShouldResemble, expectedUserStruct)
	})

	Convey("FullName method on User", t, func() {
		Convey("No FirstName and LastName should return ''", func() {
			u := new(User)
			So(u.FullName(), ShouldBeBlank)
		})
		Convey("No FirstName should return 'LastName'", func() {
			u := new(User)
			u.LastName = String("Bar")
			So(u.FullName(), ShouldEqual, "Bar")
		})
		Convey("No LastName should return 'FirstName'", func() {
			u := new(User)
			u.FirstName = String("Foo")
			So(u.FullName(), ShouldEqual, "Foo")
		})
		Convey("FirstName and LastName should return 'FirstName LastName'", func() {
			u := new(User)
			u.FirstName = String("Foo")
			u.LastName = String("Bar")
			So(u.FullName(), ShouldEqual, "Foo Bar")
		})
	})

	Convey("String method on User", t, func() {
		Convey("FirstName, LastName should return 'FirstName LastName'", func() {
			u := new(User)
			u.FirstName = String("Foo")
			u.LastName = String("Bar")
			So(u.String(), ShouldEqual, "Foo Bar")
		})

		Convey("No FirstName, LastName should return ''", func() {
			u := new(User)
			So(u.String(), ShouldBeBlank)
		})
	})
}
