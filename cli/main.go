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
	clientID     = flag.String("id", "", "Client ID")
	clientSecret = flag.String("secret", "", "Client secret")
	accessToken  = flag.String("access", "", "Access token")
	scope        = flag.String("scope", "basic", "Scope")
	cachefile    = flag.String("cache", "cache.json", "Token cache file")

	debug = flag.Bool("debug", false, "Sets debug mode")

	page     = flag.Int("page", 1, "Set which page you want to add to the query string")
	pageSize = flag.Int("pageSize", 1, "How many entires per page")

	businesses = flag.Bool("businesses", false, "LIST your businesses")
	business   = flag.String("business", "", "GET your business. Takes the ID")
	currencies = flag.Bool("currencies", false, "LIST the currencies")
	currency   = flag.String("currency", "", "GET a currency. Takes the code")
	countries  = flag.Bool("countries", false, "LIST the countries")
	country    = flag.String("country", "", "GET a country. Takes the code")
	customers  = flag.String("customers", "", "LIST the customers for a business. Takes the business ID")
	products   = flag.String("products", "", "LIST the products for a business. Takes the business ID")
	provinces  = flag.String("provinces", "", "LIST the provinces for a country. Takes the country code")
	user       = flag.Bool("user", false, "GET the currently authenticated user")
	accounts   = flag.String("accounts", "", "LIST the accounts for a business. Takes the ID")
)

var config *oauth.Config

func getToken(t *oauth.Transport) {
	var c string
	authURL := config.AuthCodeURL("state")
	log.Printf("Open in browser: %v\n", authURL)
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
		ClientId:     *clientID,
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

	if *customers != "" {
		customers, resp, err := client.Customers.List(*customers, &wave.CustomerListOptions{PageOptions: wave.PageOptions{Page: *page, PageSize: *pageSize}})
		p(customers, resp, err)
	}

	if *products != "" {
		products, resp, err := client.Products.List(*products, &wave.ProductListOptions{PageOptions: wave.PageOptions{Page: *page, PageSize: *pageSize}})
		p(products, resp, err)
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

func p(v interface{}, resp *wave.Response, err error) {
	if *debug {
		fmt.Printf("Request for \"%v\"\n", resp.Response.Request.URL)
	}
	fatal(resp, err)
	printPaginationInformation(resp)
	printResource(v)
}

func fatal(resp *wave.Response, err error) {
	if err != nil {
		log.Printf("Couldn't fetch: %+v", err)
		log.Printf("Raw response from \"%v?access_token=%v\":\n", resp.Response.Request.URL, *accessToken)
		io.Copy(os.Stderr, resp.Response.Body)
		log.Println()
		os.Exit(1)
	}
}

func printPaginationInformation(r *wave.Response) {
	fmt.Printf("Total Count: %v\nCurrent Page: %v\nNext Page: %v\nPrevious Page: %v\n", r.TotalCount, r.CurrentPage, r.NextPage, r.PreviousPage)
}

func printResource(r interface{}) {
	b, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		log.Fatalf("Error decoding resource: %+v\n", err)
	}
	os.Stdout.Write(b)
	fmt.Println()
}
