package main

import (
	"code.google.com/p/goauth2/oauth"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/nickpresta/gowave/wave"
	"log"
	"net/http"
	"os"
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
		businesses, _, err := client.Businesses.List()
		fatal(err)
		printResource(businesses)
	}

	if *business != "" {
		business, _, err := client.Businesses.Get(*business)
		fatal(err)
		printResource(business)
	}

	if *currencies {
		currencies, _, err := client.Currencies.List()
		fatal(err)
		printResource(currencies)
	}

	if *currency != "" {
		currency, _, err := client.Currencies.Get(*currency)
		fatal(err)
		printResource(currency)
	}

	if *countries {
		countries, _, err := client.Countries.List()
		fatal(err)
		printResource(countries)
	}

	if *country != "" {
		country, _, err := client.Countries.Get(*country)
		fatal(err)
		printResource(country)
	}

	if *provinces != "" {
		provinces, _, err := client.Countries.Provinces(*provinces)
		fatal(err)
		printResource(provinces)
	}
}

func fatal(err error) {
	if err != nil {
		log.Fatalf("Couldn't fetch: %+v", err)
	}
}

func printResource(r interface{}) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		log.Fatalf("Error decoding currencies: %+v\n", err)
	}
	os.Stdout.Write(b)
	fmt.Println()
}
