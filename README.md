# Gowave

[![Build Status](https://travis-ci.org/NickPresta/gowave.png?branch=master)](https://travis-ci.org/NickPresta/gowave)
[![Coverage Status](http://i.imgur.com/pK6knhY.gif)](https://drone.io/github.com/NickPresta/gowave/files/coverage.html)
[![GoDoc](https://godoc.org/github.com/NickPresta/gowave/wave?status.png)](https://godoc.org/github.com/NickPresta/gowave/wave)

gowave is a Go client library for accessing the [Wave API](https://developer.waveapps.com).

gowave requires Go version 1.2 or greater.

**The wave package is in an ALPHA state.** There is no guarantee of interface stability until the Wave API is "final".

## Installation

To download, build and install the wave package, run:

	go get github.com/NickPresta/gowave/wave


## Getting Started

Import the wave package into your source:

```go
import (
  "github.com/NickPresta/gowave/wave"
)
```

You will need to create a new Wave client, which will allow you to make requests
on behalf of a user. The wave package does not handle authentication and
requires you to provide an http.Client that can handle appropriate
authentication. The easiest and recommended way to do this is using the [goauth2](https://code.google.com/p/goauth2/) package.

```go
t := &oauth.Transport{
  Token: &oauth.Token{AccessToken: "... your access token ..."},
}

client := wave.NewClient(t.Client())
```

## Creating and Updating Resources

All structs in this library use pointer values so that there is differentiation
between an unset value and a zero value. This also allows the same structs to be
encoded and decoded without having to learn two different data types.  Helper
methods are provided to create pointer values for string, int, float64, and
bool:

```go
product := &wave.Product{
	Name: wave.String("Widgets"),
	Price: wave.Float64(42.34),
	IsSold: wave.Bool(true),
}
client.Products.Create(bID, product)
```

## Optional Parameters

Some endpoints take optional parameters -- usually LIST and GET methods. For
example, with Products, you can choose to embed the accounts directly in the
Product resource. Optional parameters are passed as a pointer to a struct. If
you do not wish to pass any options and want to take the API defaults, set the
options struct to nil in the function argument. If you do not wish to pass a
specific option, you may omit it entirely from the struct. For example:

```go
client.Products.List(bID, &wave.ProductListOptions{EmbedAccounts: true})
```

### Pagination

Pagination options are passed in the optional parameters:

```go
options = &wave.ProductListOptions{PageOptions: wave.PageOptions{Page: 5, PageSize: 10}}
client.Products.List(bID, options)
```

Again, omitting the PageOptions struct will not send any pagination parameters.

## Examples

### Fetch all Accounts for a given Business

```go
accounts, resp, err := client.Accounts.List(businessID)
if err != nil {
	panic(err)
}
// Do something with accounts
for account := range accounts {
	fmt.Println(*account.Name)
}
// The resp.Body still exists to be read, if you wanted to do further
// processing or something
io.Copy(os.Stdout, resp.Response.Body)
```

### Create a Business

```go
b := &Business{
	CompanyName:         "My New Business",
	PrimaryCurrencyCode: "CAD",
	BusinessTypeInfo: &BusinessTypeInfo{
		BusinessType:       String("consultants_professionals"),
		BusinessSubtype:    String("consultants_professionals__communications"),
		OrganizationType:   String("partnership"),
	},
	Address: &Address{
		Country: &Country{
			CountryCode: "CA",
		},
	},
}
business, _, err = client.Businesses.Create(b)
// Do something with business
```

### Delete a Customer

```go
resp, err := client.Customers.Delete(businessID, customerID)
if err == nil {
	// Customer deleted. Success!
}
```

## Thanks and Inspiration

This library is heavily inspired by [go-github](https://github.com/google/go-github), although there is no affiliation
or endorsement in any way by the go-github contributors or Google Inc itself.

## License

This client library is distributed under the BSD-style license found in the [LICENSE](./LICENSE).
