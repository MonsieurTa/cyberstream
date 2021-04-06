package fortytwo

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/MonsieurTa/hypertube/config"
	"golang.org/x/oauth2/clientcredentials"
)

type Service struct {
	client *http.Client
}

func createAuthClient() *http.Client {
	ctx := context.Background()

	conf := clientcredentials.Config{
		ClientID:     config.PROVIDER_42_CLIENT_ID,
		ClientSecret: config.PROVIDER_42_SECRET,
		TokenURL:     "https://api.intra.42.fr/oauth/token",
	}
	token, err := conf.Token(ctx)
	if err != nil {
		log.Fatalf("could not get 42 provider oauth token")
	}

	tokenJSON, err := json.MarshalIndent(token, "", "\t")
	if err != nil {
		log.Fatalf(err.Error())
	}
	fmt.Printf("Using 42 oauth2 client credential:\n%s\n", string(tokenJSON))
	return conf.Client(ctx)
}

func NewService() *Service {
	return &Service{
		client: createAuthClient(),
	}
}
