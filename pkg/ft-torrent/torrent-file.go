package torrent

type TorrentFile struct {
	Announce     string
	AnnounceList []string
	InfoHash     [20]byte
	Name         string
	PieceHashes  [][20]byte
	PieceLength  int
	Length       int
}
