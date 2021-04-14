package torrent

import "net"

type Peer struct {
	ip   net.IP
	port uint16
}

func NewPeer(ip net.IP, port uint16) Peer {
	return Peer{ip, port}
}

func (p *Peer) SetIP(ip net.IP) {
	p.ip = ip
}

func (p *Peer) SetPort(port uint16) {
	p.port = port
}

func (p Peer) IP() net.IP {
	return p.ip
}

func (p Peer) Port() uint16 {
	return p.port
}
