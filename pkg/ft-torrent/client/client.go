package client

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/MonsieurTa/hypertube/pkg/ft-torrent/common"
	"github.com/MonsieurTa/hypertube/pkg/ft-torrent/message"
)

type Client struct {
	conn     net.Conn
	bitfield common.Bitfield

	peerID [20]byte
	peer   common.Peer
	State  State
}

type State struct {
	AmChoking      bool
	AmInterested   bool
	PeerChocking   bool
	PeerInterested bool
}

var (
	err_nil_msg         = fmt.Errorf("expected bitfield, got nil message")
	err_expect_bitfield = func(msg *message.Message) error { return fmt.Errorf("expected bitfield, got ID %d", msg.ID()) }
	err_expect_infohash = func(expect, got []byte) error { return fmt.Errorf("expected info hash %x but got %x", expect, got) }
)

func NewClient(peer common.Peer, peerID [20]byte) *Client {
	return &Client{
		conn:     nil,
		bitfield: nil,
		peerID:   peerID,
		peer:     peer,
		State:    defaultState(),
	}
}

func (c *Client) Bitfield() common.Bitfield {
	return c.bitfield
}

func (c *Client) Conn() net.Conn {
	return c.conn
}

func (c *Client) Close() {
	c.conn.Close()
}

func (c *Client) Ask(infoHash [20]byte) error {
	conn, err := net.DialTimeout("tcp", c.peer.String(), 5*time.Second)
	if err != nil {
		return err
	}

	err = shakeHand(conn, c.peerID, infoHash)
	if err != nil {
		conn.Close()
		return err
	}

	bf, err := recvBitfield(conn)
	if err != nil {
		conn.Close()
		return err
	}
	c.bitfield = bf
	c.conn = conn

	log.Printf("connection established with %s\n", c.peer.String())
	return nil
}

func defaultState() State {
	return State{
		AmChoking:      true,
		AmInterested:   false,
		PeerChocking:   true,
		PeerInterested: false,
	}
}

func shakeHand(Conn net.Conn, peerID, infoHash [20]byte) error {
	hs := NewHandShake(infoHash, peerID)

	_, err := Conn.Write(hs.Serialize())
	if err != nil {
		return err
	}

	resp, err := ReadHandshake(Conn)

	if err != nil {
		return errors.New("could not read peer handshake")
	}

	if !bytes.Equal(resp.InfoHash(), infoHash[:]) {
		return err_expect_infohash(infoHash[:], resp.InfoHash())
	}
	return nil
}

func recvBitfield(conn net.Conn) (common.Bitfield, error) {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Now())

	msg, err := message.ReadMessage(conn)
	if err != nil {
		return nil, err
	}

	if msg == nil {
		return nil, err_nil_msg
	}
	if msg.ID() != message.BITFIELD {
		return nil, err_expect_bitfield(msg)
	}

	return msg.Payload(), nil
}

func (c *Client) Read() (*message.Message, error) {
	return message.ReadMessage(c.conn)
}

func (c *Client) SendRequest(index, start, length int) error {
	req := message.Request(index, start, length)
	_, err := c.conn.Write(req.Serialize())
	return err
}

func (c *Client) SendInterested() error {
	msg := message.Interested()
	_, err := c.conn.Write(msg.Serialize())
	return err
}

func (c *Client) SendNotInterested() error {
	msg := message.NotInterested()
	_, err := c.conn.Write(msg.Serialize())
	return err
}

func (c *Client) SendUnchoke() error {
	msg := message.Unchoke()
	_, err := c.conn.Write(msg.Serialize())
	return err
}

func (c *Client) SendHave(index int) error {
	msg := message.Have(index)
	_, err := c.conn.Write(msg.Serialize())
	return err
}
