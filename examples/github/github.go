// Copyright 2013 Dobrosław Żybort
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"

	"bitbucket.org/gosimple/oauth2"
	//"github.com/toqueteos/webbrowser"
)

var (
	// Register new app at https://github.com/settings/applications and provide
	// clientId (-id), clientSecret (-secret) and redirectURL (-redirect)
	// as imput arguments.
	clientId     = flag.String("id", "", "Client ID")
	clientSecret = flag.String("secret", "", "Client Secret")
	redirectURL  = flag.String("redirect", "http://httpbin.org/get", "Redirect URL")

	authURL    = "https://github.com/login/oauth/authorize"
	tokenURL   = "https://github.com/login/oauth/access_token"
	apiBaseURL = "https://api.github.com/"
)

const startInfo = `
Register new app at https://github.com/settings/applications and provide
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
	url := service.GetAuthorizeURL("")
	fmt.Println(url)

	// Open authorization url in default system browser.
	// webbrowser.Open(url)

	fmt.Printf("\nVisit URL and provide code: ")
	code := ""
	// Read access code from cmd.
	fmt.Scanf("%s", &code)
	// Get access token.
	token, _ := service.GetAccessToken(code)
	fmt.Println()

	// Prepare resource request.
	github := oauth2.Request(apiBaseURL, token.AccessToken)
	github.AccessTokenInHeader = true
	github.AccessTokenInHeaderScheme = "token"
	//github.AccessTokenInURL = true

	// Make the request.
	// Provide API end point (http://developer.github.com/v3/users/#get-the-authenticated-user)
	apiEndPoint := "user"
	githubUserData, err := github.Get(apiEndPoint)
	if err != nil {
		log.Fatal("Get:", err)
	}
	defer githubUserData.Body.Close()

	fmt.Println("User info:")
	// Write the response to standard output.
	io.Copy(os.Stdout, githubUserData.Body)

	fmt.Println()
}
