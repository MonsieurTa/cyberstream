package tracker

import (
	"encoding/binary"
	"net"

	"github.com/MonsieurTa/hypertube/pkg/ft-torrent/peer"
)

type TrackerResponse struct {
	failureReason string
	interval      int
	trackerID     string
	peers         []byte
	complete      int
	incomplete    int
}

func NewTrackerResponse(m Decoder) TrackerResponse {
	return TrackerResponse{
		failureReason: m.GetString("failure reason"),
		interval:      m.GetInt("interval"),
		trackerID:     m.GetString("tracker id"),
		peers:         []byte(m.GetString("peers")),
		complete:      m.GetInt("complete"),
		incomplete:    m.GetInt("incomplete"),
	}
}

func (tr *TrackerResponse) Peers() ([]peer.Peer, error) {
	const peerSize = 6

	if len(tr.peers)%peerSize != 0 {
		return nil, err_malformed_response
	}

	nbPeers := len(tr.peers) / peerSize
	rv := make([]peer.Peer, nbPeers)
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

func (tr *TrackerResponse) Failed() bool {
	return tr.failureReason != ""
}
