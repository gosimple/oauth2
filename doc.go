// Copyright 2013 Dobrosław Żybort
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

/*
Package oauth2 provide support for OAuth 2.0 authentication and ability
to make authenticated HTTP requests.

Check /examples/ folder for usages.

	// Initialize service.
	service := oauth2.Service(
		YOUR_CLIENT_ID,
		YOUR_CLIENT_SECRET,
		"https://github.com/login/oauth/authorize",
		"https://github.com/login/oauth/access_token"
	)

	// Set custom redirect
	service.RedirectURL = "http://you.example.org/handler"

	// Get authorization url.
	authUrl := service.GetAuthorizeURL("")

	// Send user to authUrl and get code
	code := "..."

	// Get access token.
	token, err := service.GetAccessToken(code)
	if err != nil {
		boo
	}

	// Prepare resource request.
	apiBaseURL = "https://api.github.com/"
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

Requests or bugs?
https://bitbucket.org/matrixik/go-oauth2/issues
*/
package oauth2
