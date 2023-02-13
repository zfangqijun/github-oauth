package main

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	"github.com/imroc/req/v3"
	"golang.org/x/oauth2"
)

type OAuthResult struct {
	AccessToken string `json:"access_token"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
	ErrorType   string `json:"error"`
}

func OAuthByCode(code string) ([]byte, error) {
	var clientId = "35267af0118570d03009"
	var clientSecret = "34f8bda5539ca6ddd308655563e67a6729aac9ca"
	var LoginOAuthAccessToken = "/login/oauth/access_token"

	resp, err := req.C().
		SetBaseURL("https://github.com").
		R().
		SetHeader("Accept", "application/json").
		SetQueryParams(map[string]string{
			"client_id":     clientId,
			"client_secret": clientSecret,
			"code":          code,
		}).
		Post(LoginOAuthAccessToken)

	result, err := resp.ToBytes()

	if !resp.IsSuccessState() {
		fmt.Println("bad response status:", resp.Status)
	}

	return result, err
}

func GetRepos(token string) ([]*github.Repository, error) {
	ctx := context.Background()

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)

	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	repos, resp, err := client.Repositories.List(ctx, "", nil)
	fmt.Printf("%s\n", resp.Request.Proto)

	if err != nil {
		return nil, err
	}

	return repos, err
}
