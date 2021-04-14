package torrent

import (
	"crypto/rand"
	"io"
	"os"
	"strconv"
	"strings"
)

var (
	DEFAULT_TRACKER_PORT uint16 = 6881
)

type TorrentFile struct {
	Announce     string
	AnnounceList []string
	InfoHash     [20]byte
	Name         string
	PieceHashes  [][20]byte
	PieceLength  int
	Length       int
}

func Parse(r io.Reader) (TorrentFile, error) {
	bto, err := Open(r)
	if err != nil {
		return TorrentFile{}, err
	}
	return bto.toTorrentFile()
}

func ParseFromFile(filepath string) (TorrentFile, error) {
	r, err := os.Open(filepath)
	if err != nil {
		return TorrentFile{}, err
	}
	return Parse(r)
}

func generatePeerID() ([20]byte, error) {
	peerID := [20]byte{}
	_, err := rand.Read(peerID[:])
	if err != nil {
		return [20]byte{}, err
	}
	return peerID, nil
}

func (t *TorrentFile) defaultTracker() (Tracker, error) {
	peerID, err := generatePeerID()
	if err != nil {
		return Tracker{}, err
	}
	return t.buildTracker(t.Announce, peerID)
}

func (t *TorrentFile) Trackers() ([]Tracker, error) {
	if len(t.AnnounceList) == 0 {
		tr, err := t.defaultTracker()
		if err != nil {
			return nil, err
		}
		return []Tracker{tr}, nil
	}

	output := make([]Tracker, 0, len(t.AnnounceList))
	for _, v := range t.AnnounceList {
		if !strings.HasPrefix(v, "http://") {
			continue
		}

		peerID, err := generatePeerID()
		if err != nil {
			return nil, err
		}

		tr, err := t.buildTracker(v, peerID)
		if err != nil {
			return nil, err
		}
		output = append(output, tr)
	}
	return output, nil
}

func (t *TorrentFile) buildTracker(announce string, peerID [20]byte) (Tracker, error) {
	config := TrackerConfig{
		Announce:   announce,
		Port:       DEFAULT_TRACKER_PORT,
		Uploaded:   "0",
		Downloaded: "0",
		Compact:    "1",
		Left:       strconv.Itoa(t.Length),
	}
	copy(config.PeerID[:], peerID[:])
	return NewTracker(config)
}
