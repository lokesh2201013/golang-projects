package p2p

import (
	"net"
	"sync"
)

type TCPTransport struct {
	listenAdrr string
	listener   net.Listener
	mu sync.RWMutex
	peers map[net.Addr]Peer
}

func NewTCPTransport(listenAddr string) *TCPTransport{
	return &TCPTransport{
		listenAdrr: listenAddr ,
	}
}
