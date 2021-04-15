package torrent

import (
	"encoding/binary"
	"net"
)

type TrackerResponse struct {
	interval      int
	peers         []byte
	failureReason string
}

func NewTrackerResponse(m bencodeMap) TrackerResponse {
	return TrackerResponse{
		interval:      m.GetInt("interval"),
		peers:         []byte(m.GetString("peers")),
		failureReason: m.GetString("failure reason"),
	}
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
