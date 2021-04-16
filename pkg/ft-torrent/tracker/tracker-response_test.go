package tracker

import (
	"testing"

	"github.com/MonsieurTa/hypertube/pkg/ft-torrent/bencode"
	"github.com/stretchr/testify/assert"
)

func TestTrackerResponse(t *testing.T) {
	data := map[string]interface{}{
		"interval": int64(900),
		"peers":    []byte{},
	}

	trResp := NewTrackerResponse(bencode.Decoder(data))

	expect := TrackerResponse{
		interval: 900,
		peers:    []byte{},
	}
	assert.Equal(t, expect, trResp)
}
