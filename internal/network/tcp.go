package network

import (
	"errors"
	"net"
	"time"
)

const TCP = "tcp"

type TCPListener struct {
	Log func(any)

	listener net.Listener
	queue    chan<- Packet
}

func NewTCPListener(address string, queue chan<- Packet) (*TCPListener, error) {
	listener, err := net.Listen(TCP, address)
	if err != nil {
		return nil, err
	}
	return &TCPListener{func(a any) {}, listener, queue}, nil
}

func (l *TCPListener) Address() net.Addr {
	return l.listener.Addr()
}

func (l *TCPListener) Listen() error {
	for {
		connection, err := l.listener.Accept()
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				return nil
			}
			return err
		}

		go func(connection net.Conn) {
			buffer := make([]byte, 1<<16)
			for {
				size, err := connection.Read(buffer)
				if err != nil && err.Error() == "EOF" {
					return
				}
				if size > 0 && l.queue != nil {
					data := make([]byte, size)
					copy(data, buffer)
					l.queue <- TCPPacket{
						timestamp:  time.Now().UTC(),
						connection: connection,
						data:       data,
						err:        err,
					}
				}
			}
		}(connection)
	}
}

func (l *TCPListener) Close() error {
	return l.listener.Close()
}

type TCPPacket struct {
	timestamp  time.Time
	connection net.Conn
	data       []byte
	err        error
}

func (p TCPPacket) Timestamp() time.Time {
	return p.timestamp
}

func (p TCPPacket) RemoteAddress() net.Addr {
	return p.connection.RemoteAddr()
}

func (p TCPPacket) LocalAddress() net.Addr {
	return p.connection.LocalAddr()
}

func (p TCPPacket) Data() []byte {
	return p.data
}

func (p TCPPacket) Error() error {
	return p.err
}

func (p TCPPacket) Responder() Responder {
	return TCPResponder{p.connection}
}

type TCPResponder struct {
	connection net.Conn
}

func (r TCPResponder) Write(b []byte) (int, error) {
	return r.connection.Write(b)
}
