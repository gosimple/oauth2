// Copyright 2013 Dobrosław Żybort
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package oauth2

import (
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

var _ = fmt.Printf

// ResRequest represents values needed to make authenticated HTTP requests.
type ResRequest struct {
	// Base URL for API
	apiBaseURL url.URL
	// Access token added to each resource request
	AccessToken string
	// Header allows you to add custom headers that'll be added to each
	// resource request
	Header http.Header

	// The OAuth 2.0 Authorization Framework: Bearer Token Usage
	// http://tools.ietf.org/html/rfc6750

	// Set AccessTokenInURL to true if destination service require
	// authorization in the HTTP request URI
	//		GET /resource?access_token=<YOUR_ACCESS_TOKEN> HTTP/1.1
	//		Host: server.example.com
	AccessTokenInURL bool
	// Authentication URI parameter, default: "access_token"
	AccessTokenInURLParam string

	// Set AccessTokenInHeader to true if destination service require
	// authorization in the "Authorization" request header
	//		GET /resource HTTP/1.1
	//		Host: server.example.com
	//		Authorization: Bearer <YOUR_ACCESS_TOKEN>
	AccessTokenInHeader bool
	// Authentication header scheme, default: "Bearer"
	AccessTokenInHeaderScheme string
}

// Request initializes basic values that can be used to make
// authenticated HTTP requests.
func Request(apiBaseURL, accessToken string) *ResRequest {
	req := new(ResRequest)
	req.Header = make(http.Header)
	req.AccessToken = accessToken

	req.setApiBaseURL(apiBaseURL)

	req.AccessTokenInURLParam = "access_token"
	req.AccessTokenInHeaderScheme = "Bearer"
	return req
}

// setApiBaseURL parse and update req.ApiBaseURL
func (req *ResRequest) setApiBaseURL(baseURL string) {
	apiBaseURL, err := url.Parse(strings.TrimRight(baseURL, "/"))
	if err != nil {
		panic("ApiBaseURL error: " + err.Error())
	}
	req.apiBaseURL = *apiBaseURL
}

// ApiBaseURL update req.ApiBaseURL
func (req *ResRequest) ApiBaseURL(baseURL string) {
	req.setApiBaseURL(baseURL)
}

// buildURL build full URL from req.apiBaseURL and endPoint
func (req *ResRequest) buildURL(endPoint string) string {
	endPoint = strings.TrimLeft(endPoint, "/")
	return req.apiBaseURL.String() + "/" + endPoint
}

// updateTokenInURL add access token to fullURL
// if req.AccessTokenInURL is set to true.
func (req *ResRequest) updateTokenInURL(fullURL string) (string, error) {
	if req.AccessTokenInURL {
		parsedURL, err := url.Parse(fullURL)
		if err != nil {
			return "", errors.New("Error updating token in URL")
		}
		params, _ := url.ParseQuery(parsedURL.RawQuery)
		params.Set(req.AccessTokenInURLParam, req.AccessToken)
		parsedURL.RawQuery = params.Encode()
		fullURL = parsedURL.String()
	}
	return fullURL, nil
}

// updateTokenInHeader add access token to request header
// if req.AccessTokenInHeader is set to true.
func (req *ResRequest) updateTokenInHeader(request *http.Request) (
	updatedRequest *http.Request) {
	if req.AccessTokenInHeader {
		authHeader := req.AccessTokenInHeaderScheme + " " + req.AccessToken
		request.Header.Set("Authorization", authHeader)
	}
	return request
}

// Do updates HTTP request with access token. Next sends an HTTP request
// and returns an HTTP response
//func (req *Request) Do(req *http.Request) (
//	resp *http.Response, err error) {
//	return service.Client().Do(req)
//}

// Delete issues a DELETE to the specified API endpoint.
func (req *ResRequest) Delete(endPoint string) (
	resp *http.Response, err error) {
	return req.sendRequest("DELETE", endPoint, nil)
}

// Get issues a GET to the specified API endpoint.
func (req *ResRequest) Get(endPoint string) (
	resp *http.Response, err error) {
	return req.sendRequest("GET", endPoint, nil)
}

// Head issues a HEAD to the specified API endpoint.
func (req *ResRequest) Head(endPoint string) (
	resp *http.Response, err error) {
	return req.sendRequest("HEAD", endPoint, nil)
}

// Options issues a OPTIONS to the specified API endpoint.
func (req *ResRequest) Options(endPoint string) (
	resp *http.Response, err error) {
	return req.sendRequest("OPTIONS", endPoint, nil)
}

// Patch issues a PATCH to the specified API endpoint, with data's keys
// and values urlencoded as the request body.
func (req *ResRequest) Patch(endPoint string, data url.Values) (
	resp *http.Response, err error) {
	return req.sendRequest("PATCH", endPoint, data)
}

// Post issues a POST to the specified API endpoint, with data's keys
// and values urlencoded as the request body.
func (req *ResRequest) Post(endPoint string, data url.Values) (
	resp *http.Response, err error) {
	return req.sendRequest("POST", endPoint, data)
}

// Put issues a PUT to the specified API endpoint, with data's keys and values
// urlencoded as the request body.
func (req *ResRequest) Put(endPoint string, data url.Values) (
	resp *http.Response, err error) {
	return req.sendRequest("PUT", endPoint, data)
}

// Trace issues a TRACE to the specified API endpoint.
func (req *ResRequest) Trace(endPoint string) (
	resp *http.Response, err error) {
	return req.sendRequest("TRACE", endPoint, nil)
}

// sendRequest issues OAuth-authenticated request method to the specified
// API endpoint, with data's keys and values URL-encoded as the request body.
// Caller should close resp.Body when done reading from it.
func (req *ResRequest) sendRequest(method, endPoint string, data url.Values) (
	resp *http.Response, err error) {
	fullURL := req.buildURL(endPoint)

	fullURL, err = req.updateTokenInURL(fullURL)
	if err != nil {
		return nil, err
	}

	var encData string
	var body io.ReadCloser
	if data != nil {
		encData = data.Encode()
		reader := strings.NewReader(encData)
		body = ioutil.NopCloser(reader)
	}

	//fmt.Println(fullURL)

	request, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, errors.New("Error building request")
	}

	request.Header = req.Header
	request = req.updateTokenInHeader(request)

	if data != nil {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		request.ContentLength = int64(len(encData))
	}

	//for key, val := range request.Header {
	//	fmt.Println(key, val)
	//}

	return http.DefaultClient.Do(request)
}
