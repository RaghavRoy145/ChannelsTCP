package main

import (
	"fmt"

	"github.com/labstack/gommon/log"
	"github.com/raghavroy145/chan-TCP-nels/tcpchan"
)

func main() {
	channel, err := tcpchan.New[string](":3000", ":4000")
	if err != nil {
		log.Fatal(err)
	}

	channel.SendCh <- "royboi"
	msg := channel.RecvCh
	fmt.Println("received msg from channel (:4000) over TCP: ", msg)
	select {}
}
