package tcpchan

import (
	"encoding/gob"
	"log"
	"net"
	"sync"
	"time"
)

// generics, woo
type TCPChan[T any] struct {
	listenAddr string
	remoteAddr string

	SendCh       chan T
	RecvCh       chan T
	outboundConn net.Conn
	ln           net.Listener
	wg           sync.WaitGroup
}

func New[T any](listenAddr, remoteAddr string) (*TCPChan[T], error) {
	tcpchan := &TCPChan[T]{
		listenAddr: listenAddr,
		remoteAddr: remoteAddr,
		SendCh:     make(chan T, 10),
		RecvCh:     make(chan T, 10),
	}
	ln, err := net.Listen("tcp", listenAddr)
	if err != nil {
		return nil, err
	}
	tcpchan.ln = ln

	go tcpchan.loop()
	go tcpchan.acceptLoop()
	go tcpchan.dialRemoteAndRead()
	return tcpchan, nil
}

func (t *TCPChan[T]) loop() {
	t.wg.Wait()

	for {
		msg := <-t.SendCh
		if err := gob.NewEncoder(t.outboundConn).Encode(&msg); err != nil {
			log.Println(err)
			return
		}
		log.Println("sent msg over the wire: ", msg)
	}
}

func (t *TCPChan[T]) acceptLoop() {
	defer func() {
		t.ln.Close()
	}()
	for {
		conn, err := t.ln.Accept()
		if err != nil {
			log.Println("Accept error", err)
			return
		}
		log.Printf("sender connected %s", conn.RemoteAddr())
		go t.handleConn(conn)
	}
}

func (t *TCPChan[T]) handleConn(conn net.Conn) {
	defer func() {
		t.ln.Close()
	}()
	for {
		var msg T
		if err := gob.NewDecoder(conn).Decode(&msg); err != nil {
			log.Println(err)
			return
		}
		t.RecvCh <- msg
	}
}

func (t *TCPChan[T]) dialRemoteAndRead() {
	t.wg.Add(1)
	for {
		conn, err := net.Dial("tcp", t.remoteAddr)
		if err != nil {
			log.Printf("dial error (%s) retrying: ", err)
			time.Sleep(time.Second * 3)
		} else {
			t.outboundConn = conn
			break
		}
	}
	t.wg.Done()
}
