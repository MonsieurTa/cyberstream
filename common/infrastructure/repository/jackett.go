package repository

import (
	"encoding/xml"
	"io"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type Indexers struct {
	XMLName xml.Name  `xml:"indexers"`
	Indexer []Indexer `xml:"indexer"`
}

type Indexer struct {
	ID           string       `xml:"id,attr"`
	Configured   bool         `xml:"configured,attr"`
	Title        string       `xml:"title"`
	Description  string       `xml:"description"`
	Link         string       `xml:"link"`
	Language     string       `xml:"language"`
	Capabilities Capabilities `xml:"caps"`
}

type Capabilities struct {
	Server struct {
		Title string `xml:"title,attr"`
	} `xml:"server"`
	Searching struct {
		Search      Search `xml:"search"`
		TVSearch    Search `xml:"tv-search"`
		MovieSearch Search `xml:"movie-search"`
		MusicSearch Search `xml:"music-search"`
		AudioSearch Search `xml:"audio-search"`
		BookSearch  Search `xml:"book-search"`
	} `xml:"searching"`
	Categories struct {
		Category []Category `xml:"category"`
	} `xml:"categories"`
}

type Search struct {
	Available       string `xml:"available,attr"`
	SupportedParams string `xml:"supportedParams,attr"`
}

type Category struct {
	ID     string  `xml:"id,attr"`
	Name   string  `xml:"name,attr"`
	Subcat *Subcat `xml:"subcat,omitempty"`
}

type Subcat struct {
	ID   string `xml:"id,attr"`
	Name string `xml:"name,attr"`
}

type Jackett interface {
	Indexers(configured bool) (*Indexers, error)
}

type jackett struct {
	apikey string
	host   string
	path   string
}

func NewJackett() Jackett {
	apikey := os.Getenv("JACKETT_API_KEY")
	host := os.Getenv("JACKETT_HOST") + ":" + os.Getenv("JACKETT_PORT")
	path := "/api/" + os.Getenv("JACKETT_API_VERSION")
	return &jackett{apikey, host, path}
}

func makeBaseURL(scheme, host, path string) url.URL {
	rv := url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   path + "/indexers/all/results/torznab/api",
	}
	return rv
}

func (j *jackett) Indexers(configured bool) (*Indexers, error) {
	u := makeBaseURL("http", j.host, j.path)

	params := url.Values{}
	params.Set("t", "indexers")
	params.Set("configured", strconv.FormatBool(configured))
	params.Set("apikey", j.apikey)
	u.RawQuery = params.Encode()

	url := u.String()
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	b, err := io.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	var indexers Indexers
	err = xml.Unmarshal(b, &indexers)
	if err != nil {
		return nil, err
	}
	return &indexers, nil
}
