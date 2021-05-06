package repository

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const (
	API_URL = "https://subsplease.org/api"
)

type SubsPlease struct {
	c *http.Client
}

type episode struct {
	Time        string           `json:"time"`
	ReleaseDate string           `json:"release_date"`
	Show        string           `json:"show"`
	Episode     string           `json:"episode"`
	Downloads   []downloadOption `json:"downloads"`
	Xdcc        string           `json:"xdcc"`
	ImageUrl    string           `json:"image_url"`
	Page        string           `json:"page"`
}

type downloadOption struct {
	Res    string `json:"res"`
	Magnet string `json:"magnet"`
}

func NewSubsPlease() *SubsPlease {
	return &SubsPlease{
		c: http.DefaultClient,
	}
}

func (s *SubsPlease) Latests() ([]episode, error) {
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

	var res map[string]episode
	err = json.Unmarshal(b, &res)
	if err != nil {
		return nil, err
	}

	rv := make([]episode, 0, len(res))
	for _, v := range res {
		rv = append(rv, v)
	}
	return rv, nil
}
