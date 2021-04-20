package message

import (
	"encoding/binary"
	"fmt"
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

var (
	err_wrong_size = func(expected, got int) error { return fmt.Errorf("expected size: %d, got %d", expected, got) }

	err_wrong_msg         = func(expected, got messageID) error { return fmt.Errorf("expected ID: %d, got %d", expected, got) }
	err_payload_too_short = func(len int) error { return fmt.Errorf("payload too short. %d < 8", len) }
	err_wrong_index       = func(expected, got int) error { return fmt.Errorf("expected index %d, got %d", expected, got) }
	err_wrong_offset      = func(offset, limit int) error { return fmt.Errorf("offset too high. %d >= %d", offset, limit) }
	err_wrong_datalen     = func(datalen, offset, limit int) error {
		return fmt.Errorf("data too long [%d] for offset %d with length %d", datalen, offset, limit)
	}
)

func (m *Message) ID() messageID {
	return m.id
}

func (m *Message) Payload() []byte {
	return m.payload
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

func Choke() *Message {
	return &Message{id: CHOKE}
}

func Unchoke() *Message {
	return &Message{id: UNCHOKE}
}

func Interested() *Message {
	return &Message{id: INTERESTED}
}

func NotInterested() *Message {
	return &Message{id: NOT_INTERESTED}
}

func Request(index, start, length int) *Message {
	payload := make([]byte, 12)
	binary.BigEndian.PutUint32(payload[0:4], uint32(index))
	binary.BigEndian.PutUint32(payload[4:8], uint32(start))
	binary.BigEndian.PutUint32(payload[8:12], uint32(length))
	return &Message{
		id:      REQUEST,
		payload: payload,
	}
}

func Have(index int) *Message {
	payload := make([]byte, 4)
	binary.BigEndian.PutUint32(payload, uint32(index))
	return &Message{
		id:      HAVE,
		payload: payload,
	}
}

func ParseHave(msg *Message) (int, error) {
	const size = 4

	if len(msg.payload) != size {
		return 0, err_wrong_size(size, len(msg.payload))
	}

	rv := int(binary.BigEndian.Uint32(msg.payload))
	return rv, nil
}

func ParsePiece(index int, buf []byte, msg *Message) (int, error) {
	if msg.ID() != PIECE {
		return 0, err_wrong_msg(PIECE, msg.ID())
	}

	payload := msg.Payload()
	if len(payload) < 8 {
		return 0, err_payload_too_short(len(payload))
	}

	parsedIndex := int(binary.BigEndian.Uint32(payload[0:4]))
	if parsedIndex != index {
		return 0, err_wrong_index(index, parsedIndex)
	}

	offset := int(binary.BigEndian.Uint32(payload[4:8]))
	if offset >= len(buf) {
		return 0, err_wrong_offset(offset, len(buf))
	}

	data := payload[8:]
	if offset+len(data) > len(buf) {
		return 0, err_wrong_datalen(offset+len(data), offset, len(buf))
	}
	copy(buf[offset:], data)

	return len(data), nil
}
