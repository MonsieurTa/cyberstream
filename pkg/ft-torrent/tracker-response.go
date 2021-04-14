package torrent

import (
	"encoding/binary"
	"fmt"
	"net"
)

type TrackerResponse struct {
	interval int
	peers    []byte
}

var (
	err_key = func(key string) error { return fmt.Errorf("key error: %v", key) }
)

func NewTrackerResponse(data map[string]interface{}) (TrackerResponse, error) {
	rv := TrackerResponse{}

	rawInterval, ok1 := data["interval"]
	interval, ok2 := rawInterval.(int64)
	if !ok1 || !ok2 {
		return TrackerResponse{}, err_key("interval")
	}
	rv.interval = int(interval)

	rawPeers, ok1 := data["peers"]
	peers, ok2 := rawPeers.([]byte)
	if !ok1 || !ok2 {
		return TrackerResponse{}, err_key("peers")
	}
	rv.peers = peers[:]
	return rv, nil
}

func (tr *TrackerResponse) Peers() ([]Peer, error) {
	const peerSize = 6

	if len(tr.peers)%peerSize != 0 {
		return nil, err_malformed_response
	}

	nbPeers := len(tr.peers) / peerSize
	rv := make([]Peer, nbPeers)
	for i := 0; i < nbPeers; i++ {
		offset := i * peerSize
		rv[i].SetIP(net.IP(tr.peers[offset : offset+4]))
		rv[i].SetPort(binary.BigEndian.Uint16(tr.peers[offset+4 : offset+6]))
	}
	return rv, nil
}

func (tr *TrackerResponse) Interval() int {
	return tr.interval
}
