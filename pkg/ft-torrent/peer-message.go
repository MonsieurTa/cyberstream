package torrent

import (
	"encoding/binary"
	"io"
)

const (
	CHOKE          messageID = 0
	UNCHOKE        messageID = 1
	INTERESTED     messageID = 2
	NOT_INTERESTED messageID = 3
	HAVE           messageID = 4
	BITFIELD       messageID = 5
	REQUEST        messageID = 6
	PIECE          messageID = 7
	CANCEL         messageID = 8
)

type messageID uint8

type Bitfield []byte

type Message struct {
	id      messageID
	payload []byte
}

func (m *Message) ID() messageID {
	return m.id
}

func (m *Message) Serialize() []byte {
	len := uint32(len(m.payload) + 1)

	buf := make([]byte, len+4)

	binary.BigEndian.PutUint32(buf[0:4], len)
	buf[4] = byte(m.id)
	copy(buf[5:], m.payload)

	return buf
}

func ReadMessage(r io.Reader) (*Message, error) {
	lengthBuf := make([]byte, 4)
	_, err := io.ReadFull(r, lengthBuf)
	if err != nil {
		return nil, err
	}
	length := binary.BigEndian.Uint32(lengthBuf)

	// keep-alive message
	if length == 0 {
		return nil, nil
	}

	messageBuf := make([]byte, length)
	_, err = io.ReadFull(r, messageBuf)
	if err != nil {
		return nil, err
	}

	m := Message{
		id:      messageID(messageBuf[0]),
		payload: messageBuf[1:],
	}

	return &m, nil
}

func (b Bitfield) HasPiece(index int) bool {
	byteIndex := index / 8
	offset := index % 8
	return b[byteIndex]>>(7-offset)&1 != 0
}

func (bf Bitfield) SetPiece(index int) {
	byteIndex := index / 8
	offset := index % 8
	bf[byteIndex] |= 1 << (7 - offset)
}
