OAuth2
======

Package oauth2 provide support for OAuth 2.0 authentication and ability
to make authenticated HTTP requests.

[Documentation online](http://godoc.org/github.com/gosimple/oauth2)

## Status

Working for now but can break in future (like moving ```Get```, ```Post```
methods to ```Service``` to allow auto token refresh etc.).

## Example

Check ```examples``` folder for usages.

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
		log.Fatal("Get access token error: ", err)
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
		log.Fatal("Get: ", err)
	}
	defer githubUserData.Body.Close()

### Requests or bugs?

<https://github.com/gosimple/oauth2/issues>

## Installation

	go get github.com/gosimple/oauth2

## License

The source files are distributed under the
[Mozilla Public License, version 2.0](http://mozilla.org/MPL/2.0/),
unless otherwise noted.
Please read the [FAQ](http://www.mozilla.org/MPL/2.0/FAQ.html)
if you have further questions regarding the license.
