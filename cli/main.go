package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"code.google.com/p/goauth2/oauth"
	"github.com/nickpresta/gowave/wave"
)

var (
	clientId     = flag.String("id", "", "Client ID")
	clientSecret = flag.String("secret", "", "Client secret")
	accessToken  = flag.String("access", "", "Access token")
	scope        = flag.String("scope", "basic", "Scope")
	cachefile    = flag.String("cache", "cache.json", "Token cache file")

	businesses = flag.Bool("businesses", false, "LIST your businesses")
	business   = flag.String("business", "", "GET your business. Takes the ID")
	currencies = flag.Bool("currencies", false, "LIST the currencies")
	currency   = flag.String("currency", "", "GET a currency. Takes the code")
	countries  = flag.Bool("countries", false, "LIST the countries")
	country    = flag.String("country", "", "GET a country. Takes the code")
	provinces  = flag.String("provinces", "", "LIST the provinces for a country. Takes the country code")
	user       = flag.Bool("user", true, "GET the currently authenticated user")
	accounts   = flag.String("accounts", "", "GET the accounts for a business. Takes the ID")
)

var config *oauth.Config

func getToken(t *oauth.Transport) {
	var c string
	authUrl := config.AuthCodeURL("state")
	log.Printf("Open in browser: %v\n", authUrl)
	log.Printf("Enter verification code: ")
	fmt.Scanln(&c)
	_, err := t.Exchange(c)
	if err != nil {
		log.Fatalf("An error occurred exchanging the code: %v\n", err)
	}
}

func main() {
	flag.Parse()

	config = &oauth.Config{
		ClientId:     *clientId,
		ClientSecret: *clientSecret,
		Scope:        *scope,
		AuthURL:      "https://api.waveapps.com/oauth2/authorize/",
		TokenURL:     "https://api.waveapps.com/oauth2/token/",
		TokenCache:   oauth.CacheFile(*cachefile),
		RedirectURL:  "https://wave-portal.ngrok.com/oauth2",
	}

	t := &oauth.Transport{
		Config:    config,
		Transport: http.DefaultTransport,
		Token:     &oauth.Token{AccessToken: *accessToken},
	}

	if *accessToken == "" {
		_, err := config.TokenCache.Token()
		if err != nil {
			getToken(t)
		}
	}

	client := wave.NewClient(t.Client())

	if *businesses {
		businesses, resp, err := client.Businesses.List()
		fatal(resp, err)
		printResource(businesses)
	}

	if *business != "" {
		business, resp, err := client.Businesses.Get(*business)
		fatal(resp, err)
		printResource(business)
	}

	if *currencies {
		currencies, resp, err := client.Currencies.List()
		fatal(resp, err)
		printResource(currencies)
	}

	if *currency != "" {
		currency, resp, err := client.Currencies.Get(*currency)
		fatal(resp, err)
		printResource(currency)
	}

	if *countries {
		countries, resp, err := client.Countries.List()
		fatal(resp, err)
		printResource(countries)
	}

	if *country != "" {
		country, resp, err := client.Countries.Get(*country)
		fatal(resp, err)
		printResource(country)
	}

	if *provinces != "" {
		provinces, resp, err := client.Countries.Provinces(*provinces)
		fatal(resp, err)
		printResource(provinces)
	}

	if *user {
		user, resp, err := client.Users.Get()
		fatal(resp, err)
		printResource(user)
	}

	if *accounts != "" {
		accounts, resp, err := client.Accounts.List(*accounts)
		fatal(resp, err)
		printResource(accounts)
	}
}

func fatal(resp *http.Response, err error) {
	if err != nil {
		log.Printf("Couldn't fetch: %+v", err)
		log.Printf("Raw response from %v?access_token=%v :\n", resp.Request.URL, *accessToken)
		io.Copy(os.Stderr, resp.Body)
		log.Println()
		os.Exit(1)
	}
}

func printResource(r interface{}) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		log.Fatalf("Error decoding resource: %+v\n", err)
	}
	os.Stdout.Write(b)
	fmt.Println()
}
