package torrent

import (
	"crypto/sha1"

	"github.com/marksamman/bencode"
)

func hash(b *bencodeInfo) [20]byte {
	if len(b.filesInfo) == 0 {
		return [20]byte{}
	}
	if len(b.filesInfo) == 1 {
		return hashOne(b)
	}
	return hashMultiple(b)
}

func hashOne(b *bencodeInfo) [20]byte {
	infoDict := map[string]interface{}{
		"name":         b.name,
		"pieces":       b.filesInfo[0].pieces,
		"piece length": b.filesInfo[0].pieceLength,
		"length":       b.filesInfo[0].length,
	}
	return sha1.Sum(bencode.Encode(infoDict))
}

func hashMultiple(b *bencodeInfo) [20]byte {
	infoDict := map[string]interface{}{
		"name":  b.name,
		"files": b.filesInfo,
	}
	return sha1.Sum(bencode.Encode(infoDict))
}
