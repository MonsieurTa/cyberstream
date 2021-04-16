package torrent

type pieceWork struct {
	index  int
	hash   [20]byte
	length int
}

type pieceResult struct {
	index int
	buf   []byte
}

func (t *Torrent) Download() error {
	trackers, err := t.Trackers()
	if err != nil {
		return err
	}

	workQueue := make(chan *pieceWork, len(t.PieceHashes))
	results := make(chan *pieceResult)
	defer close(workQueue)
	defer close(results)

	for i, hash := range t.PieceHashes {
		length := t.calculatePieceSize(i)
		workQueue <- &pieceWork{i, hash, length}
	}

	trackers.RequestPeers(t.InfoHash)
	return nil
}

func (t *Torrent) calculatePieceBounds(index int) (int, int) {
	start := index * t.PieceLength
	end := start + t.PieceLength
	if end > t.Length {
		end = t.Length
	}
	return start, end
}

func (t *Torrent) calculatePieceSize(index int) int {
	start, end := t.calculatePieceBounds(index)
	return end - start
}
