package tracker

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"time"

	b "github.com/MonsieurTa/hypertube/pkg/ft-torrent/bencode"
	"github.com/MonsieurTa/hypertube/pkg/ft-torrent/common"
	"github.com/marksamman/bencode"
)

type Trackers []Tracker

type Tracker struct {
	announce   string
	peerID     [20]byte
	port       uint16
	uploaded   string
	downloaded string
	compact    string
	left       string
}

var (
	DEFAULT_TRACKER_PORT uint16 = 6881

	err_malformed_response = errors.New("malformed tracker response")
	err_invalid_url        = errors.New("invalid url")
)

func NewTracker(announce string, peerID [20]byte, left int) (Tracker, error) {
	rv := Tracker{
		announce:   announce,
		peerID:     peerID,
		port:       DEFAULT_TRACKER_PORT,
		uploaded:   "0",
		downloaded: "0",
		compact:    "1",
		left:       strconv.Itoa(left),
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
	data, err := bencode.Decode(resp.Body)
	if err != nil {
		return TrackerResponse{}, err
	}

	return NewTrackerResponse(b.Decoder(data)), nil
}

func (t *Tracker) Protocol() string {
	return t.announce[0:3]
}

func (trs Trackers) RequestPeers(infoHash [20]byte) []common.Peer {
	rv := make([]common.Peer, 0, len(trs)*50)
	for _, tracker := range trs {
		resp, err := tracker.RequestPeers(infoHash)

		if err != nil || resp.Failed() {
			continue
		}
		peers, err := resp.Peers()
		if err != nil {
			fmt.Println("invalid peers")
			continue
		}
		rv = append(rv, peers...)
	}
	return rv
}
