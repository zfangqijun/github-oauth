package main

import (
	"fmt"

	"github.com/imroc/req/v3"
)

var clientId = "35267af0118570d03009"
var clientSecret = "34f8bda5539ca6ddd308655563e67a6729aac9ca"
var LoginOAuthAccessToken = "/login/oauth/access_token"

var client = req.C().SetBaseURL("https://github.com")

type OAuthResult struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func OAuthByCode(code string) (OAuthResult, error) {
	var result OAuthResult

	resp, err := client.R().
		SetHeader("Accept", "application/json").
		SetQueryParams(map[string]string{
			"client_id":     clientId,
			"client_secret": clientSecret,
			"code":          code,
		}).
		Post(LoginOAuthAccessToken)

	resp.Into(&result)

	fmt.Println(result.Scope)
	fmt.Println(result.AccessToken)
	fmt.Println(result.TokenType)

	return result, err
}
