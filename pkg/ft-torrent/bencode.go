package torrent

import (
	"crypto/sha1"
	"fmt"
	"io"

	"github.com/marksamman/bencode"
)

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

type bencodeMap map[string]interface{}

var (
	err_malformed_piece = func(len int) error { return fmt.Errorf("malformed pieces: len = %d", len) }
)

func ReadTorrentFile(r io.Reader) (*bencodeTorrent, error) {
	m, err := bencode.Decode(r)
	if err != nil {
		return nil, err
	}
	return unserialize(bencodeMap(m))
}

func unserialize(m bencodeMap) (*bencodeTorrent, error) {
	bto, err := newBencodeTorrent(m)
	if err != nil {
		return nil, err
	}

	info := m.GetDict("info")
	bto.info = newBencodeInfo(info)
	return bto, nil
}

func newBencodeTorrent(meta bencodeMap) (*bencodeTorrent, error) {
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

func newBencodeInfo(info bencodeMap) bencodeInfo {
	name := info.GetString("name")
	files := info.GetList("files")
	if files != nil {
		fileList := make([]bencodeFileInfo, len(files))
		for i, f := range files {
			df := bencodeMap(f.(map[string]interface{}))

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
			info.GetString("pieces"),
			info.GetInt("piece length"),
			info.GetInt("length"),
			"",
		}},
	}
}

func (b *bencodeInfo) hash() [20]byte {
	if len(b.filesInfo) == 1 {
		infoDict := map[string]interface{}{
			"name":         b.name,
			"pieces":       b.filesInfo[0].pieces,
			"piece length": b.filesInfo[0].pieceLength,
			"length":       b.filesInfo[0].length,
		}
		return sha1.Sum(bencode.Encode(infoDict))
	}

	infoDict := map[string]interface{}{
		"name":  b.name,
		"files": b.filesInfo,
	}
	return sha1.Sum(bencode.Encode(infoDict))
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

func (b *bencodeTorrent) toTorrentFile() ([]Torrent, error) {
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

func (m bencodeMap) GetString(key string) string {
	v, ok := m[key]
	if !ok {
		return ""
	}
	rv, ok := v.(string)
	if !ok {
		return ""
	}
	return rv
}

func (m bencodeMap) GetList(key string) []interface{} {
	v, ok := m[key]
	if !ok {
		return nil
	}
	rv, ok := v.([]interface{})
	if !ok {
		return nil
	}
	return rv
}

func (m bencodeMap) GetInt(key string) int {
	v, ok := m[key]
	if !ok {
		return 0
	}
	rv, ok := v.(int64)
	if !ok {
		return 0
	}
	return int(rv)
}

func (m bencodeMap) GetDict(key string) bencodeMap {
	v, ok := m[key]
	if !ok {
		return nil
	}
	rv, ok := v.(map[string]interface{})
	if !ok {
		return nil
	}
	return bencodeMap(rv)
}
