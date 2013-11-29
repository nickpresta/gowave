# Gowave

[![Build Status](https://drone.io/github.com/NickPresta/gowave/status.png)](https://drone.io/github.com/NickPresta/gowave/latest)
[coverage](https://drone.io/github.com/NickPresta/gowave/files/coverage.html)
[documentation]()

gowave is a Go client library for accessing the [Wave API](https://developer.waveapps.com).

gowave requires Go version 1.1 or greater.

## Usage

```go
import "github.com/nickpresta/gowave/wave"
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

I have included a copy of other licenses from code used or other projects included in the [NOTICE](./NOTICE) file.
