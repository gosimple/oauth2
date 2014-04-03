// Copyright 2013 Dobrosław Żybort
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"github.com/gosimple/oauth2"

	"flag"
	"fmt"
	"io"
	"log"
	"os"

	//"github.com/toqueteos/webbrowser"
)

var (
	// Register new app at https://bitly.com/a/oauth_apps and provide
	// clientId (-id), clientSecret (-secret) and redirectURL (-redirect)
	// as imput arguments.
	clientId     = flag.String("id", "", "Client ID")
	clientSecret = flag.String("secret", "", "Client Secret")
	redirectURL  = flag.String("redirect", "http://httpbin.org/get", "Redirect URL")

	authURL    = "https://bitly.com/oauth/authorize"
	tokenURL   = "https://api-ssl.bitly.com/oauth/access_token"
	apiBaseURL = "https://api-ssl.bitly.com/v3/"
)

const startInfo = `
Register new app at https://bitly.com/a/oauth_apps and provide
-id, -secret and -redirect as input arguments.
`

func main() {
	flag.Parse()

	if *clientId == "" || *clientSecret == "" {
		fmt.Println(startInfo)
		flag.Usage()
		os.Exit(2)
	}

	// Initialize service.
	service := oauth2.Service(
		*clientId, *clientSecret, authURL, tokenURL)
	service.RedirectURL = *redirectURL

	// Get authorization url.
	aUrl := service.GetAuthorizeURL("")
	fmt.Println()
	fmt.Printf("%v", aUrl)
	fmt.Println()

	// Open authorization url in default system browser.
	//webbrowser.Open(aUrl)

	fmt.Printf("\nVisit URL and provide code: ")
	code := ""
	// Read access code from cmd.
	fmt.Scanf("%s", &code)
	// Get access token.
	token, err := service.GetAccessToken(code)
	if err != nil {
		log.Fatalf("Get access token error: %v", err)
	}
	fmt.Println()
	fmt.Println("Token expiration time:", token.ExpirationTime)
	fmt.Println("Is token expired?:", token.Expired())

	// Prepare resource request.
	bitly := oauth2.Request(apiBaseURL, token.AccessToken)
	bitly.AccessTokenInURL = true

	// Make the request.
	// Provide API end point (http://dev.bitly.com/user_info.html#v3_user_info)
	apiEndPoint := "user/info"
	bitlyUserInfo, err := bitly.Get(apiEndPoint)
	if err != nil {
		log.Fatalf("Get: %v", err)
	}
	defer bitlyUserInfo.Body.Close()

	fmt.Println("User info response:")
	// Write the response to standard output.
	io.Copy(os.Stdout, bitlyUserInfo.Body)

	fmt.Println()
}
