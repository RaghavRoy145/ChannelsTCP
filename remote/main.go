package main

import (
	"fmt"
	"log"

	"github.com/raghavroy145/chan-TCP-nels/tcpchan"
)

func main() {
	channel, err := tcpchan.New[string](":4000", ":3000")
	if err != nil {
		log.Fatal(err)
	}

	msg := <-channel.RecvCh

	fmt.Println("received msg from channel (:3000) over TCP: ", msg)

	channel.SendCh <- msg
}
