package torrent

import (
	"crypto/rand"
	"io"
	"os"
	"strconv"
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

func (t *TorrentFile) Trackers() ([]Tracker, error) {
	output := make([]Tracker, len(t.AnnounceList))

	for i, v := range t.AnnounceList {
		peerID := [20]byte{}
		_, err := rand.Read(peerID[:])
		if err != nil {
			return nil, err
		}

		output[i], err = t.buildTracker(v, peerID)
		if err != nil {
			return nil, err
		}
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
