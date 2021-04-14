package torrent

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/IncSW/go-bencode"
)

type bencodeTorrent struct {
	Announce     string
	AnnounceList []string
	Info         bencodeInfo
}

type bencodeInfo struct {
	Pieces      string
	PieceLength int
	Length      int
	Name        string
}

var (
	err_invalid_info          = errors.New("invalid info")
	err_invalid_pieces        = errors.New("invalid pieces")
	err_invalid_piece_length  = errors.New("invalid piece length")
	err_invalid_length        = errors.New("invalid length")
	err_invalid_name          = errors.New("invalid name")
	err_invalid_announce      = errors.New("invalid announce")
	err_invalid_announce_list = errors.New("invalid announce-list")

	err_invalid_data     = errors.New("invalid data")
	err_missing_announce = errors.New("missing announce")
	err_missing_info     = errors.New("missing info")

	err_malformed_piece = func(len int) error { return fmt.Errorf("malformed pieces: len = %d", len) }
)

func Open(r io.Reader) (bencodeTorrent, error) {
	buf, err := ioutil.ReadAll(r)
	if err != nil {
		return bencodeTorrent{}, err
	}

	data, err := bencode.Unmarshal(buf)
	if err != nil {
		return bencodeTorrent{}, err
	}
	return unserialize(data)
}

func (b *bencodeInfo) hash() ([20]byte, error) {
	dict := map[string]interface{}{
		"pieces":       []byte(b.Pieces),
		"piece length": b.PieceLength,
		"length":       b.Length,
		"name":         []byte(b.Name),
	}

	buf, err := bencode.Marshal(dict)
	if err != nil {
		return [20]byte{}, err
	}
	h := sha1.Sum(buf)
	return h, nil
}

func (b *bencodeInfo) splitPieceHashes() ([][20]byte, error) {
	hashLen := 20
	buf := []byte(b.Pieces)
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

func (b *bencodeTorrent) toTorrentFile() (TorrentFile, error) {
	pieceHashes, err := b.Info.splitPieceHashes()
	if err != nil {
		return TorrentFile{}, err
	}
	infoHash, err := b.Info.hash()

	if err != nil {
		return TorrentFile{}, err
	}
	return TorrentFile{
		Announce:     b.Announce,
		AnnounceList: b.AnnounceList,
		InfoHash:     infoHash,
		Name:         b.Info.Name,
		PieceHashes:  pieceHashes,
		PieceLength:  b.Info.PieceLength,
		Length:       b.Info.Length,
	}, nil
}

func (bto *bencodeTorrent) fill(rawAnnounce interface{}, rawAnnounceList interface{}) error {
	announce, ok := rawAnnounce.([]uint8)
	if !ok {
		return err_invalid_announce
	}
	bto.Announce = string(announce)

	announceList, ok := rawAnnounceList.([]interface{})
	if ok {
		bto.AnnounceList = make([]string, len(announceList))
		for i, v := range announceList {
			s, ok := v.([]interface{})
			if !ok {
				return err_invalid_announce_list
			}
			bto.AnnounceList[i] = string(s[0].([]uint8))
		}
	}
	return nil
}

func (i *bencodeInfo) fill(rawInfo interface{}) error {
	info, ok := rawInfo.(map[string]interface{})
	if !ok {
		return err_invalid_info
	}

	rawPieces, ok1 := info["pieces"]
	pieces, ok2 := rawPieces.([]uint8)
	if !ok1 || !ok2 {
		return err_invalid_pieces
	}
	i.Pieces = string(pieces)

	rawPieceLength, ok1 := info["piece length"]
	pieceLength, ok2 := rawPieceLength.(int64)
	if !ok1 || !ok2 {
		return err_invalid_piece_length
	}
	i.PieceLength = int(pieceLength)

	rawLength, ok1 := info["length"]
	length, ok2 := rawLength.(int64)
	if !ok1 || !ok2 {
		return err_invalid_length
	}
	i.Length = int(length)

	rawName, ok1 := info["name"]
	name, ok2 := rawName.([]uint8)
	if !ok1 || !ok2 {
		return err_invalid_name
	}
	i.Name = string(name)
	return nil
}

func unserialize(data interface{}) (bencodeTorrent, error) {
	m, ok := data.(map[string]interface{})
	if !ok {
		return bencodeTorrent{}, err_invalid_data
	}

	bto := bencodeTorrent{}
	announce, ok := m["announce"]
	if !ok {
		return bencodeTorrent{}, err_missing_announce
	}

	announceList, ok := m["announce-list"]
	if ok {
		err := bto.fill(announce, announceList)
		if err != nil {
			return bencodeTorrent{}, err
		}
	}

	info, ok := m["info"]
	if !ok {
		return bencodeTorrent{}, err_missing_info
	}

	err := bto.Info.fill(info)
	if err != nil {
		return bencodeTorrent{}, err
	}
	return bto, nil
}
