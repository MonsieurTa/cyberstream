package subsplease

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	API_URL = "https://subsplease.org/api"
)

type Episode struct {
	Time        string           `json:"time"`
	ReleaseDate string           `json:"release_date"`
	Show        string           `json:"show"`
	Episode     string           `json:"episode"`
	Downloads   []DownloadOption `json:"downloads"`
	Xdcc        string           `json:"xdcc"`
	ImageUrl    string           `json:"image_url"`
	Page        string           `json:"page"`
}

type DownloadOption struct {
	Res    string `json:"res"`
	Magnet string `json:"magnet"`
}

type subsPlease struct {
	c *http.Client
}

func NewSubsPlease() SubsPlease {
	return &subsPlease{
		c: http.DefaultClient,
	}
}

func (s *subsPlease) Latests() ([]Episode, error) {
	url := API_URL + `?f=latest&tz=Europe/Paris`
	resp, err := s.c.Get(url)
	if err != nil {
		return nil, err
	}

	b, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var res map[string]Episode
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}

	rv := make([]Episode, 0, len(res))
	for _, v := range res {
		rv = append(rv, v)
	}
	return rv, nil
}
