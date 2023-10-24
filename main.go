package main

import (
	"fmt"
	"time"

	"github.com/labstack/gommon/log"
	"github.com/raghavroy145/chan-TCP-nels/tcpchan"
)

func main() {
	channelLocal, err := tcpchan.New[string](":3000", ":4000")
	if err != nil {
		log.Fatal(err)
	}
	// The reason we are using sleep here (and below) is because in a real
	// scenario, we would would have the send and receive pairs on different
	// servers, on the same machine we need to wait to ensure everyone is
	// ready and we are not blocking on a server

	go func() {
		time.Sleep(time.Second * 5)
		channelLocal.SendCh <- "royboi"
	}()
	channelRemote, err := tcpchan.New[string](":4000", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	time.Sleep(time.Second * 3)
	msg := <-channelRemote.RecvCh
	fmt.Println("received from channel over the wire: ", msg)
}
