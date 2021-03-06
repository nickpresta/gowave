// Copyright (c) 2013, Nick Presta
// All rights reserved.
//
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// +build integration

package wave

import (
	"log"
	"testing"
	. "github.com/smartystreets/goconvey/convey"
)

func TestBusinessessServiceIntegration(t *testing.T) {
	setUpIntegrations()

	// Skip until the API drops trailing slashes
	SkipConvey("Full end-to-end Business integration", t, func() {
		businesses, _, err := integrationClient.Businesses.List(nil)
		So(err, ShouldBeNil)
		bID := businesses[0].ID

		for i := 0; i < len(businesses); i++ {
			log.Printf("Business: %+v\n", businesses[i])
		}

		business, _, err := integrationClient.Businesses.Get(bID)
		So(err, ShouldBeNil)
		So(business, ShouldNotBeNil)

		b := &Business{
			CompanyName:         "CREATE TEST Business",
			PrimaryCurrencyCode: "CAD",
			BusinessTypeInfo: &BusinessTypeInfo{
				BusinessType:     String("consultants_professionals"),
				BusinessSubtype:  String("consultants_professionals__communications"),
				OrganizationType: String("partnership"),
			},
			Address: &Address{
				Country: &Country{
					CountryCode: "CA",
				},
			},
		}
		business, _, err = integrationClient.Businesses.Create(b)
		So(err, ShouldBeNil)
		So(business, ShouldNotBeNil)
		So(business.CompanyName, ShouldEqual, b.CompanyName)
		So(business.PrimaryCurrencyCode, ShouldEqual, b.PrimaryCurrencyCode)
	})
}
