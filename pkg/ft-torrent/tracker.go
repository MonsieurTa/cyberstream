package torrent

import (
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/IncSW/go-bencode"
)

type Tracker struct {
	peerID     [20]byte
	port       uint16
	uploaded   string
	downloaded string
	compact    string
	left       string
	baseURL    *url.URL
}

var (
	err_malformed_response = errors.New("malformed tracker response")
)

type TrackerConfig struct {
	Announce   string
	PeerID     [20]byte
	Port       uint16
	Uploaded   string
	Downloaded string
	Compact    string
	Left       string
}

func NewTracker(config TrackerConfig) (Tracker, error) {
	baseURL, err := url.Parse(config.Announce)
	if err != nil {
		return Tracker{}, err
	}
	rv := Tracker{
		port:       config.Port,
		uploaded:   config.Uploaded,
		downloaded: config.Downloaded,
		compact:    config.Compact,
		left:       config.Left,
		baseURL:    baseURL,
	}
	copy(rv.peerID[:], config.PeerID[:])
	return rv, nil
}

func (t *Tracker) url(infoHash [20]byte) (string, error) {
	params := url.Values{
		"info_hash":  []string{string(infoHash[:])},
		"peer_id":    []string{string(t.peerID[:])},
		"port":       []string{strconv.Itoa(int(t.port))},
		"uploaded":   []string{t.uploaded},
		"downloaded": []string{t.downloaded},
		"compact":    []string{t.compact},
		"left":       []string{t.left},
	}
	t.baseURL.RawQuery = params.Encode()
	return t.baseURL.String(), nil
}

func (t *Tracker) RequestPeers(infoHash [20]byte) (TrackerResponse, error) {
	url, err := t.url(infoHash)

	if err != nil {
		return TrackerResponse{}, err
	}

	client := http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return TrackerResponse{}, err
	}

	defer resp.Body.Close()
	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return TrackerResponse{}, err
	}

	rawData, err := bencode.Unmarshal(rawBody)
	if err != nil {
		return TrackerResponse{}, err
	}

	data, ok := rawData.(map[string]interface{})
	if !ok {
		return TrackerResponse{}, err_malformed_response
	}
	return NewTrackerResponse(data)
}
