package torrent

import (
	"fmt"
	"io"
	"strings"

	b "github.com/MonsieurTa/hypertube/pkg/ft-torrent/bencode"
	"github.com/MonsieurTa/hypertube/pkg/ft-torrent/common"
	"github.com/MonsieurTa/hypertube/pkg/ft-torrent/download"
	"github.com/MonsieurTa/hypertube/pkg/ft-torrent/tracker"
	"github.com/marksamman/bencode"
)

type Torrent struct {
	announce     string
	announceList []string
	infoHash     [20]byte
	name         string
	pieceHashes  [][20]byte
	pieceLength  int
	length       int
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
			announce:     b.announce,
			announceList: b.announceList[:],
			infoHash:     infoHash,
			name:         b.info.name + v.path,
			pieceHashes:  pieceHashes,
			pieceLength:  v.pieceLength,
			length:       v.length,
		}
		rv = append(rv, t)
	}
	return rv, nil
}

func (t *Torrent) Trackers(peerID [20]byte) (tracker.Trackers, error) {
	if len(t.announceList) == 0 {
		tr, err := t.defaultTracker(peerID)
		if err != nil {
			return nil, err
		}
		return tracker.Trackers{tr}, nil
	}

	output := make([]tracker.Tracker, 0, len(t.announceList))
	for _, v := range t.announceList {
		// TODO: wss, udp
		if !strings.HasPrefix(v, "http://") {
			continue
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
	return tracker.NewTracker(announce, peerID, t.length)
}

func (t *Torrent) defaultTracker(peerID [20]byte) (tracker.Tracker, error) {
	return t.buildTracker(t.announce, peerID)
}

func (t *Torrent) Download() error {
	d, err := download.NewDownloader(t)
	if err != nil {
		return err
	}
	return d.Download()
}

func (t *Torrent) InfoHash() [20]byte {
	return t.infoHash
}

func (t *Torrent) PieceBounds(index int) (int, int) {
	start := index * t.pieceLength
	end := start + t.pieceLength
	if end > t.length {
		end = t.length
	}
	return start, end
}

func (t *Torrent) PieceSize(index int) int {
	start, end := t.PieceBounds(index)
	return end - start
}

func (t *Torrent) Peers(peerID [20]byte) ([]common.Peer, error) {
	trackers, err := t.Trackers(peerID)
	if err != nil {
		return nil, err
	}
	return trackers.RequestPeers(t.infoHash), nil
}

// Number of pieces
func (t *Torrent) Pieces() int {
	return len(t.pieceHashes)
}

func (t *Torrent) PieceHashes() [][20]byte {
	return t.pieceHashes
}

func (t *Torrent) Length() int {
	return t.length
}
