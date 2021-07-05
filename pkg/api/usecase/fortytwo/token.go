package fortytwo

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/url"
	"os"
)

type Token struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type TokenInfo struct {
	ResourceOwnerID  int      `json:"resource_owner_id"`
	Scopes           []string `json:"scopes"`
	ExpiresInSeconds int      `json:"expires_in_seconds"`
	CreatedAt        int      `json:"created_at"`
	Application      struct {
		UID string `json:"uid"`
	} `json:"application"`
}

func (s *Service) GetToken(code, state string) (*Token, error) {
	if err := s.state.Exist(state); err != nil {
		return nil, err
	}
	s.state.Delete(state)

	params := url.Values{
		"grant_type":    {"authorization_code"},
		"client_id":     {os.Getenv("AUTH_42_CLIENT_ID")},
		"client_secret": {os.Getenv("AUTH_42_SECRET")},
		"code":          {code},
		"redirect_uri":  {os.Getenv("AUTH_42_REDIRECT_URI")},
		"state":         {state},
	}

	resp, err := s.client.PostForm("https://api.intra.42.fr/oauth/token", params)
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

func (s *Service) GetTokenOwnerInfo(accessToken string) (*TokenInfo, error) {
	baseUrl, err := url.Parse("https://api.intra.42.fr/oauth/token/info")
	if err != nil {
		return nil, err
	}

	params := url.Values{}
	params.Add("access_token", accessToken)

	baseUrl.RawQuery = params.Encode()

	resp, err := s.client.Get(baseUrl.String())
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

	var tInfo TokenInfo

	err = json.Unmarshal(body, &tInfo)
	if err != nil {
		return nil, err
	}
	return &tInfo, nil
}
