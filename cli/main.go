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
	clientSecret = flag.String("secret", "", "Client Secret")
	scope        = flag.String("scope", "basic", "Scope")
	cachefile    = flag.String("cache", "cache.json", "Token cache file")
	port         = flag.Int("port", 9001, "Webserver port")
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
	}
	_, err := config.TokenCache.Token()
	if err != nil {
		getToken(t)
	}

	client := wave.NewClient(t.Client())
	businesses, _, err := client.Businesses.List()
	if err != nil {
		log.Fatalf("Couldn't fetch businesses: %+v", err)
	}

	b, err := json.MarshalIndent(businesses, "", "  ")
	if err != nil {
		log.Fatalf("Error decoding businesses: %+v\n", err)
	}
	os.Stdout.Write(b)
	fmt.Println()
}
