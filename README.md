# Gowave

[![Build Status](https://travis-ci.org/NickPresta/gowave.png?branch=master)](https://travis-ci.org/NickPresta/gowave)
[![Coverage Status](https://coveralls.io/repos/NickPresta/gowave/badge.png?branch=master)](https://coveralls.io/r/NickPresta/gowave?branch=master)
[![GoDoc](https://godoc.org/github.com/NickPresta/gowave/wave?status.png)](https://godoc.org/github.com/NickPresta/gowave/wave)

gowave is a Go client library for accessing the [Wave API](https://developer.waveapps.com).

gowave requires Go version 1.2 or greater.

**This client library is in an ALPHA state.** There is no guarantee of interface stability until the Wave API is "final".

## Usage

```go
import "github.com/NickPresta/gowave/wave"
```

## Authentication

The gowave library does not directly handle authentication and relies on you to provide an `http.Client` that can handle authentication for you.
The easiest and recommended way to do this is using the [goauth2](https://code.google.com/p/goauth2/) library.

```go
t := &oauth.Transport{
  Token: &oauth.Token{AccessToken: "... your access token ..."},
}

client := wave.NewClient(t.Client())

// list all businesses for the authenticated user
businesses, response, err := client.Businesses.List()
```

## Thanks and Inspiration

This library is heavily inspired by [go-github](https://github.com/google/go-github), although there is no affiliation
or endorsement in any way by the go-github contributors or Google Inc itself.

## License

This client library is distributed under the BSD-style license found in the [LICENSE](./LICENSE).
