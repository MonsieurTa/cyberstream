package fortytwo

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"os"

	"github.com/MonsieurTa/hypertube/pkg/api/internal/inmem"
	"golang.org/x/oauth2/clientcredentials"
)

type Service struct {
	client *http.Client
	state  inmem.StateInMem
}

func NewService() (*Service, error) {
	authClient, err := createAuthClient()
	if err != nil {
		return nil, err
	}
	return &Service{
		client: authClient,
		state:  inmem.StateInMem{},
	}, nil
}

func (s *Service) GetAuthorizeURI() (string, error) {
	state, err := inmem.GenerateState()
	if err != nil {
		return "", err
	}

	baseUrl, err := url.Parse(os.Getenv("AUTH_42_AUTH_URI"))
	if err != nil {
		return "", err
	}

	params := url.Values{}
	params.Add("client_id", os.Getenv("AUTH_42_CLIENT_ID"))
	params.Add("redirect_uri", os.Getenv("AUTH_42_REDIRECT_URI"))
	params.Add("state", state)
	params.Add("response_type", "code")

	baseUrl.RawQuery = params.Encode()

	s.state.Save(state)
	return baseUrl.String(), nil
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
