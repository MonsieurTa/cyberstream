package client

import (
	"bytes"
	"fmt"
	"net"
	"time"

	"github.com/MonsieurTa/hypertube/pkg/ft-torrent/handshake"
	"github.com/MonsieurTa/hypertube/pkg/ft-torrent/message"
	"github.com/MonsieurTa/hypertube/pkg/ft-torrent/peer"
)

type Client struct {
	conn     net.Conn
	bitfield message.Bitfield
	state    State
}

type State struct {
	amChoking      bool
	amInterested   bool
	peerChocking   bool
	peerInterested bool
}

var (
	err_nil_msg         = fmt.Errorf("expected bitfield, got nil message")
	err_expect_bitfield = func(msg *message.Message) error { return fmt.Errorf("expected bitfield, got ID %d", msg.ID()) }
	err_expect_infohash = func(expect, got []byte) error { return fmt.Errorf("expected info hash %x but got %x", expect, got) }
)

func NewClient(peer peer.Peer, peerID, infoHash [20]byte) (*Client, error) {
	conn, err := net.DialTimeout("tcp", peer.String(), 5*time.Second)
	if err != nil {
		return nil, err
	}

	_, err = shakeHand(conn, infoHash, peerID)
	if err != nil {
		conn.Close()
		return nil, err
	}

	bf, err := recvBitfield(conn)
	if err != nil {
		return nil, err
	}
	return &Client{
		conn:     conn,
		bitfield: bf,
		state:    defaultState(),
	}, nil
}

func defaultState() State {
	return State{
		amChoking:      true,
		amInterested:   false,
		peerChocking:   true,
		peerInterested: false,
	}
}

func shakeHand(conn net.Conn, peerID, infoHash [20]byte) (*handshake.Handshake, error) {
	hs := handshake.NewHandShake(infoHash, peerID)

	_, err := conn.Write(hs.Serialize())
	if err != nil {
		return nil, err
	}

	resp, err := handshake.Read(conn)
	if err != nil {
		return nil, err
	}
	if !bytes.Equal(resp.InfoHash(), infoHash[:]) {
		return nil, err_expect_infohash(infoHash[:], resp.InfoHash())
	}
	return &resp, nil
}

func recvBitfield(conn net.Conn) (message.Bitfield, error) {
	conn.SetDeadline(time.Now().Add(5 * time.Second))
	defer conn.SetDeadline(time.Now())

	msg, err := message.Read(conn)
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
	msg, err := message.Read(c.conn)
	return msg, err
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
