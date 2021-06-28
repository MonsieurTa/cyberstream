package tcp

import (
	"log"
	"unsafe"
)

type Header struct {
	FileSize    uint64
	PieceLength uint64
	DirNameSize uint64
	DirName     []byte
}

type Response struct {
	PlaylistURLLength uint64
	PlaylistURL       []byte
}

// Size of FileSize + PieceLength + DirNameSize data types
func HeaderMetadataSize() uint64 {
	var dummy uint64
	return uint64(unsafe.Sizeof(dummy)) * 3
}

func (h *Header) Print() {
	log.Println(h)
}
