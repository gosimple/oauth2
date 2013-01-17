// Copyright 2013 Dobrosław Żybort
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package oauth2

import (
	"errors"
	//"fmt"
	"net/http"
	"net/url"
	"strings"
)

type ResRequest struct {
	apiBaseURL url.URL
	Header     http.Header
	token      string

	// http://tools.ietf.org/html/rfc6750
	accessTokenInURL        bool
	accessTokenInURLText    string
	accessTokenInHeader     bool
	accessTokenInHeaderText string
}

func Request(apiBaseURL, token string) *ResRequest {
	setApiBaseURL(apiBaseURL)

	req := new(ResRequest)
	req.Header = make(http.Header)
	req.token = token
	return req
}

func setApiBaseURL(baseURL string) {
	apiBaseURL, err := url.Parse(strings.TrimRight(baseURL, "/"))
	if err != nil {
		panic("ApiBaseURL error: " + err.Error())
	}
	req.apiBaseURL = *apiBaseURL
}

func (req *ResRequest) SetApiBaseURL(baseURL string) {
	setApiBaseURL(baseURL)
}

func (req *ResRequest) AccessTokenInURL(val bool) {
	req.accessTokenInURL = val
	req.accessTokenInURLText = "access_token"
}

func (req *ResRequest) SetAccessTokenInURLParam(param string) {
	req.accessTokenInURLText = param
}

func (req *ResRequest) AccessTokenInHeader(val bool) {
	req.accessTokenInHeader = val
	req.accessTokenInHeaderText = "Bearer"
}

func (req *ResRequest) SetAccessTokenInHeaderString(text string) {
	req.accessTokenInHeaderText = text
}

//func (req *Request) Do(req *http.Request) (
//	resp *http.Response, err error) {
//	return service.Client().Do(req)
//}

func (req *ResRequest) Get(endPoint string) (resp *http.Response, err error) {
	endPoint = strings.TrimLeft(endPoint, "/")
	fullURL := req.apiBaseURL.String() + "/" + endPoint

	if req.accessTokenInURL {
		parsedURL, err := url.Parse(fullURL)
		if err != nil {
			return nil, errors.New("Error building GET request")
		}
		params, _ := url.ParseQuery(parsedURL.RawQuery)
		params.Set(req.accessTokenInURLText, req.token)
		parsedURL.RawQuery = params.Encode()
		fullURL = parsedURL.String()
	}

	//fmt.Println(fullURL)

	request, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, errors.New("Error building GET request")
	}
	request.Header = req.Header

	if req.accessTokenInHeader {
		authHeader := req.accessTokenInHeaderText + " " + req.token
		request.Header.Set("Authorization", authHeader)
	}

	//for key, val := range request.Header {
	//	fmt.Println(key, val)
	//}

	return http.DefaultClient.Do(request)
}

//func (req *Request) Head(url string) (
//	resp *http.Response, err error) {
//	return service.Client().Head(url)
//}
//func (req *Request) Post(url string, bodyType string, body io.Reader) (
//	resp *http.Response, err error) {
//	return service.Client().Post(url, bodyType, body)
//}

//func (req *Request) PostForm(url string, data url.Values) (
//	resp *http.Response, err error) {
//	return service.Client().PostForm(url, data)
//}
