package torrent

import (
	"crypto/rand"
	"fmt"
	"io"
	"strings"

	b "github.com/MonsieurTa/hypertube/pkg/ft-torrent/bencode"
	"github.com/MonsieurTa/hypertube/pkg/ft-torrent/tracker"
	"github.com/marksamman/bencode"
)

type Torrent struct {
	Announce     string
	AnnounceList []string
	InfoHash     [20]byte
	Name         string
	PieceHashes  [][20]byte
	PieceLength  int
	Length       int
}

type bencodeTorrent struct {
	announce     string
	announceList []string
	creationDate int
	comment      string
	createdBy    string
	encoding     string
	info         bencodeInfo
}

type bencodeInfo struct {
	name      string
	filesInfo []bencodeFileInfo
}

type bencodeFileInfo struct {
	pieces      string
	pieceLength int
	length      int
	path        string
}

var (
	err_malformed_piece = func(len int) error { return fmt.Errorf("malformed pieces: len = %d", len) }
)

func Read(r io.Reader) (*bencodeTorrent, error) {
	m, err := bencode.Decode(r)
	if err != nil {
		return nil, err
	}
	return unserialize(b.Decoder(m))
}

func unserialize(m Decoder) (*bencodeTorrent, error) {
	bto, err := newBencodeTorrent(m)
	if err != nil {
		return nil, err
	}

	info := m.GetDict("info")
	bto.info = newBencodeInfo(b.Decoder(info))
	return bto, nil
}

func newBencodeTorrent(meta Decoder) (*bencodeTorrent, error) {
	list := meta.GetList("announce-list")

	announceList := make([]string, len(list))
	for i, v := range list {
		announceList[i] = v.([]interface{})[0].(string)
	}

	return &bencodeTorrent{
		announce:     meta.GetString("announce"),
		announceList: announceList,
		creationDate: meta.GetInt("creation date"),
		comment:      meta.GetString("comment"),
		createdBy:    meta.GetString("created by"),
		encoding:     meta.GetString("encoding"),
	}, nil
}

func newBencodeInfo(d Decoder) bencodeInfo {
	name := d.GetString("name")
	files := d.GetList("files")
	if files != nil {
		fileList := make([]bencodeFileInfo, len(files))
		for i, f := range files {
			df := b.Decoder(f.(map[string]interface{}))

			fileList[i].pieces = df.GetString("pieces")
			fileList[i].pieceLength = df.GetInt("piece length")
			fileList[i].length = df.GetInt("length")
			fileList[i].path = df.GetString("path")
		}
		return bencodeInfo{name, fileList}
	}

	return bencodeInfo{
		name,
		[]bencodeFileInfo{{
			d.GetString("pieces"),
			d.GetInt("piece length"),
			d.GetInt("length"),
			"",
		}},
	}
}

func (b *bencodeInfo) hash() [20]byte {
	return hash(b)
}

func (b *bencodeFileInfo) splitPieceHashes() ([][20]byte, error) {
	hashLen := 20
	buf := []byte(b.pieces)
	if len(buf)%hashLen != 0 {
		return nil, err_malformed_piece(len(buf))
	}

	nbPieces := len(buf) / hashLen
	output := make([][20]byte, nbPieces)
	for i := 0; i < nbPieces; i++ {
		copy(output[i][:], buf[i*hashLen:(i+1)*hashLen])
	}
	return output, nil
}

func (b *bencodeTorrent) toTorrent() ([]Torrent, error) {
	rv := make([]Torrent, 0, len(b.info.filesInfo))
	infoHash := b.info.hash()
	for _, v := range b.info.filesInfo {
		pieceHashes, err := v.splitPieceHashes()
		if err != nil {
			continue
		}

		t := Torrent{
			Announce:     b.announce,
			AnnounceList: b.announceList[:],
			InfoHash:     infoHash,
			Name:         b.info.name + v.path,
			PieceHashes:  pieceHashes,
			PieceLength:  v.pieceLength,
			Length:       v.length,
		}
		rv = append(rv, t)
	}
	return rv, nil
}

func (t *Torrent) Trackers() (tracker.Trackers, error) {
	if len(t.AnnounceList) == 0 {
		tr, err := t.defaultTracker()
		if err != nil {
			return nil, err
		}
		return tracker.Trackers{tr}, nil
	}

	output := make([]tracker.Tracker, 0, len(t.AnnounceList))
	for _, v := range t.AnnounceList {
		// TODO: wss, udp
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

func (t *Torrent) buildTracker(announce string, peerID [20]byte) (tracker.Tracker, error) {
	return tracker.NewTracker(announce, peerID, t.Length)
}

func (t *Torrent) defaultTracker() (tracker.Tracker, error) {
	peerID, err := generatePeerID()
	if err != nil {
		return tracker.Tracker{}, err
	}
	return t.buildTracker(t.Announce, peerID)
}

func generatePeerID() ([20]byte, error) {
	peerID := [20]byte{}
	_, err := rand.Read(peerID[:])
	if err != nil {
		return [20]byte{}, err
	}
	return peerID, nil
}
