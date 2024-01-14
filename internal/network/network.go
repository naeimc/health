package network

import (
	"net"
	"time"
)

type Listener interface {
	Address() net.Addr
	Listen() error
	Close() error
}

type Packet interface {
	Timestamp() time.Time
	RemoteAddress() net.Addr
	LocalAddress() net.Addr
	Data() []byte
	Error() error
	Responder() Responder
}

type Responder interface {
	Write([]byte) (int, error)
}
