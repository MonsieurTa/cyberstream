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
	announce   string
	peerID     [20]byte
	port       uint16
	uploaded   string
	downloaded string
	compact    string
	left       string
}

type TrackerConfig struct {
	Announce   string
	PeerID     [20]byte
	Port       uint16
	Uploaded   string
	Downloaded string
	Compact    string
	Left       string
}

var (
	err_malformed_response = errors.New("malformed tracker response")
	err_invalid_url        = errors.New("invalid url")
)

func NewTracker(config TrackerConfig) (Tracker, error) {
	rv := Tracker{
		announce:   config.Announce,
		peerID:     config.PeerID,
		port:       config.Port,
		uploaded:   config.Uploaded,
		downloaded: config.Downloaded,
		compact:    config.Compact,
		left:       config.Left,
	}
	return rv, nil
}

func (t *Tracker) url(infoHash [20]byte) (string, error) {
	baseURL, err := url.Parse(t.announce)
	if err != nil {
		return "", err
	}

	params := url.Values{
		"info_hash":  []string{string(infoHash[:])},
		"peer_id":    []string{string(t.peerID[:])},
		"port":       []string{strconv.Itoa(int(t.port))},
		"uploaded":   []string{t.uploaded},
		"downloaded": []string{t.downloaded},
		"compact":    []string{t.compact},
		"left":       []string{t.left},
	}
	baseURL.RawQuery = params.Encode()
	return baseURL.String(), nil
}

func (t *Tracker) RequestPeers(infoHash [20]byte) (TrackerResponse, error) {
	url, err := t.url(infoHash)

	if err != nil {
		return TrackerResponse{}, err_invalid_url
	}

	client := http.Client{Timeout: 5 * time.Second}
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
