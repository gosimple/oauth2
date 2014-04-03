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
	// Create new "Client ID" at https://code.google.com/apis/console/
	// under "API access" and provide clientId (-id), clientSecret (-secret)
	// and redirectURL (-redirect) as imput arguments.
	clientId     = flag.String("id", "", "Client ID")
	clientSecret = flag.String("secret", "", "Client Secret")
	redirectURL  = flag.String("redirect", "http://httpbin.org/get", "Redirect URL")
	scope        = flag.String("scope", "https://www.googleapis.com/auth/userinfo.profile", "Scope")

	authURL    = "https://accounts.google.com/o/oauth2/auth"
	tokenURL   = "https://accounts.google.com/o/oauth2/token"
	apiBaseURL = "https://www.googleapis.com/oauth2/v2/"
)

const startInfo = `
Create new "Client ID" at https://code.google.com/apis/console/
under "API access" and provide -id, -secret and -redirect as input arguments.
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

	service.Scope = *scope

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
	google := oauth2.Request(apiBaseURL, token.AccessToken)
	google.AccessTokenInHeader = true

	// Make the request.
	// Provide API end point (https://www.googleapis.com/discovery/v1/apis/oauth2/v2/rest)
	apiEndPoint := "userinfo"
	googleUserData, err := google.Get(apiEndPoint)
	if err != nil {
		log.Fatalf("Get: %v", err)
	}
	defer googleUserData.Body.Close()

	fmt.Println("User info response:")
	// Write the response to standard output.
	io.Copy(os.Stdout, googleUserData.Body)

	fmt.Println()
}
