package tcpchan

import (
	"encoding/gob"
	"log"
	"net"
	"time"
)

var defaultDialInterval = 3 * time.Second

// generics, woo
type Sender[T any] struct {
	Chan chan T

	remoteAddr   string
	outboundConn net.Conn
	dialInterval time.Duration
}

type Receiver[T any] struct {
	Chan       chan T
	listenAddr string
	listener   net.Listener
}

func NewSender[T any](remoteAddr string) (*Sender[T], error) {
	sender := &Sender[T]{
		Chan:         make(chan T),
		remoteAddr:   remoteAddr,
		dialInterval: defaultDialInterval,
	}
	go sender.dialRemote()
	go sender.loop()

	return sender, nil
}

func NewReceiver[T any](listenAddr string) (*Receiver[T], error) {
	recv := &Receiver[T]{
		Chan:       make(chan T),
		listenAddr: listenAddr,
	}

	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	recv.listener = ln

	go recv.acceptLoop()
	return recv, nil
}

func (r *Receiver[T]) acceptLoop() {
	defer func() {
		r.listener.Close()
	}()
	for {
		conn, err := r.listener.Accept()
		if err != nil {
			log.Println("Accept error", err)
			return
		}
		log.Printf("sender connected %s", conn.RemoteAddr())
		go r.handleConn(conn)
	}
}

func (r *Receiver[T]) handleConn(conn net.Conn) {
	for {
		var msg T
		if err := gob.NewDecoder(conn).Decode(&msg); err != nil {
			log.Println(err)
			continue
		}
		r.Chan <- msg
	}
}

func (s *Sender[T]) dialRemote() {
	conn, err := net.Dial("tcp", s.remoteAddr)
	if err != nil {
		log.Printf("dial error (%s) retrying in (%d): ", err, s.dialInterval)
		time.Sleep(s.dialInterval)
		s.dialRemote()
	}
	s.outboundConn = conn
}

func (s *Sender[T]) loop() {
	for {
		msg := <-s.Chan
		log.Println("sending msg over the wire: ", msg)
		if err := gob.NewEncoder(s.outboundConn).Encode(msg); err != nil {
			log.Println(err)
		}
	}
}
