package network

import (
	"fmt"
	"sync"
)

type Server struct {
	Handle func(Packet) ([]byte, error)
	Log    func(any)

	listeners []Listener

	wait     *sync.WaitGroup
	queue    <-chan Packet
	shutdown chan any
}

func NewServer(queue <-chan Packet) *Server {
	return &Server{
		Handle: func(p Packet) ([]byte, error) { return nil, nil },
		Log:    func(a any) {},

		listeners: make([]Listener, 0),

		queue:    queue,
		wait:     &sync.WaitGroup{},
		shutdown: make(chan any, 1),
	}
}

func (s *Server) Serve() error {
	s.Log("starting server")
	for _, l := range s.listeners {
		s.Log(fmt.Sprintf("starting listener: %s/%s", l.Address().Network(), l.Address().String()))
		go l.Listen()
	}

	for {
		select {
		case <-s.shutdown:
			for _, l := range s.listeners {
				s.Log(fmt.Sprintf("stopping listener: %s/%s", l.Address().Network(), l.Address().String()))
				l.Close()
			}
			s.Log("stopping server")
			s.wait.Wait()
			return nil
		case packet := <-s.queue:
			s.wait.Add(1)
			go func(packet Packet) {
				s.wait.Done()
				response, err := s.Handle(packet)
				if len(response) > 0 {
					packet.Responder().Write(response)
				}
				if err != nil {
					fmt.Println(err)
				}
			}(packet)
		}
	}
}

func (s *Server) Listener(l Listener) {
	s.listeners = append(s.listeners, l)
}

func (s *Server) Shutdown() {
	s.shutdown <- nil
}
