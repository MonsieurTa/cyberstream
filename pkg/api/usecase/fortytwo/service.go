package fortytwo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"golang.org/x/oauth2/clientcredentials"
)

type Service struct {
	client *http.Client
}

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

func NewService() (*Service, error) {
	authClient, err := createAuthClient()
	if err != nil {
		return nil, err
	}
	return &Service{
		client: authClient,
	}, nil
}

func (s *Service) GetAccessToken(code, state string) (*Token, error) {
	form := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {os.Getenv("AUTH_42_CLIENT_ID")},
		"client_secret": {os.Getenv("AUTH_42_SECRET")},
		"code":          {code},
		"redirect_uri":  {os.Getenv("AUTH_42_REDIRECT_URI")},
		"state":         {state},
	}

	resp, err := s.client.PostForm("https://api.intra.42.fr/oauth/token", form)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, errors.New(string(body))
	}

	var token Token
	err = json.Unmarshal(body, &token)
	if err != nil {
		return nil, err
	}
	return &token, nil
}

func (s *Service) GetAuthorizeURI(state string) (string, error) {
	params := "client_id=%s&redirect_uri=%s&state=%s&response_type=code"
	params = fmt.Sprintf(
		params,
		os.Getenv("AUTH_42_CLIENT_ID"),
		os.Getenv("AUTH_42_REDIRECT_URI"),
		state,
	)
	rv := os.Getenv("AUTH_42_AUTH_URI") + "?" + params
	return rv, nil
}

func createAuthClient() (*http.Client, error) {
	ctx := context.Background()

	conf := clientcredentials.Config{
		ClientID:     os.Getenv("AUTH_42_CLIENT_ID"),
		ClientSecret: os.Getenv("AUTH_42_SECRET"),
		TokenURL:     "https://api.intra.42.fr/oauth/token",
	}
	token, err := conf.Token(ctx)
	if err != nil {
		return nil, err
	}

	tokenJSON, err := json.MarshalIndent(token, "", "\t")
	if err != nil {
		return nil, err
	}
	fmt.Printf("Using 42 oauth2 client credentials:\n%s\n", string(tokenJSON))
	return conf.Client(ctx), nil
}