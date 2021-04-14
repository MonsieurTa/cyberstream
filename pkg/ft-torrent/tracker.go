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
	PeerID     [20]byte
	Port       uint16
	Uploaded   string
	Downloaded string
	Compact    string
	Left       string
	baseURL    *url.URL
}

type TrackerResponse struct {
	Interval int
	Peers    string
}

var (
	err_malformed_response = errors.New("malformed tracker response")
)

func (t *Tracker) url(infoHash [20]byte) (string, error) {
	params := url.Values{
		"info_hash":  []string{string(infoHash[:])},
		"peer_id":    []string{string(t.PeerID[:])},
		"port":       []string{strconv.Itoa(int(t.Port))},
		"uploaded":   []string{t.Uploaded},
		"downloaded": []string{t.Downloaded},
		"compact":    []string{t.Compact},
		"left":       []string{t.Left},
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

func NewTrackerResponse(data map[string]interface{}) (TrackerResponse, error) {
	rv := TrackerResponse{}

	rawInterval, ok1 := data["interval"]
	interval, ok2 := rawInterval.(int64)
	if !ok1 || !ok2 {
		return TrackerResponse{}, err_malformed_response
	}
	rv.Interval = int(interval)

	rawPeers, ok1 := data["peers"]
	peers, ok2 := rawPeers.([]byte)
	if !ok1 || !ok2 {
		return TrackerResponse{}, err_malformed_response
	}
	rv.Peers = string(peers[:])
	return rv, nil
}
