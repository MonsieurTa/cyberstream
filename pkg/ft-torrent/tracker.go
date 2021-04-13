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

type PeerResponse struct {
	Interval int
	Peers    string
}

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

func (t *Tracker) RequestPeers(infoHash [20]byte) (PeerResponse, error) {
	url, err := t.url(infoHash)

	if err != nil {
		return PeerResponse{}, err
	}

	client := http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return PeerResponse{}, err
	}

	defer resp.Body.Close()
	rawBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return PeerResponse{}, err
	}

	rawData, err := bencode.Unmarshal(rawBody)
	if err != nil {
		return PeerResponse{}, err
	}

	data, ok := rawData.(map[string]interface{})
	if !ok {
		return PeerResponse{}, errors.New("malformed tracker response")
	}
	return NewPeerResponse(data)
}

func NewPeerResponse(data map[string]interface{}) (PeerResponse, error) {
	rv := PeerResponse{}

	rawInterval, ok1 := data["interval"]
	interval, ok2 := rawInterval.(int64)
	if !ok1 || !ok2 {
		return PeerResponse{}, errors.New("malformed tracker response")
	}
	rv.Interval = int(interval)

	rawPeers, ok1 := data["peers"]
	peers, ok2 := rawPeers.([]byte)
	if !ok1 || !ok2 {
		return PeerResponse{}, errors.New("malformed tracker response")
	}
	rv.Peers = string(peers[:])
	return rv, nil
}
