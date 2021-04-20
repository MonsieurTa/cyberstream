package download

import "github.com/MonsieurTa/hypertube/pkg/ft-torrent/common"

type Torrent interface {
	InfoHash() [20]byte
	Pieces() int
	PieceHashes() [][20]byte
	PieceBounds(index int) (int, int)
	PieceSize(index int) int
	Length() int
	Peers(peerID [20]byte) ([]common.Peer, error)
}
